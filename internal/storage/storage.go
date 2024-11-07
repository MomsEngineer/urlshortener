package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/MomsEngineer/urlshortener/internal/logger"
	"github.com/MomsEngineer/urlshortener/internal/storage/db"
	"github.com/MomsEngineer/urlshortener/internal/storage/fileio"
	"github.com/MomsEngineer/urlshortener/internal/storage/memory"
)

type LinkStorage interface {
	Ping() error
	SaveLink(id, link string) error
	GetLink(id string) (string, bool, error)
}

type Storage struct {
	db     *sql.DB
	memory *memory.LinksMap
	file   *fileio.FileIO
	log    logger.Logger
}

func getFileIO(lm *memory.LinksMap, fileName string) (*fileio.FileIO, error) {
	file, err := fileio.NewFileIO(fileName)
	if err != nil {
		return nil, err
	}

	m, err := file.Read()
	if err != nil {
		return nil, err
	}

	for k, v := range m {
		lm.SaveLink(k, v)
	}

	return file, nil
}

func Create(log logger.Logger, dbDSN, fileName string) (*Storage, error) {
	storage := &Storage{
		memory: memory.NewLinksMap(),
		log:    log,
	}

	db, err := db.NewDB(dbDSN)
	if err != nil {
		log.Error("Failed to create DB", err)
	}
	storage.db = db

	if fileName == "" {
		return storage, nil
	}

	file, err := getFileIO(storage.memory, fileName)
	if err != nil {
		log.Error("Failed to create file for IO", err)
	}
	storage.file = file

	return storage, nil
}

func (s *Storage) Close() {
	if s.file != nil {
		s.file.Close()
		s.log.Debug("Closed the file:", s.file.Name)
	}

	if s.db != nil {
		s.log.Debug("Closed the realdb")
		s.db.Close()
	}
}

func (s *Storage) SaveLink(id, link string) error {
	if s.file != nil {
		if err := s.file.Write(id, link); err != nil {
			s.log.Error("Failed to write to file", err)
			return err
		}
		s.log.Debug("Saved to file:", s.file.Name)
	}

	s.memory.SaveLink(id, link)
	s.log.Debug("Saved to memory")

	return nil
}

func (s *Storage) GetLink(id string) (string, bool, error) {
	link, exist := s.memory.GetLink(id)
	s.log.Debug("Get link from memory. Link:", link)
	return link, exist, nil
}

func (s *Storage) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if s.db == nil {
		return errors.New("the db is nil")
	}

	return s.db.PingContext(ctx)
}
