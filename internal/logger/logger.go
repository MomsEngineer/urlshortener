package logger

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Logger struct {
	logger *zap.Logger
}

func Create() Logger {
	return Logger{
		logger: zap.NewExample(),
	}
}

func (l Logger) Logger() gin.HandlerFunc {
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

func (l Logger) Error(msg string, err error) {
	line, fn := getLineAndFileName()
	t := fmt.Sprintf("line %d: %s(): %s | error: %v", line, fn, msg, err)

	l.logger.Error(t)
}

func (l Logger) Debug(values ...any) {
	msg := strings.Join(convertToStringSlice(values...), " ")
	line, fn := getLineAndFileName()
	t := fmt.Sprintf("line %d: %s(): %s", line, fn, msg)

	l.logger.Debug(t)
}

func getLineAndFileName() (int, string) {
	pc, _, line, _ := runtime.Caller(2)
	fn := runtime.FuncForPC(pc)
	return line, filepath.Base(fn.Name())
}

func convertToStringSlice(values ...any) []string {
	var result []string
	for _, v := range values {
		result = append(result, fmt.Sprintf("%v", v))
	}
	return result
}
