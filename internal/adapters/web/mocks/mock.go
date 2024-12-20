package mocks

import (
	"context"
	"errors"
)

type Storage struct{}

func (s *Storage) SaveLink(context.Context, string, string) (string, error) {
	return "", nil
}

func (s *Storage) SaveLinksBatch(context.Context, string, map[string]string) error {
	return nil
}

func (s *Storage) GetLink(_ context.Context, _ string, id string) (link string, err error) {
	if id == "abc123" {
		return "https://example.com", nil
	}
	return "", errors.New("not found")
}

func (s *Storage) GetLinksByUser(context.Context, string) (map[string]string, error) {
	return nil, nil
}

func (s *Storage) Ping(_ context.Context) error {
	return nil
}

func (s *Storage) Close() error {
	return nil
}
