package mapstorage

import (
	"context"
	"errors"

	"github.com/MomsEngineer/urlshortener/internal/entities/link"
)

type MapStorage struct {
	Links []*link.Link
}

func NewMapStorage() *MapStorage {
	return &MapStorage{}
}

func (lm *MapStorage) SaveLink(_ context.Context, l *link.Link) error {
	lm.Links = append(lm.Links, l)
	return nil
}

func (lm *MapStorage) SaveLinksBatch(_ context.Context, links []*link.Link) error {
	for _, l := range links {
		lm.SaveLink(context.TODO(), l)
	}
	return nil
}

func (lm *MapStorage) GetLink(_ context.Context, link *link.Link) error {
	for _, l := range lm.Links {
		if l.ShortURL == link.ShortURL {
			link.OriginalURL = l.OriginalURL
			return nil
		}
	}

	return errors.New("not found")
}

func (lm *MapStorage) GetLinksByUser(ctx context.Context, userID string) (map[string]string, error) {
	res := make(map[string]string)
	for _, l := range lm.Links {
		if l.UserID == userID {
			res[l.ShortURL] = l.OriginalURL
		}
	}

	return res, nil
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
