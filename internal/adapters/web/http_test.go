package web

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MomsEngineer/urlshortener/internal/adapters/web/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setup() *gin.Engine {
	gin.SetMode(gin.TestMode)

	mockStorage := new(mocks.Storage)

	router := gin.New()
	router.POST("/", func(c *gin.Context) {
		HandlePost(c, mockStorage, "http://localhost:8080/")
	})
	router.POST("/api/shorten", func(c *gin.Context) {
		HandlePostAPI(c, mockStorage, "http://localhost:8080/")
	})
	router.GET("/:id", func(c *gin.Context) {
		HandleGet(c, mockStorage)
	})

	return router
}

func TestHandler(t *testing.T) {
	router := setup()

	tests := []struct {
		name           string
		method         string
		url            string
		body           []byte
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "POST request",
			method:         http.MethodPost,
			url:            "/",
			body:           []byte("https://example.com"),
			expectedStatus: http.StatusCreated,
			expectedBody:   "http://localhost:8080/",
		},
		{
			name:           "POST request",
			method:         http.MethodPost,
			url:            "/",
			body:           []byte("https://example.com"),
			expectedStatus: http.StatusCreated,
			expectedBody:   "http://localhost:8080/",
		},
		{
			name:           "POST api request",
			method:         http.MethodPost,
			url:            "/api/shorten",
			body:           []byte("{\"url\":\"https://example.com\"}"),
			expectedStatus: http.StatusCreated,
			expectedBody:   "{\"result\":\"http://localhost:8080/",
		},
		{
			name:           "Successful GET request",
			method:         http.MethodGet,
			url:            "/abc123",
			body:           nil,
			expectedStatus: http.StatusTemporaryRedirect,
			expectedBody:   "",
		},
		{
			name:           "Failed GET request",
			method:         http.MethodGet,
			url:            "/abc",
			body:           nil,
			expectedStatus: http.StatusNotFound,
			expectedBody:   "",
		},
		{
			name:           "Invalid method",
			method:         http.MethodPut,
			url:            "/",
			body:           nil,
			expectedStatus: http.StatusNotFound,
			expectedBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, bytes.NewBuffer(tt.body))
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, rr.Body.String(), tt.expectedBody)
			}
		})
	}
}
