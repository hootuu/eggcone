package rest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hootuu/gelato/logger"
	"go.uber.org/zap"
	"time"
)

type Logger struct {
	logger *zap.Logger
}

func NewLogger(code string) *Logger {
	return &Logger{
		logger: logger.GetLogger(code),
	}
}

func (logger *Logger) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		method := c.Request.Method
		path := c.Request.URL.Path

		c.Next()

		latency := time.Since(startTime)

		statusCode := c.Writer.Status()

		logger.logger.Info(fmt.Sprintf("[%s]%s", method, path),
			zap.Time("t", startTime),
			zap.Int("s", statusCode),
			zap.Duration("e", latency),
		)
	}
}
