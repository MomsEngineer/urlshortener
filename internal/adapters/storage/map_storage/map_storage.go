package memory

import (
	"context"
	"errors"
)

type MemoryStorage struct {
	Links map[string]string
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		Links: make(map[string]string),
	}
}

func (lm *MemoryStorage) SaveLink(_ context.Context, id, link string) (string, error) {
	lm.Links[id] = link
	return "", nil
}

func (lm *MemoryStorage) SaveLinksBatch(_ context.Context, links map[string]string) error {
	for k, v := range links {
		lm.Links[k] = v
	}
	return nil
}

func (lm *MemoryStorage) GetLink(_ context.Context, id string) (string, bool, error) {
	link, exists := lm.Links[id]
	return link, exists, nil
}

func (lm *MemoryStorage) Ping(_ context.Context) error {
	if lm.Links == nil {
		return errors.New("links is nil")
	}
	return nil
}

func (lm *MemoryStorage) Close() error {
	return nil
}
