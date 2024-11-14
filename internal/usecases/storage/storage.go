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
	GetLinksByUser(ctx context.Context, userID string) (map[string]string, error)
	Ping(context.Context) error
	Close() error
}

type StoregeInterface interface {
	SaveLinksBatch(cxt context.Context, userID string, links map[string]string) error
	SaveLink(ctx context.Context, userID, original string) (string, error)
	GetLink(ctx context.Context, userID, short string) (string, error)
	GetLinksByUser(ctx context.Context, userID string) (map[string]string, error)
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

func (s *Storage) SaveLinksBatch(ctx context.Context, userID string, ls map[string]string) error {
	var links []*link.Link

	for id, original := range ls {
		l, err := link.NewLink(userID, "", original)
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

func (s *Storage) SaveLink(ctx context.Context, userID, original string) (string, error) {
	l, err := link.NewLink(userID, "", original)
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

func (s *Storage) GetLink(ctx context.Context, userID, short string) (string, error) {
	l, err := link.NewLink(userID, short, "")
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

func (s *Storage) GetLinksByUser(ctx context.Context, userID string) (map[string]string, error) {
	links, err := s.store.GetLinksByUser(ctx, userID)
	if err != nil {
		log.Error("Failed to get links", err)
		return nil, err
	}

	if len(links) == 0 {
		log.Debug("Not found link for userd id", userID)
		return nil, ierror.ErrNoContent
	}

	return links, nil
}
