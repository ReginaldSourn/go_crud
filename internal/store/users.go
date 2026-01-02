package store

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUsernameExists  = errors.New("username already exists")
	ErrInvalidUsername = errors.New("invalid username")
)

type User struct {
	ID           int64
	Username     string
	PasswordHash []byte
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserStore interface {
	Create(ctx context.Context, username string, passwordHash []byte) (User, error)
	GetByID(ctx context.Context, id int64) (User, error)
	GetByUsername(ctx context.Context, username string) (User, error)
	List(ctx context.Context) ([]User, error)
	Update(ctx context.Context, id int64, username string, passwordHash []byte) (User, error)
	Delete(ctx context.Context, id int64) error
}

type inMemoryUserStore struct {
	mu     sync.RWMutex
	nextID int64
	byID   map[int64]User
	byName map[string]int64
}

func NewInMemoryUserStore() UserStore {
	return &inMemoryUserStore{
		nextID: 1,
		byID:   make(map[int64]User),
		byName: make(map[string]int64),
	}
}

func (s *inMemoryUserStore) Create(ctx context.Context, username string, passwordHash []byte) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if username == "" {
		return User{}, ErrInvalidUsername
	}
	if _, exists := s.byName[username]; exists {
		return User{}, ErrUsernameExists
	}

	now := time.Now().UTC()
	u := User{
		ID:           s.nextID,
		Username:     username,
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	s.nextID++
	s.byID[u.ID] = u
	s.byName[username] = u.ID
	return u, nil
}

func (s *inMemoryUserStore) GetByID(ctx context.Context, id int64) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	u, ok := s.byID[id]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return u, nil
}

func (s *inMemoryUserStore) GetByUsername(ctx context.Context, username string) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	id, ok := s.byName[username]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return s.byID[id], nil
}

func (s *inMemoryUserStore) List(ctx context.Context) ([]User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]User, 0, len(s.byID))
	for _, u := range s.byID {
		users = append(users, u)
	}
	return users, nil
}

func (s *inMemoryUserStore) Update(ctx context.Context, id int64, username string, passwordHash []byte) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	u, ok := s.byID[id]
	if !ok {
		return User{}, ErrUserNotFound
	}

	if username != "" && username != u.Username {
		if _, exists := s.byName[username]; exists {
			return User{}, ErrUsernameExists
		}
		delete(s.byName, u.Username)
		u.Username = username
		s.byName[username] = id
	}

	if passwordHash != nil {
		u.PasswordHash = passwordHash
	}

	u.UpdatedAt = time.Now().UTC()
	s.byID[id] = u
	return u, nil
}

func (s *inMemoryUserStore) Delete(ctx context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	u, ok := s.byID[id]
	if !ok {
		return ErrUserNotFound
	}
	delete(s.byID, id)
	delete(s.byName, u.Username)
	return nil
}
