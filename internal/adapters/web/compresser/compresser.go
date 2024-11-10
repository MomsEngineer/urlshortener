package compresser

import (
	"compress/gzip"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type compressWriter struct {
	gin.ResponseWriter
	zipWriter *gzip.Writer
}

func (c *compressWriter) Write(b []byte) (int, error) {
	n, err := c.zipWriter.Write(b)
	if err != nil {
		return n, err
	}

	return n, nil
}

func CompresserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		сontentEncoding := c.Request.Header.Get("Content-Encoding")
		if strings.Contains(сontentEncoding, "gzip") {
			gz, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.String(http.StatusInternalServerError, "Failed to unzip the body")
				c.Abort()
				return
			}
			defer gz.Close()

			c.Request.Body = gz
		}

		acceptEncoding := c.Request.Header.Get("Accept-Encoding")
		if strings.Contains(acceptEncoding, "gzip") {
			gz := gzip.NewWriter(c.Writer)
			defer gz.Close()

			c.Header("Content-Encoding", "gzip")

			c.Writer = &compressWriter{
				ResponseWriter: c.Writer,
				zipWriter:      gz,
			}
		}

		c.Next()
	}
}
