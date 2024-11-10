package compresser_test

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MomsEngineer/urlshortener/internal/adapters/web/compresser"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setup() *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(compresser.CompresserMiddleware())
	router.POST("/", func(c *gin.Context) {
		link, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, "Unable to read request body")
			return
		}
		c.String(http.StatusCreated, string(link))
	})

	return router
}

func TestCompresserMiddleware(t *testing.T) {
	router := setup()

	tests := []struct {
		name           string
		compress       bool
		decompress     bool
		body           []byte
		header         map[string]string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:       "Request Decompression",
			compress:   false,
			decompress: true,
			body:       []byte("https://example.com"),
			header: map[string]string{
				"Content-Encoding": "gzip",
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   "https://example.com",
		},
		{
			name:       "Response Compression",
			compress:   true,
			decompress: false,
			body:       []byte("https://example.com"),
			header: map[string]string{
				"Accept-Encoding": "gzip",
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   "https://example.com",
		},
		{
			name:       "Request Decompression And Response Compression",
			compress:   true,
			decompress: true,
			body:       []byte("https://example.com"),
			header: map[string]string{
				"Accept-Encoding":  "gzip",
				"Content-Encoding": "gzip",
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   "https://example.com",
		},
		{
			name:       "Request Decompression And Response Compression",
			compress:   false,
			decompress: false,
			body:       []byte("https://example.com"),
			header: map[string]string{
				"Content-Encoding": "gzip",
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to unzip the body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			if tt.decompress {
				gz := gzip.NewWriter(&buf)
				_, err := gz.Write([]byte(tt.body))
				require.NoError(t, err)
				err = gz.Close()
				require.NoError(t, err)
			} else {
				_, err := buf.Write([]byte(tt.body))
				require.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/", &buf)

			for key, value := range tt.header {
				req.Header.Set(key, value)
			}

			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			require.Equal(t, tt.expectedStatus, resp.Code)

			var body []byte
			if tt.compress {
				gz, err := gzip.NewReader(resp.Body)
				require.NoError(t, err)
				defer gz.Close()
				body, err = io.ReadAll(gz)
				require.NoError(t, err)
			} else {
				var err error
				body, err = io.ReadAll(resp.Body)
				require.NoError(t, err)
			}

			assert.Contains(t, string(body), tt.expectedBody)
		})
	}
}
