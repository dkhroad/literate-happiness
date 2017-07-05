package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hash"
)

func NewHMAC(key string) HMAC {
	h := hmac.New(sha256.New, []byte(key))
	return HMAC{h}
}

func (h HMAC) Hash(input string) string {
	h.hmac.Reset()
	h.hmac.Write([]byte(input))
	mac := h.hmac.Sum(nil)
	return base64.URLEncoding.EncodeToString(mac)
}

type HMAC struct {
	hmac hash.Hash
}
