package url

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
)

type URL struct {
	ShortURL    string
	OriginalURL string
}

func NewURL(baseURL, link string) (*URL, error) {
	short, err := GenerateID(8)
	if err != nil {
		return nil, err
	}

	shortURL := baseURL + "/" + short

	return &URL{
		ShortURL:    shortURL,
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
