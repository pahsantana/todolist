package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	logRequest = "request"
	logMethod  = "method"
	logPath    = "path"
	logStatus  = "status"
	logLatency = "latency"
)

func RequestLogger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Info(logRequest,
			zap.String(logMethod, c.Request.Method),
			zap.String(logPath, c.Request.URL.Path),
			zap.Int(logStatus, c.Writer.Status()),
			zap.Duration(logLatency, time.Since(start)),
		)
	}
}
