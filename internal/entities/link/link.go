package link

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
)

type Link struct {
	ShortURL    string
	OriginalURL string
}

func NewLink(short, link string) (*Link, error) {
	var err error
	if short == "" {
		short, err = GenerateID(8)
		if err != nil {
			return nil, err
		}
	}

	return &Link{
		ShortURL:    short,
		OriginalURL: link,
	}, nil
}

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
