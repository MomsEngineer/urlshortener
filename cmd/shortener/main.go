package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type DB struct {
	Links map[string]string
}

func generateId(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b)[:n], nil
}

func (db *DB) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		id, err := generateId(8)
		if err != nil {
			http.Error(w, "Server is not available!",
				http.StatusInternalServerError)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read request body",
				http.StatusInternalServerError)
			return
		}

		db.Links[id] = string(body)

		url := "http://localhost:8080/" + id

		w.Header().Set("content-type", "text/plain")
		w.Header().Set("Content-Length", fmt.Sprint(len(url)))
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(url))

	case http.MethodGet:
		id := strings.TrimPrefix(r.URL.Path, "/")

		w.Header().Set("Location", db.Links[id])
		w.WriteHeader(http.StatusTemporaryRedirect)

	default:
		http.Error(w, "Only GET and POST requests are allowed!",
			http.StatusBadRequest)
	}
}

func main() {
	links := make(map[string]string)
	handler := DB{
		Links: links,
	}

	if err := http.ListenAndServe("localhost:8080", &handler); err != nil {
		panic(err)
	}
}
