package server

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"

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
	group.GET("/ping", ping)
	group.GET("/routes", getRoutes)
	group.GET("/application", getApplication)
	group.GET("/relationships", getRelationships)
	group.GET("/env", listEnv)
	group.GET("/env/:name", getEnv)
	group.POST("/env/:name", setEnv)
	return group
}

func ping(c *gin.Context) {
	logrus.Trace("ping")
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "pong",
		"ts":      time.Now().UTC(),
	})
}

func getRoutes(c *gin.Context) {
	logrus.Trace("getRoutes")
	if routes, err := env.Routes(); err == nil {
		c.IndentedJSON(http.StatusOK, routes)
	} else {
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
}

func getApplication(c *gin.Context) {
	logrus.Trace("getApplication")
	if app, err := env.Application(); err == nil {
		c.IndentedJSON(http.StatusOK, app)
	} else {
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
}

func getRelationships(c *gin.Context) {
	logrus.Trace("getRelationships")
	if rels, err := env.Relationships(); err == nil {
		c.IndentedJSON(http.StatusOK, rels)
	} else {
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
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
