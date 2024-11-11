package filestorage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/MomsEngineer/urlshortener/internal/adapters/logger"
	"github.com/MomsEngineer/urlshortener/internal/entities/link"
)

var log = logger.Create()

type entry struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type reader struct {
	file    *os.File
	decoder *json.Decoder
}

type writer struct {
	file    *os.File
	encoder *json.Encoder
}

type FileStorage struct {
	r       *reader
	w       *writer
	path    string
	counter uint64
}

func NewFileStorage(path string) (*FileStorage, error) {
	fs := &FileStorage{path: path}

	r, err := newReader(path)
	if err != nil {
		log.Error("Failed to create reader", err)
		return nil, err
	}
	fs.r = r

	w, err := newWriter(path)
	if err != nil {
		fs.Close()
		log.Error("Failed to create writer", err)
		return nil, err
	}
	fs.w = w

	var lastEntry entry
	for {
		entry, err := fs.r.readEntry()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Error("Failed to read entry", err)
			return nil, err
		}
		lastEntry = *entry
	}

	if lastEntry.UUID != "" {
		counter, err := strconv.ParseUint(lastEntry.UUID, 10, 64)
		if err != nil {
			log.Error("Failed to parse counter from UUID", err)
			return nil, err
		}
		fs.counter = counter
	}

	return fs, nil
}

func (fs *FileStorage) SaveLinksBatch(_ context.Context, ls []*link.Link) error {
	for _, l := range ls {
		e := &entry{
			UUID:        strconv.FormatUint(uint64(fs.counter+1), 10),
			ShortURL:    l.ShortURL,
			OriginalURL: l.OriginalURL,
		}

		if err := fs.w.writeEntry(e); err != nil {
			log.Error("Failed to save link", err)
			return err
		}

		fs.counter++
	}

	return nil
}

func (fs *FileStorage) SaveLink(_ context.Context, l *link.Link) error {
	e := &entry{
		UUID:        strconv.FormatUint(uint64(fs.counter+1), 10),
		ShortURL:    l.ShortURL,
		OriginalURL: l.OriginalURL,
	}

	if err := fs.w.writeEntry(e); err != nil {
		log.Error("Failed to save link", err)
		return err
	}

	fs.counter++

	return nil
}

func (fs *FileStorage) GetLink(_ context.Context, l *link.Link) error {
	_, err := fs.r.file.Seek(0, 0)
	if err != nil {
		log.Error("Failed to seek to the beginning of the file", err)
		return err
	}

	fs.r.decoder = json.NewDecoder(fs.r.file)

	for {
		entry, err := fs.r.readEntry()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Error("Failed to read entry", err)
			return err
		}

		if entry.ShortURL == l.ShortURL {
			l.OriginalURL = entry.OriginalURL
			return nil
		}
	}

	return errors.New("not found")
}

func (fs *FileStorage) Ping(_ context.Context) error {
	if fs.r == nil || fs.w == nil {
		log.Error("Reader or writer is not initialized", nil)
		return fmt.Errorf("reader or writer is not initialized")
	}

	return nil
}

func (fs *FileStorage) Close() error {
	errR := fs.r.file.Close()
	errW := fs.w.file.Close()

	if errR != nil || errW != nil {
		return fmt.Errorf("failed to close reader: %v, failed to close writer: %v",
			errR, errW)
	}
	return nil
}

func newReader(fileName string) (*reader, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &reader{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (r *reader) readEntry() (*entry, error) {
	e := &entry{}
	err := r.decoder.Decode(e)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func newWriter(fileName string) (*writer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &writer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (w *writer) writeEntry(e *entry) error {
	return w.encoder.Encode(e)
}
