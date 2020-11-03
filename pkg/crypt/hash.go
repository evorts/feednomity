package crypt

import (
	sha12 "crypto/sha1"
	"fmt"
	"hash"
	"io"
)

type cryptHash struct {
	salt string
	sha1 hash.Hash
}

type ICryptHash interface {
	HashWithSalt(value string) string
	HashWithoutSalt(value string) string
	RenewHash() ICryptHash
}

func NewHashEncryption(salt string) ICryptHash {
	return &cryptHash{
		sha1: sha12.New(),
		salt: salt,
	}
}

func (c *cryptHash) HashWithSalt(value string) string {
	_, err := io.WriteString(c.sha1, fmt.Sprintf("%s%s", c.salt, value))
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", c.sha1.Sum(nil))
}

func (c *cryptHash) HashWithoutSalt(value string) string {
	_, err := io.WriteString(c.sha1, value)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", c.sha1.Sum(nil))
}

func (c *cryptHash) RenewHash() ICryptHash {
	c.sha1 = sha12.New()
	return c
}
