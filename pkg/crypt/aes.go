package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

type cryptAES struct {
	salt   string
	cipher cipher.Block
	gcm    cipher.AEAD
}

type ICryptAES interface {
	Initialize() (ICryptAES, error)
	Encrypt(value string) (string, error)
	Decrypt(value string) (string, error)
}

func NewCryptAES(salt string) ICryptAES {
	return &cryptAES{
		salt: hex.EncodeToString([]byte(salt)),
	}
}

func (c *cryptAES) Initialize() (ic ICryptAES, err error) {
	key, _ := hex.DecodeString(c.salt)
	c.cipher, err = aes.NewCipher(key)
	if err != nil {
		return
	}
	c.gcm, err = cipher.NewGCM(c.cipher)
	return c, err
}

func (c *cryptAES) Encrypt(value string) (rs string, err error) {
	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, c.gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := c.gcm.Seal(nonce, nonce, []byte(value), nil)
	return fmt.Sprintf("%x", ciphertext), nil
}

func (c *cryptAES) Decrypt(value string) (rs string, err error) {
	enc, _ := hex.DecodeString(value)
	//Get the nonce size
	nonceSize := c.gcm.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	var dec []byte
	dec, err = c.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return
	}
	rs = fmt.Sprintf("%s", dec)
	return
}
