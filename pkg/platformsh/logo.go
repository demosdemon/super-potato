package platformsh

import (
	"bytes"
	"encoding/xml"
	"html/template"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
	"github.com/llgcode/draw2d/draw2dsvg"
	"github.com/pkg/errors"
	"golang.org/x/image/tiff"
	"gopkg.in/go-playground/colors.v1"

	"github.com/demosdemon/super-potato/pkg/favicon"
)

const (
	FormatSVG  = "image/svg+xml"
	FormatGIF  = "image/gif"
	FormatICO  = "image/vnd.microsoft.icon"
	FormatJPG  = "image/jpeg"
	FormatPNG  = "image/png"
	FormatTIFF = "image/tiff"
)

const logoSVGTemplate = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 50 50">
	<defs>
		<style>
			.background {
				fill: {{ .Background }};
			}

			.foreground {
				fill: {{ .Foreground }};
			}
		</style>
	</defs>
	<rect class="background" width="50" height="50"/>
	<rect class="foreground" x="10.73" y="10.72" width="28.55" height="11.35"/>
	<rect class="foreground" x="10.73" y="25.74" width="28.55" height="5.82"/>
	<rect class="foreground" x="10.73" y="35.42" width="28.55" height="3.86"/>
</svg>
`

type SVGLogo struct {
	Background string
	Foreground string
}

func NewLogoSVG() SVGLogo {
	return SVGLogo{
		Background: "#0a0a0a",
		Foreground: "#fff",
	}
}

func (x SVGLogo) Execute(w io.Writer) error {
	tpl := template.Must(template.New("").Parse(logoSVGTemplate))
	return tpl.Execute(w, x)
}

func (x SVGLogo) Render(w http.ResponseWriter) error {
	x.WriteContentType(w)
	var buf bytes.Buffer
	err := x.Execute(&buf)
	if err != nil {
		return err
	}

	_, _ = buf.WriteTo(w)
	return nil
}

func (x SVGLogo) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", FormatSVG)
}

type RasterLogo struct {
	Size       int
	Background string
	Foreground string
}

func NewRasterLogo() RasterLogo {
	return RasterLogo{
		Size:       50,
		Background: "#0a0a0a",
		Foreground: "#fff",
	}
}

type RenderRasterLogo struct {
	RasterLogo
	ContentType string
}

func (x RasterLogo) draw(ctx draw2d.GraphicContext) error {
	width, height := float64(x.Size), float64(x.Size)

	bg, err := x.bg()
	if err != nil {
		return errors.Wrap(err, "invalid Background")
	}

	fg, err := x.fg()
	if err != nil {
		return errors.Wrap(err, "invalid Foreground")
	}

	ctx.SetFillColor(bg)
	rect(ctx, 0, 0, width, height)
	ctx.Fill()

	spanRoot := width * 10.73 / 50.0
	spanWidth := width * 28.55 / 50.0

	for _, span := range []struct {
		y      float64
		height float64
	}{
		{height * 10.72 / 50, height * 11.35 / 50},
		{height * 25.74 / 50, height * 5.82 / 50},
		{height * 35.42 / 50, height * 3.86 / 50},
	} {
		ctx.SetFillColor(fg)
		rect(ctx, spanRoot, span.y, spanWidth, span.height)
		ctx.FillStroke()
	}

	return nil
}

func (x RasterLogo) bg() (color.Color, error) {
	return parseColor(x.Background)
}

func (x RasterLogo) fg() (color.Color, error) {
	return parseColor(x.Foreground)
}

func (x RasterLogo) Negotiate(c *gin.Context) render.Render {
	format := c.NegotiateFormat(
		//FormatSVG,
		FormatICO,
		FormatPNG,
		FormatJPG,
		FormatGIF,
		FormatTIFF,
	)

	var rv RenderRasterLogo
	rv.RasterLogo = x
	rv.ContentType = format

	return RenderRasterLogo{RasterLogo: x, ContentType: format}
}

func (x RenderRasterLogo) Render(w http.ResponseWriter) error {
	x.WriteContentType(w)
	switch x.ContentType {
	case FormatSVG:
		svg := draw2dsvg.NewSvg()
		ctx := draw2dsvg.NewGraphicContext(svg)
		err := x.draw(ctx)
		if err != nil {
			return err
		}
		data, err := xml.Marshal(svg)
		if err != nil {
			return err
		}
		_, err = w.Write(data)
		if err != nil {
			return err
		}
	default:
		m := image.NewRGBA(image.Rect(0, 0, x.Size, x.Size))
		ctx := draw2dimg.NewGraphicContext(m)
		err := x.draw(ctx)
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		switch x.ContentType {
		case FormatGIF:
			err = gif.Encode(&buf, m, nil)
		case FormatICO:
			err = favicon.Encode(&buf, m)
		case FormatJPG:
			err = jpeg.Encode(&buf, m, nil)
		case FormatPNG:
			err = png.Encode(&buf, m)
		case FormatTIFF:
			err = tiff.Encode(&buf, m, nil)
		}

		if err != nil {
			return err
		}

		_, _ = buf.WriteTo(w)
	}

	return nil
}

func (x RenderRasterLogo) WriteContentType(w http.ResponseWriter) {
	if v, ok := w.Header()["Content-Type"]; ok && len(v) > 0 {
		return
	}
	w.Header().Set("Content-Type", x.ContentType)
}

type xColor struct {
	colors.Color
}

func (x xColor) RGBA() (uint32, uint32, uint32, uint32) {
	rgba := x.ToRGBA()
	return uint32(rgba.R) << 8, uint32(rgba.G) << 8, uint32(rgba.B) << 8, uint32(rgba.A*255) << 8
}

func wrapColor(x colors.Color) xColor {
	return xColor{Color: x}
}

func rect(path draw2d.PathBuilder, x, y, width, height float64) {
	x2 := x + width
	y2 := y + height
	draw2dkit.Rectangle(path, x, y, x2, y2)
}

func parseColor(s string) (color.Color, error) {
	parsed, err := colors.Parse(s)
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing color %q", s)
	}

	return wrapColor(parsed), nil
}
