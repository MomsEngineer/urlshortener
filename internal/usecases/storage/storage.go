package storage

import (
	"context"

	"github.com/MomsEngineer/urlshortener/internal/adapters/logger"
	db "github.com/MomsEngineer/urlshortener/internal/adapters/storage/db_storage"
	fs "github.com/MomsEngineer/urlshortener/internal/adapters/storage/file_storage"
	ms "github.com/MomsEngineer/urlshortener/internal/adapters/storage/map_storage"
)

var log = logger.Create()

type Storage interface {
	SaveLinksBatch(context.Context, map[string]string) error
	SaveLink(context.Context, string, string) (string, error)
	GetLink(context.Context, string) (string, bool, error)
	Ping(context.Context) error
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

	storage := ms.NewMemoryStorage()
	log.Debug("Created map storage")

	return storage, nil
}
