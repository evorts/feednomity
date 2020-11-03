package crypt

type cryptic struct {
	salt string
	hash ICryptHash
	aes  ICryptAES
}

type ICryptic interface {
	ICryptHash
	ICryptAES
}

func NewCryptic(salt string, hash ICryptHash, aes ICryptAES) ICryptic {
	return &cryptic{
		salt: salt,
		hash: hash,
		aes:  aes,
	}
}

func (c *cryptic) HashWithSalt(value string) string {
	return c.HashWithSalt(value)
}

func (c *cryptic) HashWithoutSalt(value string) string {
	return c.HashWithoutSalt(value)
}

func (c *cryptic) RenewHash() ICryptHash {
	c.RenewHash()
	return c
}

func (c *cryptic) Initialize() (ICryptAES, error) {
	return c.Initialize()
}

func (c *cryptic) Encrypt(value string) (string, error) {
	return c.Encrypt(value)
}

func (c *cryptic) Decrypt(value string) (string, error) {
	return c.Decrypt(value)
}
