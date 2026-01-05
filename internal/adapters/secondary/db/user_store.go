package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"github.com/reginaldsourn/go-crud/internal/core/domain"
	pkg "github.com/reginaldsourn/go-crud/pkg/error"
)

type GormUserStore struct {
	db *gorm.DB
}

func NewGormUserStore(db *gorm.DB) *GormUserStore {
	return &GormUserStore{db: db}
}

func (s *GormUserStore) Create(ctx context.Context, username string, passwordHash []byte) (domain.User, error) {
	if username == "" {
		return domain.User{}, pkg.ErrInvalidUsername
	}

	user := domain.User{
		Username:     username,
		PasswordHash: passwordHash,
	}
	if err := s.db.WithContext(ctx).Create(&user).Error; err != nil {
		if isDuplicateErr(err) {
			return domain.User{}, pkg.ErrUsernameExists
		}
		return domain.User{}, err
	}

	return user, nil
}

func (s *GormUserStore) GetByID(ctx context.Context, id int64) (domain.User, error) {
	var user domain.User
	if err := s.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, pkg.ErrUserNotFound
		}
		return domain.User{}, err
	}

	return user, nil
}

func (s *GormUserStore) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	var user domain.User
	if err := s.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, pkg.ErrUserNotFound
		}
		return domain.User{}, err
	}

	return user, nil
}

func (s *GormUserStore) List(ctx context.Context) ([]domain.User, error) {
	var users []domain.User
	if err := s.db.WithContext(ctx).Order("id ASC").Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (s *GormUserStore) Update(ctx context.Context, id int64, username string, passwordHash []byte) (domain.User, error) {
	var user domain.User
	if err := s.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, pkg.ErrUserNotFound
		}
		return domain.User{}, err
	}

	if username != "" {
		user.Username = username
	}
	if len(passwordHash) > 0 {
		user.PasswordHash = passwordHash
	}

	if err := s.db.WithContext(ctx).Save(&user).Error; err != nil {
		if isDuplicateErr(err) {
			return domain.User{}, pkg.ErrUsernameExists
		}
		return domain.User{}, err
	}

	return user, nil
}

func (s *GormUserStore) Delete(ctx context.Context, id int64) error {
	tx := s.db.WithContext(ctx).Delete(&domain.User{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return pkg.ErrUserNotFound
	}

	return nil
}

func isDuplicateErr(err error) bool {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	return false
}
