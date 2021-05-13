package hapi

import (
	"encoding/json"
	"github.com/evorts/feednomity/pkg/crypt"
	"time"
)

type hashHelper struct {
	aes crypt.ICryptAES
}

type IHashHelper interface {
	Generate(expireAt time.Time, realHash string, attributes interface{}) string
	Decode(value string) (*HashData, error)
}

func NewHashHelper(aes crypt.ICryptAES) IHashHelper {
	return &hashHelper{aes: aes}
}

func (h *hashHelper) Generate(expireAt time.Time, realHash string, attributes interface{}) string {
	data := HashData{
		ExpireAt:   expireAt,
		RealHash:   realHash,
		Attributes: attributes,
	}
	jData, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	if hash, err2 := h.aes.Encrypt(string(jData)); err2 == nil {
		return hash
	}
	return ""
}

func (h *hashHelper) Decode(value string) (*HashData, error) {
	jData, err := h.aes.Decrypt(value)
	if err != nil {
		return nil, err
	}
	var data HashData
	err = json.Unmarshal([]byte(jData), &data)
	return &data, err
}

