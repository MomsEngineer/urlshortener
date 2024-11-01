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

	s, _ := storage.Create()

	router := gin.New()
	router.SetTrustedProxies(nil)

	router.Use(logger.Create().Logger())
	router.Use(compresser.CompresserMiddleware())

	router.POST("/", func(c *gin.Context) {
		handlers.HandlePost(c, s, cfg.BaseURL)
	})

	router.POST("/api/shorten", func(c *gin.Context) {
		handlers.HandlePostAPI(c, s, cfg.BaseURL)
	})

	router.GET("/:id", func(c *gin.Context) {
		handlers.HandleGet(c, s)
	})

	router.Run(cfg.Address)
}
