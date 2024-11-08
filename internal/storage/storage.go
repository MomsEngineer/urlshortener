package storage

import (
	"github.com/MomsEngineer/urlshortener/internal/logger"
	db "github.com/MomsEngineer/urlshortener/internal/storage/db_storage"
	fs "github.com/MomsEngineer/urlshortener/internal/storage/file_storage"
	ms "github.com/MomsEngineer/urlshortener/internal/storage/map_storage"
)

var log = logger.Create()

type Storage interface {
	SaveLink(shortLink, originalLink string) error
	GetLink(shortLink string) (string, bool, error)
	Ping() error
	Close() error
}

func Create(dsn, filePath string) (Storage, error) {
	if dsn != "" {
		storage, err := db.NewDB(dsn)
		if err != nil {
			log.Error("Failed to create DB storage", err)
			return nil, err
		}
		log.Debug("Created DB")

		return storage, nil
	} else if filePath != "" {
		storage, err := fs.NewFileStorage(filePath)
		if err != nil {
			log.Error("Failed to create file storage", err)
			return nil, err
		}
		log.Debug("Created file storage")

		return storage, nil
	}

	storage := ms.NewMapStorage()
	log.Debug("Created map storage")

	return storage, nil
}
