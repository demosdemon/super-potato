package server

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/demosdemon/super-potato/pkg/app"
	"github.com/demosdemon/super-potato/pkg/platformsh"
)

type Server struct {
	afero.Fs
	platformsh.Environment
	*gin.Engine

	registerOnce sync.Once
}

func New(app *app.App) *Server {
	wd, _ := os.Getwd()
	fs := afero.NewBasePathFs(app, wd)

	env := platformsh.NewEnvironment("PLATFORM_")

	engine := gin.New()
	engine.Use(
		gin.Logger(),
		gin.Recovery(),
	)

	s := &Server{
		Fs:          fs,
		Environment: env,
		Engine:      engine,
	}

	s.SetFileSystem(s)
	return s
}

func (s *Server) Serve(l net.Listener) (err error) {
	if l == nil {
		l, err = s.Listener()
		if err != nil {
			return err
		}
	}

	return http.Serve(l, s.Register())
}

func (s *Server) Register() *Server {
	s.registerOnce.Do(s.register)
	return s
}

func (s *Server) register() {
	s.Use(s.certifiedUserMiddleware)
	s.GET("/", s.root)
	s.GET("/ping", s.getPing)
	s.GET("/user", s.getUser)
	s.registerGeneratedRoutes(s.Group("/env", func(c *gin.Context) {
		defer c.Next()
		if !getUser(c).Authenticated() {
			s.negotiate(c, http.StatusUnauthorized, gin.H{
				"message": "not logged in",
				"headers": Header{c.Request.Header},
			})
			c.Abort()
		}
	}))

	s.GET("/logo.svg", s.getLogoSVG)
}

func (s *Server) root(c *gin.Context) {
	fp, err := s.Open("/README.md")
	if err != nil {
		if os.IsNotExist(err) {
			_ = c.AbortWithError(404, err)
		} else {
			_ = c.AbortWithError(500, err)
		}
		return
	}

	data, err := ioutil.ReadAll(fp)
	_ = fp.Close()

	if err != nil {
		_ = c.AbortWithError(500, err)
		return
	}

	output := blackfriday.Run(data)
	c.Header("Content-Type", "text/html")
	c.Writer.WriteHeader(200)
	_, _ = c.Writer.Write(output)
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

func (s *Server) getLogoSVG(c *gin.Context) {
	c.Render(http.StatusOK, &logoSVG{
		Background: c.DefaultQuery("background", "#0a0a0a"),
		Foreground: c.DefaultQuery("foreground", "#fff"),
	})
}

const logoSVGTemplate = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 50 50">
<defs><style>.background {fill: {{ .Background }};}.foreground {fill: {{ .Foreground }};}</style>
</defs>
<rect class="background" width="50" height="50"/>
<rect class="foreground" x="10.73" y="10.72" width="28.55" height="11.35"/>
<rect class="foreground" x="10.73" y="35.42" width="28.55" height="3.86"/>
<rect class="foreground" x="10.73" y="25.74" width="28.55" height="5.82"/>
</svg>
`

type logoSVG struct {
	Background string
	Foreground string
}

func (x *logoSVG) Render(w http.ResponseWriter) error {
	x.WriteContentType(w)
	var buf bytes.Buffer
	tpl := template.Must(template.New("").Parse(logoSVGTemplate))
	err := tpl.Execute(&buf, x)
	if err != nil {
		return err
	}

	_, _ = buf.WriteTo(w)
	return nil
}

func (x *logoSVG) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "image/svg+xml")
}
