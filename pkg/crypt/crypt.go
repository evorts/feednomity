package crypt

import (
	sha12 "crypto/sha1"
	"fmt"
	"hash"
	"io"
)

type crypt struct {
	salt string
	sha1 hash.Hash
}

type ICrypt interface {
	CryptWithSalt(value string) string
	Crypt(value string) string
	Renew() ICrypt
}

func NewCrypt(salt string) ICrypt {
	return &crypt{
		sha1: sha12.New(),
		salt: salt,
	}
}

func (c *crypt) CryptWithSalt(value string) string {
	_, err := io.WriteString(c.sha1, fmt.Sprintf("%s%s", c.salt, value))
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", c.sha1.Sum(nil))
}

func (c *crypt) Crypt(value string) string {
	_, err := io.WriteString(c.sha1, value)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", c.sha1.Sum(nil))
}

func (c *crypt) Renew() ICrypt {
	c.sha1 = sha12.New()
	return c
}
