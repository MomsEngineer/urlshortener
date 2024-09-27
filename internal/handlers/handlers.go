package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/MomsEngineer/urlshortener/internal/db"
	"github.com/MomsEngineer/urlshortener/internal/utils"
	"github.com/gin-gonic/gin"
)

func saveLinkToDatabase(database db.Database, baseURL, link string) (string, error) {
	id, err := utils.GenerateID(8)
	if err != nil {
		return "", err
	}

	database.SaveLink(id, link)
	shortURL := baseURL + "/" + id

	return shortURL, nil
}

func HandlePost(c *gin.Context, database db.Database, BaseURL string) {
	link, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "Unable to read request body")
		return
	}

	shortURL, err := saveLinkToDatabase(database, BaseURL, string(link))
	if err != nil {
		c.String(http.StatusInternalServerError, "Could not convert the link.")
	}
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

func HandlePostAPI(c *gin.Context, database db.Database, BaseURL string) {
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(c.Request.Body); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	request := struct {
		URL string `json:"url"`
	}{}

	if err := json.Unmarshal(buf.Bytes(), &request); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL, err := saveLinkToDatabase(database, BaseURL, request.URL)
	if err != nil {
		c.String(http.StatusInternalServerError, "Could not convert the link.")
	}

	response := struct {
		Result string `json:"result"`
	}{
		Result: shortURL,
	}

	resp, err := json.Marshal(response)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Data(http.StatusCreated, "application/json", resp)
}
