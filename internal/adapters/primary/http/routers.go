package http

import (
	"go/version"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/reginaldsourn/go-crud/internal/adapters/auth"
	httphandlers "github.com/reginaldsourn/go-crud/internal/adapters/http/handlers"
	"github.com/reginaldsourn/go-crud/internal/adapters/primary/http/middleware"
	"github.com/reginaldsourn/go-crud/internal/core/ports"
	pkg "github.com/reginaldsourn/go-crud/pkg/error"
)

type RouterDependencies struct {
	UserStore ports.UserStore
	JWTSecret []byte
	JWTTTL    time.Duration
}

func NewRouter(deps RouterDependencies) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logging())
	router.Use(middleware.CORS(middleware.CORSOptions{}))

	userStoreAvailable := deps.UserStore != nil
	var userHandler *httphandlers.UserHandler
	if userStoreAvailable {
		userHandler = httphandlers.NewUserHandler(deps.UserStore)
	}

	ttl := deps.JWTTTL
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}

	router.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	v := "v1"
	router.GET("/versions", func(c *gin.Context) {
		toolchain := "go1.22.0"
		c.JSON(http.StatusOK, gin.H{
			"api": v,
			"go": gin.H{
				"toolchain": toolchain,
				"valid":     version.IsValid(toolchain),
				"lang":      version.Lang(toolchain),
			},
		})
	})

	api := router.Group("/api/" + v)
	{
		api.POST("/register", func(c *gin.Context) {
			if !userStoreAvailable {
				c.JSON(http.StatusServiceUnavailable, gin.H{"error": "user store not configured"})
				return
			}

			var req struct {
				Username string `json:"username" binding:"required"`
				Password string `json:"password" binding:"required"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
				return
			}

			u, err := deps.UserStore.Create(c.Request.Context(), req.Username, passwordHash)
			if err != nil {
				status := http.StatusBadRequest
				if err == pkg.ErrUsernameExists {
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

		api.POST("/login", func(c *gin.Context) {
			if !userStoreAvailable {
				c.JSON(http.StatusServiceUnavailable, gin.H{"error": "user store not configured"})
				return
			}

			var req struct {
				Username string `json:"username" binding:"required"`
				Password string `json:"password" binding:"required"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			u, err := deps.UserStore.GetByUsername(c.Request.Context(), req.Username)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
				return
			}
			if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(req.Password)); err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
				return
			}

			token, err := auth.GenerateToken(u.Username, deps.JWTSecret, ttl)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to issue token"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"token":    token,
				"username": u.Username,
			})
		})

		api.GET("/me", middleware.AuthMiddleware(deps.JWTSecret), func(c *gin.Context) {
			username, _ := c.Get("username")
			c.JSON(http.StatusOK, gin.H{
				"username": username,
			})
		})

		usersAPI := api.Group("/users", middleware.AuthMiddleware(deps.JWTSecret))
		{
			if userStoreAvailable {
				usersAPI.POST("", userHandler.Create)
				usersAPI.GET("", userHandler.List)
				usersAPI.GET("/:id", userHandler.Get)
				usersAPI.PUT("/:id", userHandler.Update)
				usersAPI.DELETE("/:id", userHandler.Delete)
			} else {
				usersAPI.Any("", serviceUnavailable)
				usersAPI.Any("/:id", serviceUnavailable)
			}
		}
	}

	return router
}

func serviceUnavailable(c *gin.Context) {
	c.JSON(http.StatusServiceUnavailable, gin.H{"error": "user store not configured"})
}
