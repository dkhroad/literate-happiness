package rand

import (
	"bytes"
	"testing"
)

func TestRandomBytes(t *testing.T) {
	t.Log("Should return random bytes")
	b, err := randomBytes(32)
	if err != nil {
		t.Errorf("\tDid not get random bytes", err)
	}
	c := make([]byte, 32)
	if bytes.Equal(b, c) {
		t.Errorf("\tDid not get random bytes", b)
	} else {
		t.Logf("\tGot random bytes", b)
	}
}

func TestRememberToken(t *testing.T) {
	token, err := RememberToken()
	if err != nil {
		t.Errorf("RemeberToken() failed with error", err)
	}

	t.Logf("RememberToken :", token)
}
