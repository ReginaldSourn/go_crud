package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/reginaldsourn/go-crud/internal/adapters/primary/http/dto"
	"github.com/reginaldsourn/go-crud/internal/core/ports"
	pkg "github.com/reginaldsourn/go-crud/pkg/error"
)

type DevicesHandler struct {
	store ports.DeviceStore
}

func NewDevicesHandler(store ports.DeviceStore) *DevicesHandler {
	return &DevicesHandler{store: store}
}

func (h *DevicesHandler) Create(c *gin.Context) {
	var req dto.CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	if req.TypeID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type_id is required"})
		return
	}

	device, err := h.store.Create(c.Request.Context(), req.Name, req.TypeID)
	if err != nil {
		status := http.StatusBadRequest
		if err == pkg.ErrDuplicateDevice {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.ToDeviceResponse(device))
}

func (h *DevicesHandler) Get(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	device, err := h.store.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ToDeviceResponse(device))
}

func (h *DevicesHandler) List(c *gin.Context) {
	devices, err := h.store.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make([]dto.DeviceResponse, 0, len(devices))
	for _, device := range devices {
		resp = append(resp, dto.ToDeviceResponse(device))
	}

	c.JSON(http.StatusOK, resp)
}

func (h *DevicesHandler) Update(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Name == nil && req.TypeID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	name := ""
	if req.Name != nil {
		if *req.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}
		name = *req.Name
	}

	typeID := int64(0)
	if req.TypeID != nil {
		if *req.TypeID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "type_id is required"})
			return
		}
		typeID = *req.TypeID
	}

	device, err := h.store.Update(c.Request.Context(), id, name, typeID)
	if err != nil {
		status := http.StatusBadRequest
		if err == pkg.ErrDeviceNotFound {
			status = http.StatusNotFound
		} else if err == pkg.ErrDuplicateDevice {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ToDeviceResponse(device))
}

func (h *DevicesHandler) Delete(c *gin.Context) {
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
