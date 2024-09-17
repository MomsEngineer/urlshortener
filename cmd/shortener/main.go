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

func generateID(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b)[:n], nil
}

func CreatePostHandler(db *DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			id, err := generateID(8)
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
		} else {
			http.Error(w, "Only GET and POST requests are allowed!",
				http.StatusBadRequest)
		}
	}
}

func CreateGetHandler(db *DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			id := strings.TrimPrefix(r.URL.Path, "/")

			w.Header().Set("Location", db.Links[id])
			w.WriteHeader(http.StatusTemporaryRedirect)
		} else {
			http.Error(w, "Only GET and POST requests are allowed!",
				http.StatusBadRequest)
		}
	}
}

func main() {
	links := make(map[string]string)
	db := DB{
		Links: links,
	}

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, CreatePostHandler(&db))
	mux.HandleFunc(`/{id}`, CreateGetHandler(&db))

	if err := http.ListenAndServe("localhost:8080", mux); err != nil {
		panic(err)
	}
}
