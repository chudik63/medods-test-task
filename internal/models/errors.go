package models

import "errors"

var (
	ErrUserNotFound    = errors.New("user with such id was not found")
	ErrEmptyUserID     = errors.New("user id is empty")
	ErrSessionNotFound = errors.New("session was not found")
)
