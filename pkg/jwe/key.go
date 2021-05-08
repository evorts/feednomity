package jwe

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/pkg/errors"
)

type (
	Key struct {
		Value string
	}
	KeyStorage interface {
		GetPrivate() *rsa.PrivateKey
	}
)

func (k *Key) GetPrivate() (*rsa.PrivateKey, error) {
	return generateRsaPrivateKeyFromPemString(k.Value)
}

func generateRsaPrivateKeyFromPemString(privatePem string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privatePem))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}
	pri, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pri, nil
}
