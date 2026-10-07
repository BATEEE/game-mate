package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ZapLogger trả về middleware ghi log HTTP request có cấu trúc qua Uber Zap.
func ZapLogger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		fields := []zap.Field{
			zap.Int("status", statusCode),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", clientIP),
			zap.Duration("latency", latency),
		}

		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("error", c.Errors.ByType(gin.ErrorTypePrivate).String()))
			log.Error("HTTP request failed", fields...)
			return
		}

		if statusCode >= 500 {
			log.Error("HTTP server error", fields...)
		} else if statusCode >= 400 {
			log.Warn("HTTP client error", fields...)
		} else {
			log.Info("HTTP request", fields...)
		}
	}
}
