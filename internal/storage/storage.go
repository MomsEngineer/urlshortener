package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

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
}

func Create(dbDSN, fileName string) (*Storage, error) {
	realDB, err := realdb.NewRealDB(dbDSN)
	if err != nil {
		fmt.Println(err)
	}

	if fileName == "" {
		return &Storage{
			realdb: realDB,
			db:     db.NewDB(),
			file:   nil,
		}, nil
	}

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
		realdb: realDB,
		db:     db,
		file:   file,
	}, nil
}

func (s *Storage) Close() {
	if s.file != nil {
		s.file.Close()
	}

	if s.realdb != nil {
		s.realdb.Close()
	}
}

func (s *Storage) SaveLink(id, link string) {
	s.db.SaveLink(id, link)
	if s.file != nil {
		s.file.Write(id, link)
	}
}

func (s *Storage) GetLink(id string) (string, bool) {
	return s.db.GetLink(id)
}

func (s *Storage) Ping() error {
	if s.realdb == nil {
		return errors.New("The real DB has not been created")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.realdb.PingContext(ctx)
}
