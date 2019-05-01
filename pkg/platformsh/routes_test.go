package platformsh_test

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/demosdemon/super-potato/pkg/platformsh"
)

func init() {
	const versionTests = 6
	size := versionTests*len(versions) + len(certificates) + len(jsonTests)
	testCases = make(tests, 0, size)

	for _, v := range versions {
		testCases = testCases.Append(
			newTLSVersionTest{v},
			tlsVersionStringTest{v},
			tlsVersionUnmarshalText{v},
			tlsVersionMarshalText{v},
			tlsVersionMarshalJSON{v},
			tlsVersionUnmarshalJSON{v},
		)
	}

	for _, v := range certificates {
		testCases = testCases.Append(v)
	}

	for _, v := range jsonTests {
		testCases = testCases.Append(v)
	}
}

type (
	tester interface {
		Name() string
		Run(*assert.Assertions)
	}

	tests []tester

	tlsVersionTest struct {
		name   string
		value  TLSVersion
		string string
	}

	newTLSVersionTest       struct{ tlsVersionTest }
	tlsVersionStringTest    struct{ tlsVersionTest }
	tlsVersionUnmarshalText struct{ tlsVersionTest }
	tlsVersionMarshalText   struct{ tlsVersionTest }
	tlsVersionUnmarshalJSON struct{ tlsVersionTest }
	tlsVersionMarshalJSON   struct{ tlsVersionTest }

	clientCertificateAuthorityTest struct {
		name    string
		pemData string
		valid   bool
	}

	jsonTest struct {
		name      string
		input     string
		expected  interface{}
		unmarshal func(text string) (interface{}, error)
	}
)

