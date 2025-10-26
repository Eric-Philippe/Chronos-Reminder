package api

import "errors"

var (
	// Registration validation errors
	ErrEmailRequired     = errors.New("email is required")
	ErrInvalidEmail      = errors.New("invalid email format")
	ErrUsernameRequired  = errors.New("username is required")
	ErrUsernameTooLong   = errors.New("username must be 128 characters or less")
	ErrPasswordRequired  = errors.New("password is required")
	ErrPasswordTooShort  = errors.New("password must be at least 8 characters long")
	ErrTimezoneRequired  = errors.New("timezone is required")
)
