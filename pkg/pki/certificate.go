package pki

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"encoding/xml"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
)

const timeFormat = "Jan 02 15:04:05 2006 MST"

type Certificate struct {
	*x509.Certificate
}

func UnmarshalBlock(s string) (*pem.Block, error) {
	if strings.Index(s, "%20") > 0 {
		logrus.WithField("s", s).Trace("detected percent encoding")
		var err error
		s, err = url.PathUnescape(s)
		if err != nil {
			return nil, err
		}
	}

	pb, _ := pem.Decode([]byte(s))
	return pb, nil
}

func (c *Certificate) reset() {
	c.Certificate = nil
}

func (c *Certificate) Unmarshal(mc MarshaledCertificate) error {
	if mc.Raw == "" {
		c.reset()
		return nil
	}

	pemBlock, err := UnmarshalBlock(mc.Raw)
	if err != nil {
		return err
	}

	cert, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		return err
	}

	c.Certificate = cert
	return nil
}

func (c Certificate) Marshal() MarshaledCertificate {
	rv := MarshaledCertificate{
		Subject:        c.Subject,
		Issuer:         c.Issuer,
		SerialNumber:   SerialNumber{c.SerialNumber},
		NotBefore:      c.NotBefore,
		NotAfter:       c.NotAfter,
		KeyUsage:       KeyUsage(c.KeyUsage),
		ExtKeyUsage:    ExtKeyUsage(c.ExtKeyUsage),
		SubjectKeyID:   NewSerialNumber(c.SubjectKeyId),
		AuthorityKeyID: NewSerialNumber(c.AuthorityKeyId),
		EmailAddresses: c.EmailAddresses,
	}

	raw := pem.EncodeToMemory(
		&pem.Block{
			Type:    "CERTIFICATE",
			Bytes:   c.Raw,
			Headers: rv.StringMap(),
		},
	)

	rv.Raw = url.PathEscape(string(raw))
	return rv
}

func (c *Certificate) UnmarshalText(text []byte) error {
	if text == nil || len(text) == 0 {
		c.reset()
		return nil
	}

	mc := MarshaledCertificate{Raw: string(text)}
	return c.Unmarshal(mc)
}

func (c Certificate) MarshalText() ([]byte, error) {
	if c.Certificate == nil {
		return nil, nil
	}

	mc := c.Marshal()
	raw, _ := url.PathUnescape(mc.Raw)
	return []byte(raw), nil
}

func (c *Certificate) UnmarshalJSON(data []byte) error {
	var mc MarshaledCertificate
	err := json.Unmarshal(data, &mc)
	if err != nil {
		return err
	}
	return c.Unmarshal(mc)
}

func (c Certificate) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Marshal())
}

func (c *Certificate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var mc MarshaledCertificate
	err := d.DecodeElement(&mc, &start)
	if err != nil {
		return err
	}
	return c.Unmarshal(mc)
}

func (c Certificate) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(c.Marshal(), start)
}

func (c *Certificate) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var mc MarshaledCertificate
	err := unmarshal(&mc)
	if err != nil {
		return err
	}
	return c.Unmarshal(mc)
}

func (c Certificate) MarshalYAML() (interface{}, error) {
	return c.Marshal(), nil
}
