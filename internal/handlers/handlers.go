package handlers

import (
	"io"
	"net/http"

	"github.com/MomsEngineer/urlshortener/internal/db"
	"github.com/MomsEngineer/urlshortener/internal/utils"
	"github.com/gin-gonic/gin"
)

func HandlePost(c *gin.Context, database db.Database, BaseURL string) {
	link, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "Unable to read request body")
		return
	}

	id, err := utils.GenerateID(8)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to generate short link!")
		return
	}

	database.SaveLink(id, string(link))
	shortURL := BaseURL + "/" + id

	c.String(http.StatusCreated, shortURL)
}

func HandleGet(c *gin.Context, database db.Database) {
	id := c.Param("id")
	link, exists := database.GetLink(id)
	if !exists {
		c.String(http.StatusNotFound, "Link not found")
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, link)
}
