package middleware

import (
	"log/slog"
	"mini-shop/gateway/logging"
	"time"

	"github.com/gin-gonic/gin"
)

func Logging(logger *logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		ctx := c.Request.Context()

		logger.Info(
			ctx,
			"http request completed",
			slog.String("method", c.Request.Method),
			slog.String("route", c.FullPath()),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("duration", time.Since(start)),
		)
	}
}
