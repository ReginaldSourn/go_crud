package ports

import (
	"context"

	"github.com/reginaldsourn/go-crud/internal/core/domain"
)

type DeviceStore interface {
	Create(ctx context.Context, name string, typeID int64) (domain.Device, error)
	GetByID(ctx context.Context, id int64) (domain.Device, error)
	List(ctx context.Context) ([]domain.Device, error)
	Update(ctx context.Context, id int64, name string, typeID int64) (domain.Device, error)
	Delete(ctx context.Context, id int64) error
}
