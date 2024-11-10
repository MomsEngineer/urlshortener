package main

import (
	"github.com/MomsEngineer/urlshortener/internal/adapters/config"
	"github.com/MomsEngineer/urlshortener/internal/adapters/web"
	"github.com/MomsEngineer/urlshortener/internal/usecases/storage"
)

func main() {
	cfg := config.NewConfig()

	s, err := storage.Create(cfg.DataBaseDSN, cfg.FilePath)
	if err != nil {
		panic("could not create a storage")
	}
	defer s.Close()

	router := web.NewRouter()
	web.SetupRoutes(router, s, cfg.BaseURL)

	router.Run()
}
