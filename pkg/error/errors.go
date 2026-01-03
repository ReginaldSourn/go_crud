package pkg

import "errors"

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidEmail    = errors.New("invalid email")
	ErrDuplicateEmail  = errors.New("email already exists")
	ErrUsernameExists  = errors.New("username already exists")
	ErrInvalidUsername = errors.New("invalid username")
)

var (
	ErrDeviceNotFound   = errors.New("device not found")
	ErrInvalidDeviceID  = errors.New("invalid device ID")
	ErrDuplicateDevice  = errors.New("device already exists")
	ErrInvalidDeviceKey = errors.New("invalid device key")
)
