package server

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"expvar"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/cloudflare/cfssl/certdb"
	"github.com/cloudflare/cfssl/ocsp"
	"github.com/cloudflare/cfssl/signer"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/mongo"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/russross/blackfriday/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/demosdemon/super-potato/pkg/app"
	"github.com/demosdemon/super-potato/pkg/platformsh"
)

const (
	OneYear  = 60 * 60 * 24 * 365
	OneMonth = 60 * 60 * 24 * 30
)

type Server struct {
	*app.App      `flag:"-"`
	SessionCookie string `flag:"session-cookie" desc:"The name of the session cookie." env:"PKI_SESSION_COOKIE"`

	once       sync.Once
	start      time.Time
	engine     *gin.Engine
	accessor   certdb.Accessor
	signer     signer.Signer
	ocspSigner ocsp.Signer
}

func (s *Server) Use() string {
	return "serve"
}

func (s *Server) Args(cmd *cobra.Command, args []string) error {
	return cobra.NoArgs(cmd, args)
}

func (s *Server) Run(cmd *cobra.Command, args []string) error {
	return s.Serve()
}

func (s *Server) init() {
	s.start = time.Now().Truncate(time.Second)
	s.engine = gin.New()
	s.register(s.engine)
}

func (s *Server) Init() {
	s.once.Do(s.init)
}

func (s *Server) Serve() error {
	s.Init()

	l, err := s.Listener()
	if err != nil {
		return errors.Wrap(err, "unable to open listener")
	}
	defer l.Close()

	done := make(chan error)
	defer close(done)

	go func() {
		srv := http.Server{Handler: s.engine}
		go func() {
			done <- srv.Serve(l)
		}()

		<-s.Done()
		logrus.WithError(s.Err()).Debug("context done")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logrus.WithError(err).Warn("error shutting down server")
		}
	}()

	err = <-done
	if err != nil {
		logrus.WithError(err).Warning("server shutdown")
	}
	return err
}

func (s *Server) GetSecret() []byte {
	entropy, err := s.ProjectEntropy()
	if err == nil {
		logrus.WithField("entropy", entropy).Debug("found project entropy")
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

func (s *Server) GetSessionStore() sessions.Store {
	secret := s.GetSecret() // TODO: rotate secret?
	var store sessions.Store
	if rels, err := s.Relationships(); err == nil {
		db, err := rels.MongoDB("sessions")
		if err == nil {
			col := db.C("sessions")
			store = mongo.NewStore(col, OneYear, true, secret)
		} else {
			logrus.WithError(err).Warn("unable to connect to mongo server")
		}
	} else {
		logrus.WithField("err", err).Warn("unable to determine relationships")
	}
	if store == nil {
		logrus.Warn("using cookie session store")
		store = cookie.NewStore(secret)
	}
	store.Options(sessions.Options{
		MaxAge: OneMonth,
		Secure: true,
	})
	return store
}

func (s *Server) GetDB() (*sqlx.DB, error) {
	rels, err := s.Relationships()
	if err != nil {
		return nil, errors.Wrap(err, "unable to locate relationships")
	}

	dbOpen, err := rels.Postgresql("database")
	if err != nil {
		return nil, errors.Wrap(err, "unable to get database connection string")
	}

	db, err := sqlx.Open("postgres", dbOpen)
	if err != nil {
		return nil, errors.Wrap(err, "unable to connect to postgres")
	}

	return db, nil
}

func (s *Server) register(r gin.IRouter) {
	r.Use(
		// order is important
		gin.Logger(),
		gin.Recovery(),
		sessions.Sessions(s.SessionCookie, s.GetSessionStore()),
		s.certifiedUserMiddleware,
		s.sessionDuration,
	)

	r.GET("", s.root)
	r.GET("ping", s.getPing)
	r.GET("user", s.getUser)
	r.GET("debug/vars", s.requireAuth, s.getDebugVars)
	r.GET("favicon.ico", s.serverLifetime, s.getFaviconICO)
	r.GET("logo.svg", s.cacheControl, s.getLogoSVG)
	r.GET("logo.png", s.serverLifetime, s.getLogoPNG)
	s.registerGeneratedRoutes(r.Group("env", s.requireAuth))
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
	err := session.Save()
	if err != nil {
		logrus.WithError(err).Warn("unable to save session")
	}

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
