package storage

import (
	"errors"

	"github.com/MomsEngineer/urlshortener/internal/logger"
	fs "github.com/MomsEngineer/urlshortener/internal/storage/file_storage"
	ms "github.com/MomsEngineer/urlshortener/internal/storage/map_storage"
)

var log = logger.Create()

type Storage interface {
	SaveLink(id, link string) error
	GetLink(id string) (string, bool, error)
	Ping() error
	Close() error
}

func Create(filePath string) (Storage, error) {
	if filePath != "" {
		f, err := fs.NewFileStorage(filePath)
		if err != nil {
			log.Error("Failed to create file storage", err)
			return nil, errors.New("failed to create file storage")
		}
		log.Debug(filePath)
		return f, nil
	}
	storage := ms.NewMapStorage()
	return storage, nil
}