var (
	versions = []tlsVersionTest{
		{"zero", 0, "unknown TLSVersion 0x0000"},
		{"v0", TLSv10, "TLSv1.0"},
		{"v1", TLSv11, "TLSv1.1"},
		{"v2", TLSv12, "TLSv1.2"},
		{"v3", TLSv13, "TLSv1.3"},
		{"v4", TLSv14, "TLSv1.4"},
		{"v5", TLSv14 + 1, "unknown TLSVersion 0x0306"},
	}

	certificates = []clientCertificateAuthorityTest{
		{
			name:    "zero",
			pemData: "",
			valid:   false,
		},
		{
			name: "GTE CyberTrust Global Root",
			pemData: `-----BEGIN CERTIFICATE-----
MIICWjCCAcMCAgGlMA0GCSqGSIb3DQEBBAUAMHUxCzAJBgNVBAYTAlVTMRgwFgYD
VQQKEw9HVEUgQ29ycG9yYXRpb24xJzAlBgNVBAsTHkdURSBDeWJlclRydXN0IFNv
bHV0aW9ucywgSW5jLjEjMCEGA1UEAxMaR1RFIEN5YmVyVHJ1c3QgR2xvYmFsIFJv
b3QwHhcNOTgwODEzMDAyOTAwWhcNMTgwODEzMjM1OTAwWjB1MQswCQYDVQQGEwJV
UzEYMBYGA1UEChMPR1RFIENvcnBvcmF0aW9uMScwJQYDVQQLEx5HVEUgQ3liZXJU
cnVzdCBTb2x1dGlvbnMsIEluYy4xIzAhBgNVBAMTGkdURSBDeWJlclRydXN0IEds
b2JhbCBSb290MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCVD6C28FCc6HrH
iM3dFw4usJTQGz0O9pTAipTHBsiQl8i4ZBp6fmw8U+E3KHNgf7KXUwefU/ltWJTS
r41tiGeA5u2ylc9yMcqlHHK6XALnZELn+aks1joNrI1CqiQBOeacPwGFVw1Yh0X4
04Wqk2kmhXBIgD8SFcd5tB8FLztimQIDAQABMA0GCSqGSIb3DQEBBAUAA4GBAG3r
GwnpXtlR22ciYaQqPEh346B8pt5zohQDhT37qw4wxYMWM4ETCJ57NE7fQMh017l9
3PR2VX2bY1QY6fDq81yx2YtCHrnAlU66+tXifPVoYb+O7AWXX1uw16OFNMQkpw0P
lZPvy5TYnh+dXIVtx6quTx8itc2VrbqnzPmrC3p/
-----END CERTIFICATE-----
`,
			valid: true,
		},
		{
			name: "excess data",
			pemData: `-----BEGIN CERTIFICATE-----
MIICWjCCAcMCAgGlMA0GCSqGSIb3DQEBBAUAMHUxCzAJBgNVBAYTAlVTMRgwFgYD
VQQKEw9HVEUgQ29ycG9yYXRpb24xJzAlBgNVBAsTHkdURSBDeWJlclRydXN0IFNv
bHV0aW9ucywgSW5jLjEjMCEGA1UEAxMaR1RFIEN5YmVyVHJ1c3QgR2xvYmFsIFJv
b3QwHhcNOTgwODEzMDAyOTAwWhcNMTgwODEzMjM1OTAwWjB1MQswCQYDVQQGEwJV
UzEYMBYGA1UEChMPR1RFIENvcnBvcmF0aW9uMScwJQYDVQQLEx5HVEUgQ3liZXJU
cnVzdCBTb2x1dGlvbnMsIEluYy4xIzAhBgNVBAMTGkdURSBDeWJlclRydXN0IEds
b2JhbCBSb290MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCVD6C28FCc6HrH
iM3dFw4usJTQGz0O9pTAipTHBsiQl8i4ZBp6fmw8U+E3KHNgf7KXUwefU/ltWJTS
r41tiGeA5u2ylc9yMcqlHHK6XALnZELn+aks1joNrI1CqiQBOeacPwGFVw1Yh0X4
04Wqk2kmhXBIgD8SFcd5tB8FLztimQIDAQABMA0GCSqGSIb3DQEBBAUAA4GBAG3r
GwnpXtlR22ciYaQqPEh346B8pt5zohQDhT37qw4wxYMWM4ETCJ57NE7fQMh017l9
3PR2VX2bY1QY6fDq81yx2YtCHrnAlU66+tXifPVoYb+O7AWXX1uw16OFNMQkpw0P
lZPvy5TYnh+dXIVtx6quTx8itc2VrbqnzPmrC3p/
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIDAjCCAmsCEH3Z/gfPqB63EHln+6eJNMYwDQYJKoZIhvcNAQEFBQAwgcExCzAJ
BgNVBAYTAlVTMRcwFQYDVQQKEw5WZXJpU2lnbiwgSW5jLjE8MDoGA1UECxMzQ2xh
c3MgMyBQdWJsaWMgUHJpbWFyeSBDZXJ0aWZpY2F0aW9uIEF1dGhvcml0eSAtIEcy
MTowOAYDVQQLEzEoYykgMTk5OCBWZXJpU2lnbiwgSW5jLiAtIEZvciBhdXRob3Jp
emVkIHVzZSBvbmx5MR8wHQYDVQQLExZWZXJpU2lnbiBUcnVzdCBOZXR3b3JrMB4X
DTk4MDUxODAwMDAwMFoXDTI4MDgwMTIzNTk1OVowgcExCzAJBgNVBAYTAlVTMRcw
FQYDVQQKEw5WZXJpU2lnbiwgSW5jLjE8MDoGA1UECxMzQ2xhc3MgMyBQdWJsaWMg
UHJpbWFyeSBDZXJ0aWZpY2F0aW9uIEF1dGhvcml0eSAtIEcyMTowOAYDVQQLEzEo
YykgMTk5OCBWZXJpU2lnbiwgSW5jLiAtIEZvciBhdXRob3JpemVkIHVzZSBvbmx5
MR8wHQYDVQQLExZWZXJpU2lnbiBUcnVzdCBOZXR3b3JrMIGfMA0GCSqGSIb3DQEB
AQUAA4GNADCBiQKBgQDMXtERXVxp0KvTuWpMmR9ZmDCOFoUgRm1HP9SFIIThbbP4
pO0M8RcPO/mn+SXXwc+EY/J8Y8+iR/LGWzOOZEAEaMGAuWQcRXfH2G71lSk8UOg0
13gfqLptQ5GVj0VXXn7F+8qkBOvqlzdUMG+7AUcyM83cV5tkaWH4mx0ciU9cZwID
AQABMA0GCSqGSIb3DQEBBQUAA4GBAFFNzb5cy5gZnBWyATl4Lk0PZ3BwmcYQWpSk
U01UbSuvDV1Ai2TT1+7eVmGSX6bEHRBhNtMsJzzoKQm5EWR0zLVznxxIqbxhAe7i
F6YM40AIOw7n60RzKprxaZLvcRTDOaxxp5EJb+RxBrO6WVcmeQD2+A2iMzAo1KpY
oJ2daZH9
-----END CERTIFICATE-----`,
			valid: false,
		},
		{
			name: "not a certificate",
			// intentionally weak private key... this is just a test
			pemData: `-----BEGIN RSA PRIVATE KEY-----
MC4CAQACBTQ6QL+/AgMBAAECBQpbhD15AgMH2IsCAwaoHQIDBAh7AgI0ZQIDBA0b
-----END RSA PRIVATE KEY-----`,
			valid: false,
		},
	}

	jsonTests = []jsonTest{
		{
			name:  "RoutesSchema",
			input: `{"https://master-7rqtwti-fteigbda5stns.eu-3.platformsh.site/": {"primary": true, "id": null, "attributes": {}, "type": "upstream", "tls": {"client_authentication": null, "min_version": null, "client_certificate_authorities": [], "strict_transport_security": {"preload": null, "include_subdomains": null, "enabled": null}}, "cache": {"default_ttl": 0, "cookies": ["*"], "enabled": true, "headers": ["Accept", "Accept-Language"]}, "ssi": {"enabled": false}, "upstream": "app", "original_url": "https://{default}/", "restrict_robots": true, "http_access": {"addresses": [], "basic_auth": {}}}, "http://master-7rqtwti-fteigbda5stns.eu-3.platformsh.site/": {"primary": false, "id": null, "attributes": {}, "type": "redirect", "tls": {"client_authentication": null, "min_version": null, "client_certificate_authorities": [], "strict_transport_security": {"preload": null, "include_subdomains": null, "enabled": null}}, "to": "https://master-7rqtwti-fteigbda5stns.eu-3.platformsh.site/", "original_url": "http://{default}/", "restrict_robots": true, "http_access": {"addresses": [], "basic_auth": {}}}}`,
			expected: RoutesSchema{
				mustURL("https://master-7rqtwti-fteigbda5stns.eu-3.platformsh.site/"): {
					Primary:    true,
					Attributes: make(map[string]string),
					Type:       "upstream",
					TLS: TLSSettingsSchema{
						ClientCertificateAuthorities: make([]ClientCertificateAuthority, 0),
					},
					Cache: CacheSchema{
						DefaultTTL: 0,
						Cookies:    []string{"*"},
						Enabled:    true,
						Headers:    []string{"Accept", "Accept-Language"},
					},
					SSI: SSISchema{
						Enabled: false,
					},
					Upstream:       "app",
					OriginalURL:    "https://{default}/",
					RestrictRobots: true,
					HTTPAccess: HTTPAccessSchema{
						Addresses: make([]string, 0),
						BasicAuth: make(map[string]string),
					},
				},
				mustURL("http://master-7rqtwti-fteigbda5stns.eu-3.platformsh.site/"): {
					Primary:    false,
					Attributes: make(map[string]string),
					Type:       "redirect",
					TLS: TLSSettingsSchema{
						ClientCertificateAuthorities: make([]ClientCertificateAuthority, 0),
					},
					To:             "https://master-7rqtwti-fteigbda5stns.eu-3.platformsh.site/",
					OriginalURL:    "http://{default}/",
					RestrictRobots: true,
					HTTPAccess: HTTPAccessSchema{
						Addresses: make([]string, 0),
						BasicAuth: make(map[string]string),
					},
				},
			},
			unmarshal: func(text string) (interface{}, error) {
				var rv = make(RoutesSchema)
				err := json.Unmarshal([]byte(text), &rv)
				return rv, err
			},
		},
	}

	testCases tests
)

