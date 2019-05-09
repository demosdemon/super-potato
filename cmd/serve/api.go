package serve

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func api(group gin.IRoutes) gin.IRoutes {
	logrus.Trace("api")
	group = addRoutes(group)
	group.Any("/anything/*path", anything)
	group.GET("/env", listEnv)
	group.GET("/env/:name", getEnv)
	group.GET("/ping", ping)
	group.GET("/user", getUser)
	group.POST("/env/:name", setEnv)

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

func getUser(c *gin.Context) {
	user := getCertifiedUser(c)
	if user == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	negotiate(c, http.StatusOK, user)
}
