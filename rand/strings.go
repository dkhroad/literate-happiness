package rand

import (
	"crypto/rand"
	"encoding/base64"
)

const NumBytes = 32

func randomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	n, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// url encoded random string
func RememberToken() (string, error) {
	tokenBytes, err := randomBytes(NumBytes)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(tokenBytes), nil

}
