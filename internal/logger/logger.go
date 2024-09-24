package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GinLogger struct {
	logger *zap.Logger
}

func Create() GinLogger {
	return GinLogger{
		logger: zap.NewExample(),
	}
}

func (l GinLogger) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		uri := c.Request.URL.Path

		method := c.Request.Method

		c.Next()

		duration := time.Since(start)

		status := c.Writer.Status()

		size := c.Writer.Size()

		l.logger.Info("Request and Response:",
			zap.String("url", uri),
			zap.String("method", method),
			zap.Duration("duration", duration),
			zap.Int("status", status),
			zap.Int("size", size))
	}
}
