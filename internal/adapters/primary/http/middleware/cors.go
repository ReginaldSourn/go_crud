package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type CORSOptions struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           time.Duration
}

// CORS applies basic CORS headers for browser clients.
func CORS(opts CORSOptions) gin.HandlerFunc {
	allowedOrigins := opts.AllowedOrigins
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{"*"}
	}

	allowedMethods := opts.AllowedMethods
	if len(allowedMethods) == 0 {
		allowedMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	}

	allowedHeaders := opts.AllowedHeaders
	if len(allowedHeaders) == 0 {
		allowedHeaders = []string{"Authorization", "Content-Type"}
	}

	exposedHeaders := opts.ExposedHeaders

	allowCredentials := opts.AllowCredentials
	maxAge := opts.MaxAge

	methods := strings.Join(allowedMethods, ", ")
	headers := strings.Join(allowedHeaders, ", ")
	exposed := strings.Join(exposedHeaders, ", ")

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			c.Next()
			return
		}

		if !isOriginAllowed(origin, allowedOrigins) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		allowOrigin := origin
		if containsOrigin(allowedOrigins, "*") && !allowCredentials {
			allowOrigin = "*"
		}

		c.Header("Access-Control-Allow-Origin", allowOrigin)
		c.Header("Access-Control-Allow-Methods", methods)
		c.Header("Access-Control-Allow-Headers", headers)
		c.Header("Vary", "Origin")

		if exposed != "" {
			c.Header("Access-Control-Expose-Headers", exposed)
		}
		if allowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if maxAge > 0 {
			c.Header("Access-Control-Max-Age", formatSeconds(maxAge))
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func isOriginAllowed(origin string, allowed []string) bool {
	if containsOrigin(allowed, "*") {
		return true
	}
	return containsOrigin(allowed, origin)
}

func containsOrigin(allowed []string, origin string) bool {
	for _, value := range allowed {
		if value == origin {
			return true
		}
	}
	return false
}

func formatSeconds(d time.Duration) string {
	return strconv.FormatInt(int64(d.Seconds()), 10)
}
