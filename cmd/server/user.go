package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type User struct {
	ClientCertificate Certificate
	DistinguishedName string
}

const UserCacheKey = "github.com/demosdemon/super-potato/cmd/server/CertifiedUser"

func getCertifiedUser(c *gin.Context) *User {
	if v, ok := c.Get(UserCacheKey); ok {
		if v, ok := v.(*User); ok {
			return v
		}
	}

	xClientCert := c.GetHeader("X-Client-Cert")
	if xClientCert == "" {
		xClientCert, _ = env.XClientCert()
		if xClientCert == "" {
			return nil
		}
	}

	var user User
	if err := user.ClientCertificate.UnmarshalText([]byte(xClientCert)); err != nil {
		logrus.WithFields(logrus.Fields{
			"xClientCert": xClientCert,
			"error":       err,
		}).Error("unable to decode certificate")
		return nil
	}

	user.DistinguishedName = c.GetHeader("X-Client-Dn")
	if user.DistinguishedName == "" {
		user.DistinguishedName, _ = env.XClientDN()
		if user.DistinguishedName == "" {
			user.DistinguishedName = user.ClientCertificate.Subject.String()
		}
	}

	c.Set(UserCacheKey, &user)
	return &user
}
