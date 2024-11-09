package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	ierrors "github.com/MomsEngineer/urlshortener/internal/errors"
	"github.com/MomsEngineer/urlshortener/internal/logger"
	"github.com/MomsEngineer/urlshortener/internal/storage"
	"github.com/MomsEngineer/urlshortener/internal/utils"
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

func saveLinkToStorage(ctx context.Context, ls storage.Storage, baseURL, link string) (string, error) {
	short, err := utils.GenerateID(8)
	if err != nil {
		return "", err
	}

	oldShort, err := ls.SaveLink(ctx, short, link)
	if err != nil {
		short = oldShort
	}
	shortURL := baseURL + "/" + short

	return shortURL, err
}

func HandlePost(c *gin.Context, ls storage.Storage, baseURL string) {
	link, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "Unable to read request body")
		return
	}

	shortURL, err := saveLinkToStorage(c.Request.Context(), ls, baseURL, string(link))
	if err != nil {
		if errors.Is(err, ierrors.ErrDuplicate) {
			log.Error("Error: Duplicate entry for "+string(link), err)
			c.String(http.StatusConflict, shortURL)
			return
		}

		c.String(http.StatusInternalServerError, "Failed to save link")
		return
	}
	c.String(http.StatusCreated, shortURL)
}

func HandleGet(c *gin.Context, ls storage.Storage) {
	id := c.Param("id")
	link, exists, err := ls.GetLink(c.Request.Context(), id)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	} else if !exists {
		c.String(http.StatusNotFound, "Link not found")
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, link)
}

func HandlePing(c *gin.Context, ls storage.Storage) {
	if err := ls.Ping(c.Request.Context()); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func HandlePostAPI(c *gin.Context, ls storage.Storage, baseURL string) {
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

	shortURL, err := saveLinkToStorage(c.Request.Context(), ls, baseURL, request.URL)
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
		Result: shortURL,
	}

	c.JSON(retCode, response)
}

func HandlePostBatch(c *gin.Context, ls storage.Storage, baseURL string) {
	var requests []BatchRequest
	err := json.NewDecoder(c.Request.Body).Decode(&requests)
	if err != nil {
		log.Error("Failed to decode request", err)
		c.String(http.StatusBadRequest, "Failed to decode request")
		return
	}

	links := make(map[string]string)
	var responses []BatchResponse
	for _, r := range requests {
		short, err := utils.GenerateID(8)
		if err != nil {
			log.Error("Failed to generate short link", err)
			c.String(http.StatusInternalServerError, "Failed to generate link")
			return
		}
		links[short] = r.OriginalURL

		responses = append(responses,
			BatchResponse{
				CorrelationID: r.CorrelationID,
				ShortURL:      baseURL + "/" + short,
			})
	}

	if err = ls.SaveLinksBatch(c.Request.Context(), links); err != nil {
		log.Error("Failed to save link", err)
		c.String(http.StatusInternalServerError, "Failed to save link")
		return
	}

	c.JSON(http.StatusCreated, responses)
}
