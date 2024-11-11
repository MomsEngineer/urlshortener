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

func HandleGet(c *gin.Context, s storage.StoregeInterface) {
	id := c.Param("id")
	link, err := s.GetLink(c.Request.Context(), id)
	if err != nil {
		c.String(http.StatusNotFound, "Link not found")
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, link)
}

func HandlePing(c *gin.Context, s storage.StoregeInterface) {
	if err := s.Ping(c.Request.Context()); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func HandlePost(c *gin.Context, s storage.StoregeInterface, baseURL string) {
	link, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "Unable to read request body")
		return
	}

	shortURL, err := s.SaveLink(c.Request.Context(), string(link))
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

	shortURL, err := s.SaveLink(c.Request.Context(), request.URL)
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
	var requests []BatchRequest
	err := json.NewDecoder(c.Request.Body).Decode(&requests)
	if err != nil {
		log.Error("Failed to decode request", err)
		c.String(http.StatusBadRequest, "Failed to decode request")
		return
	}

	links := make(map[string]string)
	for _, r := range requests {
		links[r.CorrelationID] = r.OriginalURL
	}

	if err = s.SaveLinksBatch(c.Request.Context(), links); err != nil {
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