func mustURL(text string) url.URL {
	rv, err := url.Parse(text)
	if err != nil {
		panic(err)
	}
	return *rv
}

func (x tlsVersionTest) valid() bool {
	return strings.HasPrefix(x.string, "TLSv")
}

func (x tests) Run(t *testing.T) {
	t.Parallel()
	for _, c := range x {
		c := c
		t.Run(c.Name(), func(t *testing.T) {
			c.Run(assert.New(t))
		})
	}
}

func (x tests) Append(c ...tester) tests {
	return append(x, c...)
}

func (t newTLSVersionTest) Name() string {
	return "NewTLSVersion/" + t.name
}

func (t newTLSVersionTest) Run(x *assert.Assertions) {
	output, err := NewTLSVersion(t.string)

	if t.valid() {
		x.NoError(err)
		x.Equal(t.value, output)
	} else {
		x.Error(err)
		x.Equal(TLSVersion(0), output)
	}
}

func (t tlsVersionStringTest) Name() string {
	return "TLSVersion_String/" + t.name
}

func (t tlsVersionStringTest) Run(x *assert.Assertions) {
	output := t.value.String()
	x.Equal(t.string, output)
}

func (t tlsVersionUnmarshalText) Name() string {
	return "TLSVersion_UnmarshalText/" + t.name
}

