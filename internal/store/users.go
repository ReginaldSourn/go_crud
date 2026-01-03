package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"github.com/reginaldsourn/go-crud/internal/model"
)

type gormUserStore struct {
	db *gorm.DB
}

type User struct {
	ID           int64
	Username     string
	PasswordHash []byte
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUsernameExists  = errors.New("username already exists")
	ErrInvalidUsername = errors.New("invalid username")
)

type UserStore interface {
	Create(ctx context.Context, username string, passwordHash []byte) (User, error)
	GetByID(ctx context.Context, id int64) (User, error)
	GetByUsername(ctx context.Context, username string) (User, error)
	List(ctx context.Context) ([]User, error)
	Update(ctx context.Context, id int64, username string, passwordHash []byte) (User, error)
	Delete(ctx context.Context, id int64) error
}

func NewGormUserStore(db *gorm.DB) UserStore {
	return &gormUserStore{db: db}
}

func (s *gormUserStore) Create(ctx context.Context, username string, passwordHash []byte) (User, error) {
	if username == "" {
		return User{}, ErrInvalidUsername
	}

	u := model.User{
		Username:     username,
		Email:        username,
		PasswordHash: passwordHash,
	}
	if err := s.db.WithContext(ctx).Create(&u).Error; err != nil {
		if isUniqueViolation(err) {
			return User{}, ErrUsernameExists
		}
		return User{}, err
	}

	return toStoreUser(u), nil
}

func (s *gormUserStore) GetByID(ctx context.Context, id int64) (User, error) {
	var u model.User
	if err := s.db.WithContext(ctx).First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}
	return toStoreUser(u), nil
}

func (s *gormUserStore) GetByUsername(ctx context.Context, username string) (User, error) {
	var u model.User
	if err := s.db.WithContext(ctx).Where("username = ?", username).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}
	return toStoreUser(u), nil
}

func (s *gormUserStore) List(ctx context.Context) ([]User, error) {
	var users []model.User
	if err := s.db.WithContext(ctx).Order("id").Find(&users).Error; err != nil {
		return nil, err
	}

	out := make([]User, 0, len(users))
	for _, u := range users {
		out = append(out, toStoreUser(u))
	}
	return out, nil
}

func (s *gormUserStore) Update(ctx context.Context, id int64, username string, passwordHash []byte) (User, error) {
	var u model.User
	if err := s.db.WithContext(ctx).First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
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
			return User{}, ErrUsernameExists
		}
		return User{}, err
	}

	return toStoreUser(u), nil
}

func (s *gormUserStore) Delete(ctx context.Context, id int64) error {
	res := s.db.WithContext(ctx).Delete(&model.User{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func toStoreUser(u model.User) User {
	return User{
		ID:           u.ID,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
