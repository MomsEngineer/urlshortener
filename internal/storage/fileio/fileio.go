package fileio

import (
	"encoding/json"
	"io"
	"os"
	"strconv"
)

type entry struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Reader struct {
	file    *os.File
	decoder *json.Decoder
}

type Writer struct {
	file    *os.File
	encoder *json.Encoder
	counter uint
}

type FileIO struct {
	r *Reader
	w *Writer
}

func NewFileIO(fileName string) (*FileIO, error) {
	r, err := NewReader(fileName)
	if err != nil {
		return nil, err
	}

	w, err := NewWriter(fileName)
	if err != nil {
		r.close()
		return nil, err
	}

	return &FileIO{
		r: r,
		w: w,
	}, nil
}

func (f *FileIO) Read() (map[string]string, error) {
	defer f.r.close()
	m := map[string]string{}

	for {
		entry, err := f.r.readEntry()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		m[entry.ShortURL] = entry.OriginalURL
	}

	return m, nil
}

func (f *FileIO) Write(shortURL, originalURL string) error {
	return f.w.writeEntry(&entry{
		UUID:        strconv.FormatUint(uint64(f.w.counter+1), 10),
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	})
}

func (f *FileIO) Close() error {
	return f.w.close()
}

func NewReader(fileName string) (*Reader, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Reader{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (r *Reader) readEntry() (*entry, error) {
	e := &entry{}
	err := r.decoder.Decode(e)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (r *Reader) close() error {
	return r.file.Close()
}

func NewWriter(fileName string) (*Writer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &Writer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (w *Writer) writeEntry(e *entry) error {
	return w.encoder.Encode(e)
}

func (w *Writer) close() error {
	return w.file.Close()
}
