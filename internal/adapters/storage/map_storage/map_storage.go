package mapstorage

import (
	"context"
	"errors"

	"github.com/MomsEngineer/urlshortener/internal/entities/link"
)

type MapStorage struct {
	Links map[string]string
}

func NewMapStorage() *MapStorage {
	return &MapStorage{
		Links: make(map[string]string),
	}
}

func (lm *MapStorage) SaveLink(_ context.Context, l *link.Link) error {
	lm.Links[l.ShortURL] = l.OriginalURL
	return nil
}

func (lm *MapStorage) SaveLinksBatch(_ context.Context, links []*link.Link) error {
	for _, l := range links {
		lm.Links[l.ShortURL] = l.OriginalURL
	}
	return nil
}

func (lm *MapStorage) GetLink(_ context.Context, link *link.Link) error {
	if o, exists := lm.Links[link.ShortURL]; exists {
		link.OriginalURL = o
		return nil
	}

	return errors.New("not found")
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
