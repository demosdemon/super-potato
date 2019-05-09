package server

import (
	"bytes"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
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
	group.GET("/user", getUser)

	return group
}

func ping(c *gin.Context) {
	logrus.Trace("ping")
	negotiate(c, http.StatusOK, gin.H{
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

	negotiate(c, http.StatusOK, keys)
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
				negotiate(c, http.StatusBadRequest, gin.H{
					"error":   err,
					"message": "unable to decode base64 value",
				})
				return
			}

			var obj interface{}
			err = json.Unmarshal(decoded, &obj)
			if err != nil {
				negotiate(c, http.StatusBadRequest, gin.H{
					"error":   err,
					"message": "unable to decode JSON value",
				})
				return
			}

			negotiate(c, http.StatusOK, gin.H{name: obj})
		} else {
			negotiate(c, http.StatusOK, gin.H{name: val})
		}
	} else {
		negotiate(c, http.StatusNotFound, gin.H{
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
		negotiate(c, http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "error reading request data",
		})
		return
	}
	err = os.Setenv(name, string(value))
	if err != nil {
		negotiate(c, http.StatusInternalServerError, gin.H{
			"error":   err,
			"message": "error setting environment variable",
		})
		return
	}
	negotiate(c, http.StatusCreated, gin.H{name: string(value)})
}

func anything(c *gin.Context) {
	body := bind(c)

	data := gin.H{
		"method":            c.Request.Method,
		"url":               c.Request.URL.String(),
		"proto":             c.Request.Proto,
		"headers":           Header{c.Request.Header},
		"body":              body,
		"content_length":    c.Request.ContentLength,
		"transfer_encoding": c.Request.TransferEncoding,
		"host":              c.Request.Host,
		"remote_addr":       c.Request.RemoteAddr,
		"request_uri":       c.Request.RequestURI,
		"client_ip":         c.ClientIP(),
	}

	negotiate(c, http.StatusOK, data)
}

func negotiate(c *gin.Context, code int, data interface{}) {
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

	c.Render(code, renderer)
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

type Header struct {
	http.Header
}

func (h Header) Keys() []string {
	keys := make([]string, 0, len(h.Header))
	for k := range h.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (h Header) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Attr = append(start.Attr, xml.Attr{
		Name: xml.Name{
			Local: "length",
		},
		Value: fmt.Sprintf("%d", len(h.Header)),
	})
	err := e.EncodeToken(start)
	if err != nil {
		return err
	}

	keys := h.Keys()

	for _, k := range keys {
		v := h.Get(k)
		keyStart := xml.StartElement{
			Name: xml.Name{
				Local: k,
			},
		}
		err := e.EncodeToken(keyStart)
		if err != nil {
			return err
		}

		value, err := url.QueryUnescape(v)
		if err != nil {
			value = v
		}

		data := xml.CharData(value)
		err = e.EncodeToken(data)
		if err != nil {
			return err
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

type User struct {
	ClientCertificate *x509.Certificate
	DistinguishedName string
}

const UserCacheKey = "github.com/demosdemon/super-potato/cmd/server/CertifiedUser"

func getCertifiedUser(c *gin.Context) *User {
	if v, ok := c.Get(UserCacheKey); ok {
		if v, ok := v.(*User); ok {
			return v
		}
	}

	xClientCert := c.GetHeader("X-Client-Cert")
	if xClientCert == "" {
		return nil
	}

	xClientCert, err := url.QueryUnescape(xClientCert)
	if err != nil {
		panic(err)
	}

	pemBlock, _ := pem.Decode([]byte(xClientCert))
	if pemBlock == nil {
		panic("invalid PEM data")
	}

	var user User

	user.ClientCertificate, err = x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		panic(err)
	}

	user.DistinguishedName = c.GetHeader("X-Client-Dn")
	if user.DistinguishedName == "" {
		user.DistinguishedName = user.ClientCertificate.Subject.String()
	}

	c.Set(UserCacheKey, &user)
	return &user
}

func getUser(c *gin.Context) {
	user := getCertifiedUser(c)
	if user == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	negotiate(c, http.StatusOK, user)
}
