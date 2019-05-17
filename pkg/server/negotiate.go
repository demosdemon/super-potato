package server

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/render"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const prettyMarkdownTemplate = "# %s\n\n```%s\n%s\n```\n"

type pretty struct {
	Markdown
	Title  string
	Format string
	Data   interface{}

	bufMu sync.Mutex
	buf   io.Reader
}

func newPretty(c *gin.Context, code int, data interface{}) *pretty {
	format := c.NegotiateFormat(
		binding.MIMEJSON,
		binding.MIMEYAML,
	)

	logrus.WithField("format", format).Trace("pretty")

	p := pretty{
		Title:  http.StatusText(code),
		Format: format,
		Data:   data,
	}
	p.Markdown.input = &p
	return &p
}

func (s *Server) negotiate(c *gin.Context, code int, data interface{}) {
	format := c.NegotiateFormat(
		binding.MIMEJSON,
		binding.MIMEHTML,
		binding.MIMEXML,
		binding.MIMEXML2,
		binding.MIMEYAML,
	)

	renderer := getRenderer(format, data, c, code)

	c.Render(code, renderer)
}

func getRenderer(format string, data interface{}, c *gin.Context, code int) render.Render {
	logrus.WithField("format", format).WithField("data", data).Trace("getRenderer")

	switch format {
	case binding.MIMEJSON:
		return render.IndentedJSON{Data: data}
	case binding.MIMEHTML:
		return newPretty(c, code, data)
	case binding.MIMEXML, binding.MIMEXML2:
		return render.XML{Data: data}
	case binding.MIMEYAML:
		return render.YAML{Data: data}
	default:
		logrus.WithField("format", format).Warn("unknown format")
		return render.JSON{Data: data}
	}
}

func (p *pretty) Read(buf []byte) (int, error) {
	logrus.Trace("pretty read")

	p.bufMu.Lock()
	defer p.bufMu.Unlock()

	if p.buf == nil {
		r, err := p.render()
		if err != nil {
			logrus.WithError(err).Trace()
			return 0, err
		}
		logrus.WithField("r", r).Trace("post render")
		p.buf = strings.NewReader(r)
	}

	return p.buf.Read(buf)
}

func (p *pretty) render() (string, error) {
	logrus.Trace("pretty render")

	var format string
	var formatted []byte
	var err error
	switch p.Format {
	case binding.MIMEJSON:
		format = "json"
		formatted, err = json.MarshalIndent(p.Data, "", "    ")
	case binding.MIMEXML, binding.MIMEXML2:
		format = ""
		formatted, err = xml.MarshalIndent(p.Data, "", "    ")
	case binding.MIMEYAML:
		format = "yaml"
		formatted, err = yaml.Marshal(p.Data)
	}

	if err != nil {
		return "", err
	}

	return fmt.Sprintf(prettyMarkdownTemplate, p.Title, format, string(formatted)), nil
}
