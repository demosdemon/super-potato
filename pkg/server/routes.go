package server

import (
	"bytes"
	"encoding/json"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/demosdemon/super-potato/pkg/platformsh"
)

func (s *Server) root(c *gin.Context) {
	fp, err := s.Open("/README.md")
	if err != nil {
		if os.IsNotExist(err) {
			c.AbortWithError(404, err)
		} else {
			c.AbortWithError(500, err)
		}
		return
	}
	defer fp.Close()

	r := NewMarkdown(fp)
	c.Render(200, r)
}

func (s *Server) getPing(c *gin.Context) {
	rv := gin.H{
		"now":     time.Now(),
		"query":   c.Request.URL.Query(),
		"ip":      c.ClientIP(),
		"message": "pong",
	}
	logrus.WithFields(logrus.Fields(rv)).Trace("getPing")

	s.negotiate(c, http.StatusOK, rv)
}

func (s *Server) getUser(c *gin.Context) {
	rv := gin.H{
		"user": getUser(c),
	}
	logrus.WithFields(logrus.Fields(rv)).Trace("getUser")
	s.negotiate(c, http.StatusOK, rv)
}

func (s *Server) getDebugVars(c *gin.Context) {
	var buf bytes.Buffer
	buf.WriteString("{")
	first := true
	expvar.Do(func(kv expvar.KeyValue) {
		if !first {
			buf.WriteString(",")
		}
		first = false
		fmt.Fprintf(&buf, "%q: %s", kv.Key, kv.Value)
	})
	buf.WriteString("}")

	var kvp gin.H
	if err := json.Unmarshal(buf.Bytes(), &kvp); err != nil {
		c.AbortWithError(500, err)
	}

	s.negotiate(c, 200, kvp)
}

func (s *Server) getFaviconICO(c *gin.Context) {
	logo := platformsh.NewRasterLogo()
	logo.Size = 32
	render := logo.Negotiate(c)
	c.Render(http.StatusOK, render)
}

func (s *Server) getLogoSVG(c *gin.Context) {
	logo := platformsh.NewLogoSVG()
	logo.Background = c.DefaultQuery("background", logo.Background)
	logo.Foreground = c.DefaultQuery("foreground", logo.Foreground)
	c.Render(http.StatusOK, logo)
}

func (s *Server) getLogoPNG(c *gin.Context) {
	logo := platformsh.NewRasterLogo()
	logo.Size = 100
	render := platformsh.RenderRasterLogo{
		RasterLogo:  logo,
		ContentType: platformsh.FormatPNG,
	}
	c.Render(http.StatusOK, render)
}

func (s *Server) getHeaders(c *gin.Context) {
	h := gin.H{}
	for k, v := range c.Request.Header {
		h[k] = strings.Join(v, "; ")
	}
	s.negotiate(c, 200, h)
}
