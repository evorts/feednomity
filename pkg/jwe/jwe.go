package jwe

import (
	"crypto/rsa"
	"time"

	"gopkg.in/square/go-jose.v2"

	"github.com/pkg/errors"
	"gopkg.in/square/go-jose.v2/jwt"
)

const (
	ISSUER                = "onelabs.co"
	DefaultExpiration int = 1 //in hour
)

type (
	IManager interface {
		Encode(pub jwt.Claims, pri PrivateClaims) (token string, err error)
		Decode(token string) (pub jwt.Claims, pri PrivateClaims, err error)
		GetExpire() time.Duration
	}
	JWE struct {
		i bool
		k *rsa.PrivateKey
		b jwt.NestedBuilder
		e int64 //expiration
	}
	PrivateClaims struct {
		ClientId    string                 `json:"client_id"`
		Id          int64                  `json:"id"`
		Username    string                 `json:"username"`
		DisplayName string                 `json:"display_name"`
		Attributes  map[string]interface{} `json:"attributes"`
		Email       string                 `json:"email"`
		Phone       string                 `json:"phone"`
		AccessRole  string                 `json:"access_role"`
		JobRole     string                 `json:"job_role"`
		Assignment  string                 `json:"assignment"`
		GroupId     int64                  `json:"group_id"`
		OrgId       int64                  `json:"org_id"`
		OrgGroupIds []int64                `json:"org_group_ids"`
	}
)

var (
	ErrInvalidScope = errors.New("invalid scope in claims")
	ErrNoTokenFound = errors.New("auth: no credentials attached in request")
)

func NewJWE(key *rsa.PrivateKey, exp int64) IManager {
	return &JWE{k: key, e: exp}
}

func (j *JWE) builder() (builder jwt.NestedBuilder, err error) {
	if j.i {
		return j.b, nil
	}

	signingKey := jose.SigningKey{Algorithm: jose.RS256, Key: j.k}

	// create a Square.jose RSA signer, used to sign the JWT
	signerOpts := (&jose.SignerOptions{}).WithContentType("JWT")
	rsaSigner, err := jose.NewSigner(signingKey, signerOpts)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	encOpts := (&jose.EncrypterOptions{}).WithContentType("JWT")
	enc, err := jose.NewEncrypter(jose.A128GCM, jose.Recipient{Algorithm: jose.RSA_OAEP, Key: &j.k.PublicKey}, encOpts)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	builder = jwt.SignedAndEncrypted(rsaSigner, enc)

	j.b = builder
	j.i = true

	return
}

func (j *JWE) Encode(pub jwt.Claims, pri PrivateClaims) (token string, err error) {
	var b jwt.NestedBuilder
	b, err = j.builder()
	if err != nil {
		return
	}

	if token, err = b.Claims(pub).Claims(pri).CompactSerialize(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

func (j *JWE) Decode(token string) (pub jwt.Claims, pri PrivateClaims, err error) {
	var parsed *jwt.NestedJSONWebToken
	parsed, err = jwt.ParseSignedAndEncrypted(token)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	decrypted, err := parsed.Decrypt(j.k)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	if err = decrypted.Claims(&j.k.PublicKey, &pub, &pri); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

func (j *JWE) GetExpire() time.Duration {
	if j.e < 1 {
		return time.Duration(DefaultExpiration) * time.Hour
	}
	return time.Duration(j.e) * time.Hour
}
