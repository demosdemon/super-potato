package pki

import (
	"crypto/x509/pkix"
	"strings"
	"time"
)

type MarshaledCertificate struct {
	Raw            string `xml:"raw,attr"`
	Subject        pkix.Name
	Issuer         pkix.Name
	SerialNumber   SerialNumber
	NotBefore      time.Time
	NotAfter       time.Time
	KeyUsage       KeyUsage
	ExtKeyUsage    ExtKeyUsage
	SubjectKeyID   SerialNumber
	AuthorityKeyID SerialNumber
	EmailAddresses []string
}

func (c MarshaledCertificate) StringMap() map[string]string {
	rv := make(map[string]string, 10)
	rv["Subject"] = c.Subject.String()
	rv["Issuer"] = c.Issuer.String()
	rv["SerialNumber"] = c.SerialNumber.String()
	rv["NotBefore"] = c.NotBefore.Format(timeFormat)
	rv["NotAfter"] = c.NotAfter.Format(timeFormat)
	rv["KeyUsage"] = c.KeyUsage.String()
	rv["ExtKeyUsage"] = c.ExtKeyUsage.String()
	rv["SubjectKeyID"] = c.SubjectKeyID.String()
	rv["AuthorityKeyID"] = c.AuthorityKeyID.String()
	rv["EmailAddresses"] = strings.Join(c.EmailAddresses, ", ")
	return rv
}
