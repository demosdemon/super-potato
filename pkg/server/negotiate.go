package server

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/render"
	"github.com/russross/blackfriday/v2"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	prettyMarkdownTemplate = "# %s\n\n```%s\n%s\n```\n"
	prettyHTMLTemplate     = `<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<title>Result</title>
	<style>
		html {
			font-family: "Helvetica Neue", Helvetica, Arial, sans-serif;
			font-weight: 300;
			background-color: #eee;
		}

		h1 {
			font-weight: 100;
		}

		body {
			margin: 3em;
		}

		.logo {
			display: block;
			margin: 10px auto;
			width: 100px;
			height: 100px;
		}
	</style>
	<link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/highlight.js/9.15.6/styles/default.min.css">
</head>
<body>
	<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 50 50" class="logo">
		<defs>
			<style>
				.background {
					fill: #ccc;
				}

				.foreground {
					fill: #eee;
				}
			</style>
		</defs>
		<rect class="background" width="50" height="50"/>
		<rect class="foreground" x="10.73" y="10.72" width="28.55" height="11.35"/>
		<rect class="foreground" x="10.73" y="35.42" width="28.55" height="3.86"/>
		<rect class="foreground" x="10.73" y="25.74" width="28.55" height="5.82"/>
	</svg>
	<div class="markdown">%s</div>
	<script src="//cdnjs.cloudflare.com/ajax/libs/highlight.js/9.15.6/highlight.min.js"></script>
	<script>hljs.initHighlighting();</script>
</body>
</html>
`
)

type pretty struct {
	Title  string
	Format string
	Data   interface{}
}

func newPretty(c *gin.Context, code int, data interface{}) pretty {
	format := c.NegotiateFormat(
		binding.MIMEJSON,
		binding.MIMEYAML,
	)

	return pretty{
		Title:  http.StatusText(code),
		Format: format,
		Data:   data,
	}
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

func (p pretty) Render(w http.ResponseWriter) error {
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
		return err
	}

	md := fmt.Sprintf(prettyMarkdownTemplate, p.Title, format, string(formatted))
	html := blackfriday.Run(
		[]byte(md),
		blackfriday.WithExtensions(
			blackfriday.CommonExtensions|blackfriday.AutoHeadingIDs|blackfriday.HardLineBreak,
		),
	)
	_, err = fmt.Fprintf(w, prettyHTMLTemplate, string(html))
	return err
}

func (p pretty) WriteContentType(w http.ResponseWriter) {
	if v, ok := w.Header()["Content-Type"]; ok && len(v) > 0 {
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}
