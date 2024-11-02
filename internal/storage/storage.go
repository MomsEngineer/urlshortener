package storage

import (
	"github.com/MomsEngineer/urlshortener/internal/storage/db"
	"github.com/MomsEngineer/urlshortener/internal/storage/fileio"
)

type LinkStorage interface {
	SaveLink(id, link string)
	GetLink(id string) (string, bool)
}

type Storage struct {
	db   *db.DB
	file *fileio.FileIO
}

func Create(fileName string) (*Storage, error) {
	file, err := fileio.NewFileIO(fileName)
	if err != nil {
		return nil, err
	}

	db := db.NewDB()

	m, err := file.Read()
	if err != nil {
		return nil, err
	}

	for k, v := range m {
		db.SaveLink(k, v)
	}

	return &Storage{
		db:   db,
		file: file,
	}, nil
}

func (s *Storage) Close() error {
	return s.file.Close()
}

func (s *Storage) SaveLink(id, link string) {
	s.db.SaveLink(id, link)
	s.file.Write(id, link)
}

func (s *Storage) GetLink(id string) (string, bool) {
	return s.db.GetLink(id)
}
