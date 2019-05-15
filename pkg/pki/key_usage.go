package pki

import (
	"crypto/x509"
	"strings"
)

type KeyUsage x509.KeyUsage

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

func (u KeyUsage) String() string {
	value := x509.KeyUsage(u)
	usage := make([]string, 0, len(keyUsageMap))
	for k, v := range keyUsageMap {
		if value&k != 0 {
			usage = append(usage, v)
		}
	}

	return strings.Join(usage, ", ")
}
