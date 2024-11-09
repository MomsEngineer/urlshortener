package mocks

import "context"

type Storage struct{}

func (s *Storage) SaveLink(_ context.Context, _, _ string) (string, error) {
	return "", nil
}

func (s *Storage) SaveLinksBatch(_ context.Context, _ map[string]string) error {
	return nil
}

func (s *Storage) GetLink(_ context.Context, id string) (link string, exists bool, err error) {
	if id == "abc123" {
		link, exists = "https://example.com", true
	} else {
		link, exists = "", false
	}
	return
}

func (s *Storage) Ping(_ context.Context) error {
	return nil
}

func (s *Storage) Close() error {
	return nil
}
