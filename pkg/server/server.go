package server

import (
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

	start time.Time
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
		start:       time.Now().Truncate(time.Second),
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

	s.GET("/favicon.ico", s.serverLifetimeMiddleware, s.getFaviconICO)
	s.GET("/logo.svg", s.getLogoSVG)
	s.GET("/logo.png", s.serverLifetimeMiddleware, s.getLogoPNG)
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
	logo := platformsh.NewLogoSVG()
	logo.Background = c.DefaultQuery("background", logo.Background)
	logo.Foreground = c.DefaultQuery("foreground", logo.Foreground)
	c.Render(http.StatusOK, logo)
}

func (s *Server) getFaviconICO(c *gin.Context) {
	logo := platformsh.NewRasterLogo()
	logo.Size = 32
	render := logo.Negotiate(c)
	c.Render(http.StatusOK, render)
}

func (s *Server) getLogoPNG(c *gin.Context) {
	logo := platformsh.NewRasterLogo()
	logo.Size = 50
	render := platformsh.RenderRasterLogo{
		RasterLogo:  logo,
		ContentType: platformsh.FormatPNG,
	}
	c.Render(http.StatusOK, render)
}

func (s *Server) serverLifetimeMiddleware(c *gin.Context) {
	switch c.Request.Method {
	case "GET", "HEAD":
		// hurray!
	default:
		c.AbortWithStatus(http.StatusMethodNotAllowed)
		return
	}

	ifModifiedSince := c.GetHeader("If-Modified-Since")
	if ifModifiedSince != "" {
		if parsed, err := time.Parse(time.RFC1123, ifModifiedSince); err == nil {
			if !parsed.Before(s.start) {
				c.Header("X-Cache", "HIT")
				c.AbortWithStatus(http.StatusNotModified)
			} else {
				c.Header("X-Cache", "MISS")
			}
		} else {
			logrus.WithField("If-Modified-Since", ifModifiedSince).Warn("invalid If-Modified-Since header")
		}
	}

	c.Header("Last-Modified", s.start.Format(time.RFC1123))
	c.Next()
}
