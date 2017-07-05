package hash

import (
	"testing"
)

func TestHash(t *testing.T) {
	h := NewHMAC("my-secret-key")
	hash := h.Hash("hash-this")
	t.Logf("hash :", hash)
}
