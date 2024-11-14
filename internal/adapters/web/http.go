package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/MomsEngineer/urlshortener/internal/adapters/logger"
	ierrors "github.com/MomsEngineer/urlshortener/internal/errors"
	"github.com/MomsEngineer/urlshortener/internal/usecases/storage"
	"github.com/gin-gonic/gin"
)

var log = logger.Create()

type BatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func getUserIDFromContext(c *gin.Context) (string, error) {
	userID, exists := c.Get("userID")
	if !exists {
		log.Error("User ID not found in context", nil)
		return "", errors.New("user ID not found in context")
	}

	userIDStr, ok := userID.(string)
	if !ok {
		log.Error("User ID is not a string", nil)
		return "", errors.New("user ID is not a string")
	}

	log.Debug("User ID:", userIDStr)

	return userIDStr, nil
}

func HandleGet(c *gin.Context, s storage.StoregeInterface) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	id := c.Param("id")
	link, err := s.GetLink(c.Request.Context(), userID, id)
	if err != nil {
		c.String(http.StatusNotFound, "Link not found")
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, link)
}

func HandleGetUserURL(c *gin.Context, s storage.StoregeInterface, baseURL string) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	links, err := s.GetLinksByUser(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, ierrors.ErrNoContent) {
			log.Debug("Get links by user", err)
			c.Status(http.StatusNoContent)
			return
		}
		log.Error("Failed to get link by user", err)
		c.String(http.StatusInternalServerError, "Failed to get link by user")
		return
	}

	responses := []struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}{}

	for short, original := range links {
		responses = append(responses, struct {
			ShortURL    string `json:"short_url"`
			OriginalURL string `json:"original_url"`
		}{
			ShortURL:    baseURL + "/" + short,
			OriginalURL: original,
		})
	}

	c.JSON(http.StatusOK, responses)

}

func HandlePing(c *gin.Context, s storage.StoregeInterface) {
	if err := s.Ping(c.Request.Context()); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func HandlePost(c *gin.Context, s storage.StoregeInterface, baseURL string) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	link, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "Unable to read request body")
		return
	}

	shortURL, err := s.SaveLink(c.Request.Context(), userID, string(link))
	if err != nil {
		if errors.Is(err, ierrors.ErrDuplicate) {
			log.Error("Error: Duplicate entry for "+string(link), err)
			c.String(http.StatusConflict, baseURL+"/"+shortURL)
			return
		}

		c.String(http.StatusInternalServerError, "Failed to save link")
		return
	}
	c.String(http.StatusCreated, baseURL+"/"+shortURL)
}

func HandlePostAPI(c *gin.Context, s storage.StoregeInterface, baseURL string) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

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

	retCode := http.StatusCreated

	shortURL, err := s.SaveLink(c.Request.Context(), userID, request.URL)
	if err != nil {
		if errors.Is(err, ierrors.ErrDuplicate) {
			log.Error("Error: Duplicate entry for "+string(request.URL), err)
			retCode = http.StatusConflict
		} else {
			c.String(http.StatusInternalServerError, "Failed to save link")
			return
		}
	}

	response := struct {
		Result string `json:"result"`
	}{
		Result: baseURL + "/" + shortURL,
	}

	c.JSON(retCode, response)
}

func HandlePostBatch(c *gin.Context, s storage.StoregeInterface, baseURL string) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	var requests []BatchRequest
	err = json.NewDecoder(c.Request.Body).Decode(&requests)
	if err != nil {
		log.Error("Failed to decode request", err)
		c.String(http.StatusBadRequest, "Failed to decode request")
		return
	}

	links := make(map[string]string)
	for _, r := range requests {
		links[r.CorrelationID] = r.OriginalURL
	}

	if err = s.SaveLinksBatch(c.Request.Context(), userID, links); err != nil {
		log.Error("Failed to save links batch", err)
		c.String(http.StatusInternalServerError, "Failed to save link")
		return
	}

	var responses []BatchResponse
	for id, short := range links {
		responses = append(responses,
			BatchResponse{
				CorrelationID: id,
				ShortURL:      baseURL + "/" + short,
			})
	}

	c.JSON(http.StatusCreated, responses)
}
