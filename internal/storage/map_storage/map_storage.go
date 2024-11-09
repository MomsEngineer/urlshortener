package mapstorage

import (
	"context"
	"errors"
)

type MapStorage struct {
	Links map[string]string
}

func NewMapStorage() *MapStorage {
	return &MapStorage{
		Links: make(map[string]string),
	}
}

func (lm *MapStorage) SaveLink(_ context.Context, id, link string) (string, error) {
	lm.Links[id] = link
	return "", nil
}

func (lm *MapStorage) SaveLinksBatch(_ context.Context, links map[string]string) error {
	for k, v := range links {
		lm.Links[k] = v
	}
	return nil
}

func (lm *MapStorage) GetLink(_ context.Context, id string) (string, bool, error) {
	link, exists := lm.Links[id]
	return link, exists, nil
}

func (lm *MapStorage) Ping(_ context.Context) error {
	if lm.Links == nil {
		return errors.New("links is nil")
	}
	return nil
}

func (lm *MapStorage) Close() error {
	return nil
}
