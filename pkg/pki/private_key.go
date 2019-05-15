package pki

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/sirupsen/logrus"
)

type PrivateKey struct {
	*rsa.PrivateKey
	secret []byte
}

func NewPrivateKey(bits int) *PrivateKey {
	pk, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		logrus.WithField("err", err).WithField("bits", bits).Panic("error generating private key")
	}
	return &PrivateKey{PrivateKey: pk}
}

func NewPrivateKeyWithSecret(bits int, secret []byte) *PrivateKey {
	rv := NewPrivateKey(bits)
	rv.secret = secret
	return rv
}

func (pk *PrivateKey) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		pk.PrivateKey = nil
		return nil
	}

	pemBlock, err := UnmarshalBlock(string(text))
	if err != nil {
		return err
	}

	var data []byte
	if len(pk.secret) > 0 {
		data, err = x509.DecryptPEMBlock(pemBlock, pk.secret)
		if err != nil {
			return err
		}
	} else {
		data = pemBlock.Bytes
	}

	pk.PrivateKey, err = x509.ParsePKCS1PrivateKey(data)
	return err
}

func (pk PrivateKey) MarshalText() ([]byte, error) {
	if pk.PrivateKey == nil {
		return nil, nil
	}

	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(pk.PrivateKey),
	}

	var err error

	if len(pk.secret) > 0 {
		block, err = x509.EncryptPEMBlock(rand.Reader, block.Type, block.Bytes, pk.secret, x509.PEMCipherAES256)
		if err != nil {
			return nil, err
		}
	}

	return pem.EncodeToMemory(block), nil
}

func (pk *PrivateKey) SetSecret(secret []byte) {
	pk.secret = secret
}
