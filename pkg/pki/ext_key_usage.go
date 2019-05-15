package pki

import (
	"crypto/x509"
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
	value := []x509.ExtKeyUsage(u)
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
