package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/reginaldsourn/go-crud/internal/adapters/auth"
	"github.com/reginaldsourn/go-crud/internal/adapters/primary/http/dto"
	"github.com/reginaldsourn/go-crud/internal/core/ports"
)

type AuthHandler struct {
	store     ports.UserStore
	jwtSecret []byte
	jwtTTL    time.Duration
}

func NewAuthHandler(store ports.UserStore, secret []byte, ttl time.Duration) *AuthHandler {
	return &AuthHandler{
		store:     store,
		jwtSecret: secret,
		jwtTTL:    ttl,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.store.GetByUsername(c.Request.Context(), req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := auth.GenerateToken(u.Username, h.jwtSecret, h.jwtTTL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to issue token"})
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{
		Token:    token,
		Username: u.Username,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
