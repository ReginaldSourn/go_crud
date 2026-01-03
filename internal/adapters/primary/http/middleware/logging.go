package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logging logs basic request details with latency and status.
func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery
		if rawQuery != "" {
			path = path + "?" + rawQuery
		}

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		log.Printf("%s %s %d %s %s", method, path, status, latency, clientIP)
	}
}
