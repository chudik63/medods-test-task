package models

import "errors"

var (
	ErrUserNotFound           = errors.New("user with such id was not found")
	ErrEmptyUserID            = errors.New("user id is empty")
	ErrInvalidUserID          = errors.New("user id should be a valid uuid")
	ErrSessionNotFound        = errors.New("session was not found")
	ErrInvalidSession         = errors.New("refresh session is invalid")
	ErrTokenExpired           = errors.New("token is expired")
	ErrInvalidToken           = errors.New("token is invalid")
	ErrMismatchedHashAndToken = errors.New("token does not match with the hash")

	ErrSMTPEmptyTo        = errors.New("empty to address")
	ErrSMTPEmptyMail      = errors.New("empty subject or body")
	ErrSMTPInvalidToEmail = errors.New("invalid to email")
	ErrEmailFormat        = errors.New("wrong email format")
)
