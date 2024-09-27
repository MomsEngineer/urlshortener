package main

import (
	"github.com/MomsEngineer/urlshortener/internal/config"
	"github.com/MomsEngineer/urlshortener/internal/db"
	"github.com/MomsEngineer/urlshortener/internal/handlers"
	"github.com/MomsEngineer/urlshortener/internal/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	cfg := config.NewConfig()
	db := db.NewDB()

	router := gin.New()
	router.SetTrustedProxies(nil)

	router.Use(logger.Create().Logger())

	router.POST("/", func(c *gin.Context) {
		handlers.HandlePost(c, db, cfg.BaseURL)
	})

	router.POST("/api/shorten", func(c *gin.Context) {
		handlers.HandlePostAPI(c, db, cfg.BaseURL)
	})

	router.GET("/:id", func(c *gin.Context) {
		handlers.HandleGet(c, db)
	})

	router.Run(cfg.Address)
}
