package storage

import "github.com/MomsEngineer/urlshortener/internal/storage/db"

type LinkStorage interface {
	SaveLink(id, link string)
	GetLink(id string) (string, bool)
}

type Storage struct {
	db *db.DB
}

func Create() (*Storage, error) {
	return &Storage{db: db.NewDB()}, nil
}

func (s *Storage) SaveLink(id, link string) {
	s.db.SaveLink(id, link)
}

func (s *Storage) GetLink(id string) (string, bool) {
	return s.db.GetLink(id)
}
