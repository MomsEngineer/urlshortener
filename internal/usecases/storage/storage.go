package storage

import (
	"context"
	"errors"

	"github.com/MomsEngineer/urlshortener/internal/adapters/logger"
	db "github.com/MomsEngineer/urlshortener/internal/adapters/storage/db_storage"
	fs "github.com/MomsEngineer/urlshortener/internal/adapters/storage/file_storage"
	ms "github.com/MomsEngineer/urlshortener/internal/adapters/storage/map_storage"
	"github.com/MomsEngineer/urlshortener/internal/entities/link"
	ierror "github.com/MomsEngineer/urlshortener/internal/errors"
)

var log = logger.Create()

type StoreInterface interface {
	SaveLinksBatch(context.Context, []*link.Link) error
	SaveLink(context.Context, *link.Link) error
	GetLink(context.Context, *link.Link) error
	Ping(context.Context) error
	Close() error
}

type StoregeInterface interface {
	SaveLinksBatch(context.Context, map[string]string) error
	SaveLink(context.Context, string) (string, error)
	GetLink(context.Context, string) (string, error)
	Ping(context.Context) error
	Close() error
}

type Storage struct {
	store StoreInterface
}

func Create(dsn, filePath string) (StoregeInterface, error) {
	if dsn != "" {
		store, err := db.NewDB(dsn)
		if err != nil {
			log.Error("Failed to create DB storage", err)
			return nil, err
		}
		log.Debug("Created DB")

		return &Storage{store: store}, nil
	} else if filePath != "" {
		store, err := fs.NewFileStorage(filePath)
		if err != nil {
			log.Error("Failed to create file storage", err)
			return nil, err
		}
		log.Debug("Created file storage")

		return &Storage{store: store}, nil
	}

	store := ms.NewMapStorage()
	log.Debug("Created map storage")

	return &Storage{store: store}, nil
}

func (s *Storage) SaveLinksBatch(ctx context.Context, ls map[string]string) error {
	var links []*link.Link

	for id, original := range ls {
		l, err := link.NewLink("", original)
		if err != nil {
			log.Error("Failed to create new link", err)
			return err
		}

		links = append(links, l)
		ls[id] = l.ShortURL
	}

	if err := s.store.SaveLinksBatch(ctx, links); err != nil {
		log.Error("Failed to save links batch", err)
		return err
	}

	return nil
}

func (s *Storage) SaveLink(ctx context.Context, original string) (string, error) {
	l, err := link.NewLink("", original)
	if err != nil {
		log.Error("Failed to create new link", err)
		return "", err
	}

	if err := s.store.SaveLink(ctx, l); err != nil {
		if errors.Is(err, ierror.ErrDuplicate) {
			return l.ShortURL, err
		}
		log.Error("Failed to get link", err)
		return "", err
	}

	return l.ShortURL, nil
}

func (s *Storage) GetLink(ctx context.Context, short string) (string, error) {
	l, err := link.NewLink(short, "")
	if err != nil {
		log.Error("Failed to create new link", err)
		return "", err
	}

	if err := s.store.GetLink(ctx, l); err != nil {
		log.Error("Failed to get link", err)
		return "", err
	}

	return l.OriginalURL, nil
}

func (s *Storage) Close() error {
	return s.store.Close()
}

func (s *Storage) Ping(ctx context.Context) error {
	return s.store.Ping(ctx)
}
