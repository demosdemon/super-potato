package server

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
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
	var body interface{}

	if err := c.ShouldBind(&body); err != nil {
		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
			return
		}
		body = string(data)
	}

	fmt := c.NegotiateFormat(
		binding.MIMEJSON,
		binding.MIMEXML,
		binding.MIMEXML2,
		binding.MIMEYAML,
	)

	data := gin.H{
		"method":            c.Request.Method,
		"url":               c.Request.URL.String(),
		"proto":             c.Request.Proto,
		"headers":           c.Request.Header,
		"body":              body,
		"content_length":    c.Request.ContentLength,
		"transfer_encoding": c.Request.TransferEncoding,
		"host":              c.Request.Host,
		"remote_addr":       c.Request.RemoteAddr,
		"request_uri":       c.Request.RequestURI,
		"client_ip":         c.ClientIP(),
	}

	var renderer render.Render
	switch fmt {
	case binding.MIMEXML, binding.MIMEXML2:
		renderer = render.XML{Data: data}
	case binding.MIMEYAML:
		renderer = render.YAML{Data: data}
	default:
		renderer = render.IndentedJSON{Data: data}
	}

	c.Render(http.StatusOK, renderer)
}
