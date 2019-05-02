package platformsh

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	TLSv10 TLSVersion = iota + 0x0301
	TLSv11
	TLSv12
	TLSv13
	TLSv14
)

var tlsNameMapping = map[TLSVersion]string{
	TLSv10: "TLSv1.0",
	TLSv11: "TLSv1.1",
	TLSv12: "TLSv1.2",
	TLSv13: "TLSv1.3",
	TLSv14: "TLSv1.4",
}

type (
	Duration struct {
		time.Duration
	}

	RouteIdentification struct {
		Scheme string `json:"scheme"`
		Host   string `json:"host"`
		Path   string `json:"path"`
	}

	Cache struct {
		Enabled    bool     `json:"enabled"`
		DefaultTTL int      `json:"default_ttl"`
		Cookies    []string `json:"cookies"`
		Headers    []string `json:"headers"`
	}

	SSI struct {
		Enabled bool `json:"enabled"`
	}

	RedirectPath struct {
		Regexp       bool     `json:"regexp"`
		To           string   `json:"to"`
		Prefix       bool     `json:"prefix"`
		AppendSuffix bool     `json:"append_suffix"`
		Code         int      `json:"code"`
		Expires      Duration `json:"expires"`
	}

	RedirectPaths map[string]RedirectPath

	Redirects struct {
		Expires Duration      `json:"expires"`
		Paths   RedirectPaths `json:"paths"`
	}

	TLSSTS struct {
		Enabled           bool `json:"enabled"`
		IncludeSubdomains bool `json:"include_subdomains"`
		Preload           bool `json:"preload"`
	}

	TLSVersion uint16

	TLSSettings struct {
		StrictTransportSecurity      TLSSTS                       `json:"strict_transport_security"`
		MinVersion                   *TLSVersion                  `json:"min_version"`
		ClientAuthentication         string                       `json:"client_authentication"`
		ClientCertificateAuthorities []ClientCertificateAuthority `json:"client_certificate_authorities"`
	}

	ClientCertificateAuthority struct {
		*x509.Certificate
	}

	HTTPAccess struct {
		Addresses []string          `json:"addresses"`
		BasicAuth map[string]string `json:"basic_auth"`
	}

	Route struct {
		Primary        bool              `json:"primary"`
		ID             *string           `json:"id"`
		OriginalURL    string            `json:"original_url"`
		Attributes     map[string]string `json:"attributes"`
		Type           string            `json:"type"`
		Redirects      Redirects         `json:"redirects"`
		TLS            TLSSettings       `json:"tls"`
		HTTPAccess     HTTPAccess        `json:"http_access"`
		RestrictRobots bool              `json:"restrict_robots"`

		// Upstream Routes
		Cache    Cache  `json:"cache"`
		SSI      SSI    `json:"ssi"`
		Upstream string `json:"upstream"`

		// Redirect Routes
		To string `json:"to"`
	}

	Routes map[url.URL]Route

	RouteRepresentation struct {
		Project     string              `json:"project"`
		Environment string              `json:"environment"`
		Route       RouteIdentification `json:"route"`
	}
)

func (v *Duration) UnmarshalText(text []byte) (err error) {
	logrus.Trace("Duration.UnmarshalText")
	v.Duration, err = time.ParseDuration(string(text))
	return err
}

func (v Duration) MarshalText() ([]byte, error) {
	logrus.Trace("Duration.MarshalText")
	return []byte(v.String()), nil
}

func NewTLSVersion(name string) (TLSVersion, error) {
	logrus.Trace("NewTLSVersion")
	for k, v := range tlsNameMapping {
		if v == name {
			return k, nil
		}
	}

	return 0, fmt.Errorf("unknown TLSVersion %q", name)
}

func (v TLSVersion) String() string {
	logrus.Trace("TLSVersion.String")
	if name, ok := tlsNameMapping[v]; ok {
		return name
	}

	return fmt.Sprintf("unknown TLSVersion 0x%04x", uint16(v))
}

func (v *TLSVersion) UnmarshalText(text []byte) (err error) {
	logrus.Trace("TLSVersion.UnmarshalText")
	*v, err = NewTLSVersion(string(text))
	return err
}

func (v TLSVersion) MarshalText() ([]byte, error) {
	logrus.Trace("TLSVersion.MarshalText")
	if rv, ok := tlsNameMapping[v]; ok {
		return []byte(rv), nil
	}

	return nil, errors.New(v.String())
}

func (v *ClientCertificateAuthority) UnmarshalText(text []byte) error {
	logrus.Trace("ClientCertificateAuthority.UnmarshalText")
	block, rest := pem.Decode(text)
	if block == nil {
		return errors.New("invalid PEM data")
	}
	if rest != nil && len(rest) > 0 {
		return errors.New("excess data after decoding the PEM block")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}

	*v = ClientCertificateAuthority{Certificate: cert}
	return nil
}

func (v ClientCertificateAuthority) MarshalText() ([]byte, error) {
	logrus.Trace("ClientCertificateAuthority.MarshalText")
	var block = pem.Block{
		Type:  "CERTIFICATE",
		Bytes: v.Raw,
	}
	data := pem.EncodeToMemory(&block)
	return data, nil
}

func (r *Routes) UnmarshalJSON(text []byte) error {
	logrus.Trace("Routes.UnmarshalJSON")
	var intermediate map[string]Route
	err := json.Unmarshal(text, &intermediate)
	if err != nil {
		return err
	}

	*r = make(Routes, len(intermediate))
	for k, v := range intermediate {
		kURL, err := url.Parse(k)
		if err != nil {
			return err
		}

		(*r)[*kURL] = v
	}

	return nil
}

func (r Routes) MarshalJSON() ([]byte, error) {
	logrus.Trace("Routes.MarshalJSON")
	intermediate := make(map[string]Route, len(r))

	for kURL, v := range r {
		intermediate[kURL.String()] = v
	}

	return json.Marshal(intermediate)
}
