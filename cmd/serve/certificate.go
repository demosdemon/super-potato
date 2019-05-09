package serve

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"encoding/xml"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/demosdemon/super-potato/pkg/platformsh"
)

type Certificate struct {
	*x509.Certificate
	Rest []byte
}

func (c *Certificate) reset() {
	*c = Certificate{}
}

func (c *Certificate) UnmarshalText(text []byte) error {
	if text == nil || len(text) == 0 {
		c.reset()
		return nil
	}

	certData, err := url.PathUnescape(string(text))
	if err != nil {
		certData = string(text)
	}

	pemBlock, rest := pem.Decode([]byte(certData))
	if pemBlock == nil {
		return errors.New("invalid PEM data")
	}

	if pemBlock.Type != "CERTIFICATE" {
		return fmt.Errorf("expected a CERTIFICATE, got %s", pemBlock.Type)
	}

	cert, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		return err
	}

	c.Certificate = cert
	c.Rest = rest
	return nil
}

func (c Certificate) MarshalText() ([]byte, error) {
	if c.Certificate == nil {
		return nil, nil
	}

	pemBlock := pem.Block{
		Type:  "CERTIFICATE",
		Bytes: c.Raw,
	}
	text := pem.EncodeToMemory(&pemBlock)
	text = append(text, c.Rest...)
	return []byte(url.PathEscape(string(text))), nil
}

func (c *Certificate) UnmarshalJSON(data []byte) error {
	var obj platformsh.JSONObject
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	if raw, ok := obj["raw"]; ok {
		return c.UnmarshalText([]byte(raw.(string)))
	}

	return InvalidCertificate{obj}
}

func (c Certificate) MarshalJSON() ([]byte, error) {
	v := c.marshal()
	return json.Marshal(v)
}

func (c *Certificate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var obj platformsh.JSONObject
	if err := d.DecodeElement(&obj, &start); err != nil {
		return err
	}

	if raw, ok := obj["raw"]; ok {
		return c.UnmarshalText([]byte(raw.(string)))
	}

	return InvalidCertificate{obj}
}

func (c *Certificate) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	v := c.marshal()
	return e.EncodeElement(v, start)
}

func (c Certificate) marshal() interface{} {
	raw, _ := c.MarshalText()
	return gin.H{
		"raw":              string(raw),
		"serial_number":    c.SerialNumber.String(),
		"issuer":           c.Issuer.String(),
		"not_before":       c.NotBefore,
		"not_after":        c.NotAfter,
		"subject":          c.Subject.String(),
		"key_usage":        keyUsageString(c.KeyUsage),
		"ext_key_usage":    extKeyUsageString(c.ExtKeyUsage),
		"subject_key_id":   encodeBytes(c.SubjectKeyId),
		"authority_key_id": encodeBytes(c.AuthorityKeyId),
		"email":            c.EmailAddresses,
	}
}

type InvalidCertificate struct {
	Input platformsh.JSONObject
}

func (InvalidCertificate) Error() string {
	return "invalid certificate data"
}

func (c InvalidCertificate) String() string {
	return fmt.Sprintf("invalid certificate data: %#v", c.Input)
}

var keyUsageMap = map[x509.KeyUsage]string{
	x509.KeyUsageDigitalSignature:  "digital signature",
	x509.KeyUsageContentCommitment: "content commitment",
	x509.KeyUsageKeyEncipherment:   "key encipherment",
	x509.KeyUsageDataEncipherment:  "data encipherment",
	x509.KeyUsageKeyAgreement:      "key agreement",
	x509.KeyUsageCertSign:          "cert sign",
	x509.KeyUsageCRLSign:           "CRL sign",
	x509.KeyUsageEncipherOnly:      "encipher only",
	x509.KeyUsageDecipherOnly:      "decipher only",
}

func keyUsageString(value x509.KeyUsage) string {
	usage := make([]string, 0, len(keyUsageMap))
	for k, v := range keyUsageMap {
		if value&k != 0 {
			usage = append(usage, v)
		}
	}

	return strings.Join(usage, ", ")
}

var extKeyUsageMap = map[x509.ExtKeyUsage]string{
	x509.ExtKeyUsageAny:                            "any",
	x509.ExtKeyUsageServerAuth:                     "serve auth",
	x509.ExtKeyUsageClientAuth:                     "client auth",
	x509.ExtKeyUsageCodeSigning:                    "code signing",
	x509.ExtKeyUsageEmailProtection:                "email protection",
	x509.ExtKeyUsageIPSECEndSystem:                 "IPSEC end system",
	x509.ExtKeyUsageIPSECTunnel:                    "IPSEC tunnel",
	x509.ExtKeyUsageIPSECUser:                      "IPSEC user",
	x509.ExtKeyUsageTimeStamping:                   "time stamping",
	x509.ExtKeyUsageOCSPSigning:                    "OCSP signing",
	x509.ExtKeyUsageMicrosoftServerGatedCrypto:     "Microsoft serve gated crypto",
	x509.ExtKeyUsageNetscapeServerGatedCrypto:      "Netscape serve gated crypto",
	x509.ExtKeyUsageMicrosoftCommercialCodeSigning: "Microsoft commercial code signing",
	x509.ExtKeyUsageMicrosoftKernelCodeSigning:     "Microsoft kernel code signing",
}

func extKeyUsageString(value []x509.ExtKeyUsage) string {
	usage := make([]string, len(value))
	for idx, v := range value {
		if name, ok := extKeyUsageMap[v]; ok {
			usage[idx] = name
		} else {
			usage[idx] = fmt.Sprintf("%02x", int(v))
		}
	}
	return strings.Join(usage, ", ")
}

func encodeBytes(b []byte) string {
	rv := make([]string, len(b))
	for idx, v := range b {
		rv[idx] = fmt.Sprintf("%02x", v)
	}
	return strings.Join(rv, ":")
}
