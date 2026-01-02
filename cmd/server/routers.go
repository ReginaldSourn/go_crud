package main

import (
	"go/version"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	// Create a new Gin router
	router := gin.Default()

	// Define a simple GET endpoint
	router.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})
	v := "v1"

	// Versions endpoint (uses go/version to validate/parse a Go toolchain version string)
	router.GET("/versions", func(c *gin.Context) {
		toolchain := "go1.22.0"
		c.JSON(200, gin.H{
			"api": v,
			"go": gin.H{
				"toolchain": toolchain,
				"valid":     version.IsValid(toolchain),
				"lang":      version.Lang(toolchain),
			},
		})
	})

	// API v1 routes
	api := router.Group("/api/" + v)
	{
		// Register route
		api.POST("/register", func(c *gin.Context) {
			var req struct {
				Username string `json:"username" binding:"required"`
				Password string `json:"password" binding:"required"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			c.JSON(201, gin.H{
				"message":  "registered",
				"username": req.Username,
			})
		})

		// Login route
		api.POST("/login", func(c *gin.Context) {
			var req struct {
				Username string `json:"username" binding:"required"`
				Password string `json:"password" binding:"required"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			// Dummy auth response (replace with real authentication/token issuance)
			c.JSON(200, gin.H{
				"message":  "logged_in",
				"username": req.Username,
				"token":    "dummy-token",
			})
		})

	}

	return router
}
