package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/MomsEngineer/urlshortener/internal/logger"
	"github.com/MomsEngineer/urlshortener/internal/storage/db"
	"github.com/MomsEngineer/urlshortener/internal/storage/fileio"
	"github.com/MomsEngineer/urlshortener/internal/storage/realdb"
)

type LinkStorage interface {
	Ping() error
	SaveLink(id, link string)
	GetLink(id string) (string, bool)
}

type Storage struct {
	realdb *sql.DB
	db     *db.DB
	file   *fileio.FileIO
	log    logger.Logger
}

func getFileIO(db *db.DB, fileName string) (*fileio.FileIO, error) {
	file, err := fileio.NewFileIO(fileName)
	if err != nil {
		return nil, err
	}

	m, err := file.Read()
	if err != nil {
		return nil, err
	}

	for k, v := range m {
		db.SaveLink(k, v)
	}

	return file, nil
}

func Create(log logger.Logger, dbDSN, fileName string) (*Storage, error) {
	storage := &Storage{
		db:  db.NewDB(),
		log: log,
	}

	realDB, err := realdb.NewRealDB(dbDSN)
	if err != nil {
		log.Error("Failed to create DB", err)
	}
	storage.realdb = realDB

	if fileName == "" {
		return storage, nil
	}

	file, err := getFileIO(storage.db, fileName)
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

	if s.realdb != nil {
		s.log.Debug("Closed the realdb")
		s.realdb.Close()
	}
}

func (s *Storage) SaveLink(id, link string) {
	if s.file != nil {
		s.file.Write(id, link)
		s.log.Debug("Saved to file:", s.file.Name)
	}

	s.db.SaveLink(id, link)
	s.log.Debug("Saved to db")
}

func (s *Storage) GetLink(id string) (string, bool) {
	return s.db.GetLink(id)
}

func (s *Storage) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return s.realdb.PingContext(ctx)
}
