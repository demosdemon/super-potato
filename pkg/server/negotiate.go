package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/render"
	"github.com/russross/blackfriday/v2"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func (s *Server) negotiate(c *gin.Context, code int, data interface{}) {
	format := c.NegotiateFormat(
		binding.MIMEJSON,
		binding.MIMEHTML,
		binding.MIMEXML,
		binding.MIMEXML2,
		binding.MIMEYAML,
	)

	var renderer render.Render
	switch format {
	case binding.MIMEJSON:
		renderer = render.IndentedJSON{Data: data}
	case binding.MIMEHTML:
		renderer = pretty{Data: data}
	case binding.MIMEXML, binding.MIMEXML2:
		renderer = render.XML{Data: data}
	case binding.MIMEYAML:
		renderer = render.YAML{Data: data}
	default:
		logrus.WithField("format", format).Warn("unknown format")
		renderer = pretty{Data: data}
	}

	c.Render(code, renderer)
}

type pretty struct {
	Data interface{}
}

func (p pretty) Render(w http.ResponseWriter) error {
	const tpl = "# Result\n\n```yaml\n%s\n```\n"
	data, err := yaml.Marshal(p.Data)
	if err != nil {
		return err
	}

	md := fmt.Sprintf(tpl, string(data))
	html := blackfriday.Run([]byte(md))
	_, err = w.Write(html)
	return err
}

func (p pretty) WriteContentType(w http.ResponseWriter) {
	if v, ok := w.Header()["Content-Type"]; ok && len(v) > 0 {
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}
