package middlewares

import (
	"net/http"
)

// DevicesMiddleware is a middleware that checks if the device ID is present in the request.
func DevicesMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the device ID is present in the request
		deviceID := r.Header.Get("Device-ID")
		if deviceID == "" {
			http.Error(w, "Device ID is required", http.StatusBadRequest)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
