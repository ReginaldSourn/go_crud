package main

import (
	"go/version"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"github.com/reginaldsourn/go-crud/internal/auth"
	"github.com/reginaldsourn/go-crud/internal/handlers"
	"github.com/reginaldsourn/go-crud/internal/middlewares"
	"github.com/reginaldsourn/go-crud/internal/store"
)

func SetupRouter() *gin.Engine {
	// Create a new Gin router
	router := gin.Default()
	// Load environment variables from .env (optional)
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found; using existing environment variables")
	}
	users := store.NewInMemoryUserStore()
	userHandler := handlers.NewUserHandler(users)
	secret := []byte(os.Getenv("JWT_SECRET"))

	log.Println("JWT_SECRET length:", len(secret))
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

			passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
				return
			}

			u, err := users.Create(c.Request.Context(), req.Username, passwordHash)
			if err != nil {
				status := http.StatusBadRequest
				if err == store.ErrUsernameExists {
					status = http.StatusConflict
				}
				c.JSON(status, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusCreated, gin.H{
				"id":       u.ID,
				"username": u.Username,
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

			u, err := users.GetByUsername(c.Request.Context(), req.Username)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
				return
			}
			if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(req.Password)); err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
				return
			}

			token, err := auth.GenerateToken(u.Username, secret, 24*time.Hour)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to issue token"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"token":    token,
				"username": u.Username,
			})
		})

		api.GET("/me", middlewares.AuthMiddleware(secret), func(c *gin.Context) {
			username, _ := c.Get("username")
			c.JSON(http.StatusOK, gin.H{
				"username": username,
			})
		})

		usersAPI := api.Group("/users", middlewares.AuthMiddleware(secret))
		{
			usersAPI.POST("", userHandler.Create)
			usersAPI.GET("", userHandler.List)
			usersAPI.GET("/:id", userHandler.Get)
			usersAPI.PUT("/:id", userHandler.Update)
			usersAPI.DELETE("/:id", userHandler.Delete)
		}
	}

	return router
}
