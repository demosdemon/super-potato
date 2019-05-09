package serve

import (
	"bytes"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/render"
)

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
