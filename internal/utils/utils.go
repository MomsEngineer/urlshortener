package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
)

func GenerateID(n int) (string, error) {
	if n <= 0 {
		return "", errors.New("n must be greater 0")
	}

	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b)[:n], nil
}
