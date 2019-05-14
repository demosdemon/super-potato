package server

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"expvar"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/mongo"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/russross/blackfriday/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/demosdemon/super-potato/pkg/app"
	"github.com/demosdemon/super-potato/pkg/platformsh"
)

const (
	OneYear  = 60 * 60 * 24 * 365
	OneMonth = 60 * 60 * 24 * 30
)

type Server struct {
	afero.Fs
	platformsh.Environment
	*gin.Engine

	registerOnce sync.Once

	start time.Time
}

func GetSecret(env platformsh.Environment) []byte {
	entropy, err := env.ProjectEntropy()
	if err == nil {
		if rv, err := base32.StdEncoding.DecodeString(entropy); err == nil {
			return rv
		} else {
			logrus.WithField("err", err).WithField("entropy", entropy).Warn("error decoding project entropy")
		}
	} else {
		logrus.WithField("err", err).Warn("project entropy not found")
	}

	secret := make([]byte, 40)
	_, _ = rand.Read(secret)
	logrus.WithField("secret", base32.StdEncoding.EncodeToString(secret)).Warn("using random secret")
	return secret
}

func mongoURL(rel platformsh.Relationship) string {
	var b strings.Builder
	b.WriteString("mongodb://")
	//if rel.Username != "" {
	//	b.WriteString(rel.Username)
	//	if rel.Password != "" {
	//		b.WriteString(":")
	//		b.WriteString(rel.Password)
	//	}
	//	b.WriteString("@")
	//}
	b.WriteString(rel.Host)
	if rel.Port > 0 {
		fmt.Fprintf(&b, ":%d", rel.Port)
	}
	if rel.Path != "" {
		b.WriteString("/")
		b.WriteString(rel.Path)
	}
	return b.String()
}

func GetSessionStore(env platformsh.Environment) sessions.Store {
	secret := GetSecret(env)
	var store sessions.Store
	if rels, err := env.Relationships(); err == nil {
		if rels, ok := rels["sessions"]; ok && len(rels) > 0 {
			rel := rels[0]
			sess, err := mgo.Dial(mongoURL(rel))
			if err != nil {
				logrus.WithField("err", err).Panic("failed to connect to mongo")
			}
			col := sess.DB("").C("sessions")
			store = mongo.NewStore(col, OneYear, true, secret)
		} else {
			logrus.Warn("unable to locate `sessions` relationship")
		}
	} else {
		logrus.WithField("err", err).Warn("unable to determine relationships")
	}
	if store == nil {
		store = cookie.NewStore(secret)
	}
	store.Options(sessions.Options{
		MaxAge: OneMonth,
		Secure: true,
	})
	return store
}

func New(app *app.App) *Server {
	wd, _ := os.Getwd()
	fs := afero.NewBasePathFs(app, wd)

	env := platformsh.NewEnvironment("PLATFORM_")
	store := GetSessionStore(env)

	engine := gin.New()
	engine.Use(
		gin.Logger(),
		gin.Recovery(),
		sessions.Sessions("super-potato", store),
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
	s.Use(s.sessionDuration)

	s.GET("/", s.root)
	s.GET("/ping", s.getPing)
	s.GET("/user", s.getUser)
	s.GET("/debug/vars", s.requireAuth, s.getDebugVars)
	s.registerGeneratedRoutes(s.Group("/env", s.requireAuth))
	s.GET("/favicon.ico", s.serverLifetime, s.getFaviconICO)
	s.GET("/logo.svg", s.getLogoSVG)
	s.GET("/logo.png", s.serverLifetime, s.getLogoPNG)
}

func (s *Server) sessionDuration(c *gin.Context) {
	session := sessions.Default(c)

	ts, _ := session.Get("ts").(string)
	if ts == "" {
		ts = time.Now().Format(time.RFC3339Nano)
		session.Set("ts", ts)
	}
	start, _ := time.Parse(time.RFC3339Nano, ts)

	count, _ := session.Get("count").(int)
	count += 1
	session.Set("count", count)
	_ = session.Save()

	elapsed := time.Now().Sub(start)
	if elapsed >= time.Minute {
		c.Header("X-Session-Duration", fmt.Sprintf("%v", elapsed))
	}

	c.Header("X-Session-Count", fmt.Sprintf("%d", count))
	c.Header("X-Session-Start", start.Format(time.RFC1123))
	c.Next()
}

func (s *Server) requireAuth(c *gin.Context) {
	if !getUser(c).Authenticated() {
		s.negotiate(c, http.StatusUnauthorized, gin.H{
			"message": "not logged in",
			"headers": Header{c.Request.Header},
		})
		c.Abort()
	}
	c.Next()
}

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

	data, err := ioutil.ReadAll(fp)
	fp.Close()

	if err != nil {
		_ = c.AbortWithError(500, err)
		return
	}

	output := blackfriday.Run(data)
	c.Header("Content-Type", "text/html")
	c.Writer.WriteHeader(200)
	c.Writer.Write(output)
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
	logo.Size = 100
	render := platformsh.RenderRasterLogo{
		RasterLogo:  logo,
		ContentType: platformsh.FormatPNG,
	}
	c.Render(http.StatusOK, render)
}

func (s *Server) serverLifetime(c *gin.Context) {
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

func (s *Server) cacheControl(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=604800")
	c.Next()
}

func (s *Server) etag(c *gin.Context) {
	c.Next()

}
