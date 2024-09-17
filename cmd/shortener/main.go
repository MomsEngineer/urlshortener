package main

import (
	"net/http"

	"github.com/MomsEngineer/urlshortener/internal/app"
	"github.com/MomsEngineer/urlshortener/internal/db"
)

func main() {
	db := db.NewDB()

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, app.Handler(db))

	if err := http.ListenAndServe("localhost:8080", mux); err != nil {
		panic(err)
	}
}
