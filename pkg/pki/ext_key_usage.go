package pki

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"strings"
)

type ExtKeyUsage []x509.ExtKeyUsage

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

func (u ExtKeyUsage) String() string {
	return strings.Join(u.Marshal(), ", ")
}

func (u ExtKeyUsage) Marshal() []string {
	value := []x509.ExtKeyUsage(u)
	usage := make([]string, len(value))
	for idx, v := range value {
		if name, ok := extKeyUsageMap[v]; ok {
			usage[idx] = name
		} else {
			usage[idx] = fmt.Sprintf("%02x", int(v))
		}
	}
	return usage
}

func (u *ExtKeyUsage) Unmarshal(usage []string) error {
	*u = make(ExtKeyUsage, len(usage))
	usageMap := make(map[string]x509.ExtKeyUsage, len(extKeyUsageMap))
	for k, v := range extKeyUsageMap {
		usageMap[v] = k
	}
	for idx, use := range usage {
		if v, ok := usageMap[use]; ok {
			(*u)[idx] = v
		}
	}
	return nil
}

func (u ExtKeyUsage) MarshalText() ([]byte, error) {
	return []byte(u.String()), nil
}

func (u *ExtKeyUsage) UnmarshalText(text []byte) error {
	usage := strings.Split(string(text), ", ")
	return u.Unmarshal(usage)
}

func (u ExtKeyUsage) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.Marshal())
}

func (u *ExtKeyUsage) UnmarshalJSON(data []byte) error {
	var usage []string
	err := json.Unmarshal(data, &usage)
	if err != nil {
		return err
	}
	return u.Unmarshal(usage)
}

func (u ExtKeyUsage) MarshalYAML() (interface{}, error) {
	return u.Marshal(), nil
}

func (u *ExtKeyUsage) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var usage []string
	err := unmarshal(&usage)
	if err != nil {
		return err
	}
	return u.Unmarshal(usage)
}
