package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"mini-shop/gateway/metrics"
)

func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		start := time.Now()

		metrics.RequestsInFlight.Inc()
		defer metrics.RequestsInFlight.Dec()

		c.Next()

		duration := time.Since(start).Seconds()

		metrics.RequestsTotal.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
			strconv.Itoa(c.Writer.Status()),
		).Inc()

		metrics.RequestDuration.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
		).Observe(duration)
	}
}
