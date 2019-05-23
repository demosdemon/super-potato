package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const UserCacheKey = "super-potato/pkg/server/CertifiedUser"

func getUser(c *gin.Context) User {
	if u, ok := c.Get(UserCacheKey); ok {
		if v, ok := u.(User); ok {
			return v
		}
	}
	return TheAnonymousUser
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
	r.GET("headers", s.getHeaders)
	s.registerGeneratedRoutes(r.Group("env", s.requireAuth))
}

func (s *Server) certifiedUserMiddleware(c *gin.Context) {
	defer c.Next()

	if getUser(c) != TheAnonymousUser {
		return
	}

	xClientCert := c.GetHeader("X-Client-Cert")
	if xClientCert == "" {
		xClientCert, _ = s.XClientCert()
		if xClientCert == "" {
			return
		}
	}

	var user CertifiedUser
	if err := user.ClientCertificate.UnmarshalText([]byte(xClientCert)); err != nil {
		logrus.WithError(err).WithField("xClientCert", xClientCert).Error("unable to decode certificate")
		return
	}

	user.DistinguishedName = c.GetHeader("X-Client-Dn")
	if user.DistinguishedName == "" {
		user.DistinguishedName, _ = s.XClientDN()
		if user.DistinguishedName == "" {
			user.DistinguishedName = user.ClientCertificate.Subject.String()
		}
	}

	c.Set(UserCacheKey, &user)
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

func (s *Server) cacheControl(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=604800")
	c.Next()
}
