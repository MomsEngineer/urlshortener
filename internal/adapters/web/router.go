package web

import (
	"github.com/MomsEngineer/urlshortener/internal/adapters/logger"
	"github.com/MomsEngineer/urlshortener/internal/adapters/web/compresser"
	"github.com/MomsEngineer/urlshortener/internal/adapters/web/cookie"
	"github.com/MomsEngineer/urlshortener/internal/usecases/storage"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	return gin.New()
}

func SetupRoutes(router *gin.Engine, s storage.StoregeInterface, baseURL string) {
	router.Use(logger.Create(logger.InfoLevel).Logger())
	router.Use(compresser.CompresserMiddleware())
	router.Use(cookie.CookieMiddleware())

	router.POST("/", func(c *gin.Context) {
		HandlePost(c, s, baseURL)
	})

	router.POST("/api/shorten", func(c *gin.Context) {
		HandlePostAPI(c, s, baseURL)
	})

	router.POST("/api/shorten/batch", func(c *gin.Context) {
		HandlePostBatch(c, s, baseURL)
	})

	router.GET("/:id", func(c *gin.Context) {
		HandleGet(c, s)
	})

	router.GET("/api/user/urls", func(c *gin.Context) {
		HandleGetUserURL(c, s, baseURL)
	})

	router.GET("/ping", func(c *gin.Context) {
		HandlePing(c, s)
	})
}
