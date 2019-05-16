package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/demosdemon/super-potato/pkg/pki"
)

type User interface {
	UserName() string
	EmailAddress() string
	Authenticated() bool
}

type AnonymousUser struct{}

func (AnonymousUser) UserName() string {
	return "Anonymous"
}

func (AnonymousUser) EmailAddress() string {
	return "anonymous@invalid.domain"
}

func (AnonymousUser) Authenticated() bool {
	return false
}

var TheAnonymousUser = AnonymousUser{}

type CertifiedUser struct {
	ClientCertificate pki.Certificate
	DistinguishedName string
}

func (u *CertifiedUser) UserName() string {
	return u.DistinguishedName
}

func (u *CertifiedUser) EmailAddress() string {
	return u.ClientCertificate.EmailAddresses[0]
}

func (u *CertifiedUser) Authenticated() bool {
	// TODO: validate ClientCertificate is signed by an approved issuer
	return true
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
		logrus.WithFields(logrus.Fields{
			"xClientCert": xClientCert,
			"error":       err,
		}).Error("unable to decode certificate")
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

const UserCacheKey = "super-potato/pkg/server/CertifiedUser"

func getUser(c *gin.Context) User {
	if u, ok := c.Get(UserCacheKey); ok {
		if v, ok := u.(User); ok {
			return v
		}
	}
	return TheAnonymousUser
}
