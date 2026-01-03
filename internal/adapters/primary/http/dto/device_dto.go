package dto

import (
	"time"

	"github.com/reginaldsourn/go-crud/internal/core/domain"
)

type CreateDeviceRequest struct {
	Name   string `json:"name" binding:"required"`
	TypeID int64  `json:"type_id" binding:"required"`
}

type UpdateDeviceRequest struct {
	Name   *string `json:"name"`
	TypeID *int64  `json:"type_id"`
}

type DeviceResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	TypeID    int64  `json:"type_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func ToDeviceResponse(d domain.Device) DeviceResponse {
	return DeviceResponse{
		ID:        d.ID,
		Name:      d.Name,
		TypeID:    d.TypeID,
		CreatedAt: d.CreatedAt.Format(time.RFC3339),
		UpdatedAt: d.UpdatedAt.Format(time.RFC3339),
	}
}
