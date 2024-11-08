package main

import (
	"github.com/MomsEngineer/urlshortener/internal/compresser"
	"github.com/MomsEngineer/urlshortener/internal/config"
	"github.com/MomsEngineer/urlshortener/internal/handlers"
	"github.com/MomsEngineer/urlshortener/internal/logger"
	"github.com/MomsEngineer/urlshortener/internal/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	cfg := config.NewConfig()

	log := logger.Create()

	s, err := storage.Create(cfg.DataBaseDSN, cfg.FilePath)
	if err != nil {
		panic("could not create a storage")
	}
	defer s.Close()

	router := gin.New()
	router.SetTrustedProxies(nil)

	router.Use(log.Logger())
	router.Use(compresser.CompresserMiddleware())

	router.POST("/", func(c *gin.Context) {
		handlers.HandlePost(c, s, cfg.BaseURL)
	})

	router.POST("/api/shorten", func(c *gin.Context) {
		handlers.HandlePostAPI(c, s, cfg.BaseURL)
	})

	router.POST("/api/shorten/batch", func(c *gin.Context) {
		handlers.HandlePostBatch(c, s, cfg.BaseURL)
	})

	router.GET("/:id", func(c *gin.Context) {
		handlers.HandleGet(c, s)
	})

	router.GET("/ping", func(c *gin.Context) {
		handlers.HandlePing(c, s)
	})

	router.Run(cfg.Address)
}
