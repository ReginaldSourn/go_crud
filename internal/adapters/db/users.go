package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"github.com/reginaldsourn/go-crud/internal/core/domain"
	"github.com/reginaldsourn/go-crud/internal/core/ports"
)

type gormUserStore struct {
	db *gorm.DB
}

func NewGormUserStore(db *gorm.DB) ports.UserStore {
	return &gormUserStore{db: db}
}

func (s *gormUserStore) Create(ctx context.Context, username string, passwordHash []byte) (domain.User, error) {
	if username == "" {
		return domain.User{}, ports.ErrInvalidUsername
	}

	u := domain.User{
		Username:     username,
		Email:        username,
		PasswordHash: passwordHash,
	}
	if err := s.db.WithContext(ctx).Create(&u).Error; err != nil {
		if isUniqueViolation(err) {
			return domain.User{}, ports.ErrUsernameExists
		}
		return domain.User{}, err
	}

	return u, nil
}

func (s *gormUserStore) GetByID(ctx context.Context, id int64) (domain.User, error) {
	var u domain.User
	if err := s.db.WithContext(ctx).First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, ports.ErrUserNotFound
		}
		return domain.User{}, err
	}
	return u, nil
}

func (s *gormUserStore) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	var u domain.User
	if err := s.db.WithContext(ctx).Where("username = ?", username).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, ports.ErrUserNotFound
		}
		return domain.User{}, err
	}
	return u, nil
}

func (s *gormUserStore) List(ctx context.Context) ([]domain.User, error) {
	var users []domain.User
	if err := s.db.WithContext(ctx).Order("id").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (s *gormUserStore) Update(ctx context.Context, id int64, username string, passwordHash []byte) (domain.User, error) {
	var u domain.User
	if err := s.db.WithContext(ctx).First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, ports.ErrUserNotFound
		}
		return domain.User{}, err
	}

	if username != "" && username != u.Username {
		u.Username = username
		u.Email = username
	}
	if passwordHash != nil {
		u.PasswordHash = passwordHash
	}

	if err := s.db.WithContext(ctx).Save(&u).Error; err != nil {
		if isUniqueViolation(err) {
			return domain.User{}, ports.ErrUsernameExists
		}
		return domain.User{}, err
	}

	return u, nil
}

func (s *gormUserStore) Delete(ctx context.Context, id int64) error {
	res := s.db.WithContext(ctx).Delete(&domain.User{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ports.ErrUserNotFound
	}
	return nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
