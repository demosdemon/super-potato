package server

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/render"
	"github.com/sirupsen/logrus"

	"github.com/demosdemon/super-potato/pkg/platformsh"
)

var env = platformsh.DefaultEnvironment

func Execute() {
	l, err := env.Listener()
	if err != nil {
		panic(err)
	}

	engine := gin.New()
	engine.Use(
		gin.Logger(),
		gin.Recovery(),
	)

	api(engine.Group("/api"))

	_ = http.Serve(l, engine)
}

func api(group gin.IRoutes) gin.IRoutes {
	logrus.Trace("api")
	group = addRoutes(group)
	group.GET("/ping", ping)
	group.GET("/env", listEnv)
	group.GET("/env/:name", getEnv)
	group.POST("/env/:name", setEnv)
	group.Any("/anything/*path", anything)

	return group
}

func ping(c *gin.Context) {
	logrus.Trace("ping")
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "pong",
		"ts":      time.Now().UTC(),
	})
}

func listEnv(c *gin.Context) {
	logrus.Trace("listEnv")
	environ := os.Environ()
	keys := make([]string, len(environ))
	for idx, kvp := range environ {
		keys[idx] = strings.SplitN(kvp, "=", 2)[0]
	}

	c.IndentedJSON(http.StatusOK, keys)
}

func getEnv(c *gin.Context) {
	logrus.Trace("getEnv")
	name := c.Param("name")
	decodeQuery := c.Query("decode")
	var decode bool
	_ = json.Unmarshal([]byte(decodeQuery), &decode)

	if val, ok := os.LookupEnv(name); ok {
		if decode {
			decoded, err := base64.StdEncoding.DecodeString(val)
			if err != nil {
				c.IndentedJSON(http.StatusBadRequest, gin.H{
					"error":   err,
					"message": "unable to decode base64 value",
				})
				return
			}

			var obj interface{}
			err = json.Unmarshal(decoded, &obj)
			if err != nil {
				c.IndentedJSON(http.StatusBadRequest, gin.H{
					"error":   err,
					"message": "unable to decode JSON value",
				})
				return
			}

			c.IndentedJSON(http.StatusOK, gin.H{name: obj})
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{name: val})
		}
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error": "not found",
			"key":   name,
		})
	}
}

func setEnv(c *gin.Context) {
	logrus.Trace("setEnv")
	name := c.Param("name")
	value, err := c.GetRawData()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "error reading request data",
		})
		return
	}
	err = os.Setenv(name, string(value))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error":   err,
			"message": "error setting environment variable",
		})
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{name: string(value)})
}

func anything(c *gin.Context) {
	body := bind(c)

	data := gin.H{
		"method":            c.Request.Method,
		"url":               c.Request.URL.String(),
		"proto":             c.Request.Proto,
		"headers":           Header(c.Request.Header),
		"body":              body,
		"content_length":    c.Request.ContentLength,
		"transfer_encoding": c.Request.TransferEncoding,
		"host":              c.Request.Host,
		"remote_addr":       c.Request.RemoteAddr,
		"request_uri":       c.Request.RequestURI,
		"client_ip":         c.ClientIP(),
	}

	format := c.NegotiateFormat(
		binding.MIMEJSON,
		binding.MIMEXML,
		binding.MIMEXML2,
		binding.MIMEYAML,
	)

	var renderer render.Render
	switch format {
	case binding.MIMEXML, binding.MIMEXML2:
		renderer = render.XML{Data: data}
	case binding.MIMEYAML:
		renderer = render.YAML{Data: data}
	default:
		renderer = render.IndentedJSON{Data: data}
	}

	c.Render(http.StatusOK, renderer)
}

func bind(c *gin.Context) interface{} {
	var body interface{}
	b := binding.Default(c.Request.Method, c.ContentType())
	if bb, ok := b.(binding.BindingBody); ok {
		if err := c.ShouldBindBodyWith(&body, bb); err != nil || body == nil {
			if data, ok := c.Get(gin.BodyBytesKey); ok {
				if b, ok := data.([]byte); ok {
					body = string(b)
				}
			}
		}
	} else {
		if data, err := c.GetRawData(); err == nil {
			_ = c.Request.Body.Close()
			c.Request.Body = ioutil.NopCloser(bytes.NewReader(data))
			if err := b.Bind(c.Request, &body); err != nil || body == nil {
				body = string(data)
			}
		}
	}
	return body
}

type Header http.Header

func (h Header) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
	start.Attr = append(start.Attr, xml.Attr{
		Name: xml.Name{
			Local: "length",
		},
		Value: fmt.Sprintf("%d", len(h)),
	})
	err = e.EncodeToken(start)
	if err != nil {
		return err
	}

	for k, v := range h {
		keyStart := xml.StartElement{
			Name: xml.Name{
				Local: k,
			},
			Attr: []xml.Attr{
				{
					Name: xml.Name{
						Local: "length",
					},
					Value: fmt.Sprintf("%d", len(v)),
				},
			},
		}
		err = e.EncodeToken(keyStart)
		if err != nil {
			return err
		}

		for idx, value := range v {
			valueStart := xml.StartElement{
				Name: xml.Name{
					Local: "String",
				},
				Attr: []xml.Attr{
					{
						Name: xml.Name{
							Local: "index",
						},
						Value: fmt.Sprintf("%d", idx),
					},
				},
			}
			err = e.EncodeToken(valueStart)
			if err != nil {
				return err
			}

			parsed, err := url.QueryUnescape(value)
			if err != nil {
				parsed = value
			}

			data := xml.CharData(parsed)
			err = e.EncodeToken(data)
			if err != nil {
				return err
			}

			err = e.EncodeToken(valueStart.End())
			if err != nil {
				return err
			}
		}

		err = e.EncodeToken(keyStart.End())
		if err != nil {
			return err
		}
	}

	err = e.EncodeToken(start.End())
	if err != nil {
		return err
	}

	return nil
}
