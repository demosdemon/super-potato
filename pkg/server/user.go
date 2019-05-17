package server

import (
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
