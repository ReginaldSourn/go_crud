package ports

import (
	"context"

	"github.com/reginaldsourn/go-crud/internal/core/domain"
)

type UserStore interface {
	Create(ctx context.Context, username string, passwordHash []byte) (domain.User, error)
	GetByID(ctx context.Context, id int64) (domain.User, error)
	GetByUsername(ctx context.Context, username string) (domain.User, error)
	List(ctx context.Context) ([]domain.User, error)
	Update(ctx context.Context, id int64, username string, passwordHash []byte) (domain.User, error)
	Delete(ctx context.Context, id int64) error
}