func (t tlsVersionUnmarshalText) Run(x *assert.Assertions) {
	var rv TLSVersion
	err := rv.UnmarshalText([]byte(t.string))
	if t.valid() {
		x.NoError(err)
		x.Equal(t.value, rv)
	} else {
		x.Error(err)
		x.Equal(TLSVersion(0), rv)
	}
}

func (t tlsVersionMarshalText) Name() string {
	return "TLSVersion_MarshalText/" + t.name
}

func (t tlsVersionMarshalText) Run(x *assert.Assertions) {
	rv, err := t.value.MarshalText()
	if t.valid() {
		x.NoError(err)
		x.Equal([]byte(t.string), rv)
	} else {
		x.Error(err)
		x.Nil(rv)
	}
}

func (t tlsVersionUnmarshalJSON) Name() string {
	return "TLSVersion_UnmarshalJSON/" + t.name
}

func (t tlsVersionUnmarshalJSON) Run(x *assert.Assertions) {
	var rv struct {
		V TLSVersion `json:"v"`
	}

	jsonString := fmt.Sprintf(`{"v":%q}`, t.string)
	err := json.Unmarshal([]byte(jsonString), &rv)

	if t.valid() {
		x.NoError(err)
		x.Equal(t.value, rv.V)
	} else {
		x.Error(err)
		x.Equal(TLSVersion(0), rv.V)
	}
}

func (t tlsVersionMarshalJSON) Name() string {
	return "TLSVersion_MarshalJSON/" + t.name
}

func (t tlsVersionMarshalJSON) Run(x *assert.Assertions) {
	var input = struct {
		V TLSVersion `json:"v"`
	}{t.value}

	output, err := json.Marshal(input)
	if t.valid() {
		jsonString := fmt.Sprintf(`{"v":%q}`, t.string)
		x.NoError(err)
		x.Equal([]byte(jsonString), output)
	} else {
		x.Error(err)
		x.Nil(output)
	}
}

func (t clientCertificateAuthorityTest) Name() string {
	return "ClientCertificateAuthority/" + t.name
}

func (t clientCertificateAuthorityTest) Run(x *assert.Assertions) {
	var rv ClientCertificateAuthority
	err := rv.UnmarshalText([]byte(t.pemData))
	if t.valid {
		x.NoError(err)
		text, err := rv.MarshalText()
		x.NoError(err)
		x.Equal(t.pemData, string(text))
	} else {
		x.Error(err)
	}
}

func (t jsonTest) Name() string {
	return "JSONTest/" + t.name
}

func (t jsonTest) Run(x *assert.Assertions) {
	rv, err := t.unmarshal(t.input)
	x.NoError(err)
	x.Equal(t.expected, rv)
}

func TestRoutes(t *testing.T) {
	testCases.Run(t)
}
