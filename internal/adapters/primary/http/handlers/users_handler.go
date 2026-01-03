package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/reginaldsourn/go-crud/internal/adapters/primary/http/dto"
	"github.com/reginaldsourn/go-crud/internal/core/ports"
	pkg "github.com/reginaldsourn/go-crud/pkg/error"
)

type UsersHandler struct {
	store ports.UserStore
}

func NewUsersHandler(store ports.UserStore) *UsersHandler {
	return &UsersHandler{store: store}
}

func (h *UsersHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	u, err := h.store.Create(c.Request.Context(), req.Username, passwordHash)
	if err != nil {
		status := http.StatusBadRequest
		if err == pkg.ErrUsernameExists {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.ToUserResponse(u))
}

func (h *UsersHandler) Get(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	u, err := h.store.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ToUserResponse(u))
}

func (h *UsersHandler) List(c *gin.Context) {
	users, err := h.store.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		resp = append(resp, dto.ToUserResponse(u))
	}

	c.JSON(http.StatusOK, resp)
}

func (h *UsersHandler) Update(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Username == nil && req.Password == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	username := ""
	if req.Username != nil {
		if *req.Username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
			return
		}
		username = *req.Username
	}

	var passwordHash []byte
	if req.Password != nil {
		if *req.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "password is required"})
			return
		}
		passwordHash, err = bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}
	}

	u, err := h.store.Update(c.Request.Context(), id, username, passwordHash)
	if err != nil {
		status := http.StatusBadRequest
		if err == pkg.ErrUserNotFound {
			status = http.StatusNotFound
		} else if err == pkg.ErrUsernameExists {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ToUserResponse(u))
}

func (h *UsersHandler) Delete(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.store.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
