package app

import (
	"io"
	"net/http"
	"strings"

	"github.com/MomsEngineer/urlshortener/internal/db"
	"github.com/MomsEngineer/urlshortener/internal/utils"
)

func Handler(database db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlePost(w, r, database)
		case http.MethodGet:
			handleGet(w, r, database)
		default:
			http.Error(w, "Only GET and POST requests are allowed!",
				http.StatusBadRequest)
		}
	}
}

func handlePost(w http.ResponseWriter, r *http.Request, database db.Database) {
	link, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body",
			http.StatusInternalServerError)
		return
	}

	id, err := utils.GenerateID(8)
	if err != nil {
		http.Error(w, "Failed to generate short link!",
			http.StatusInternalServerError)
		return
	}

	database.SaveLink(id, string(link))
	shortUrl := "http://localhost:8080/" + id

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortUrl))
}

func handleGet(w http.ResponseWriter, r *http.Request, database db.Database) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	link, exists := database.GetLink(id)
	if !exists {
		http.Error(w, "Link not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Location", link)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
