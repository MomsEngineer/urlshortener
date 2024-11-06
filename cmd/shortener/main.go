package main

import (
	"fmt"

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

	s, err := storage.Create(cfg.DataBaseDSN, cfg.FilePath)
	if err != nil {
		fmt.Println("Could not create a storage")
	}
	defer s.Close()

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

	router.GET("/ping", func(c *gin.Context) {
		handlers.HandlePing(c, s)
	})

	router.Run(cfg.Address)
}
