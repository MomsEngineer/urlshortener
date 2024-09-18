package main

import (
	"github.com/MomsEngineer/urlshortener/internal/app"
	"github.com/MomsEngineer/urlshortener/internal/db"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	db := db.NewDB()

	router := gin.New()
	router.SetTrustedProxies(nil)

	router.POST("/", func(c *gin.Context) {
		app.HandlePost(c, db)
	})

	router.GET("/:id", func(c *gin.Context) {
		app.HandleGet(c, db)
	})

	router.Run("localhost:8080")
}
