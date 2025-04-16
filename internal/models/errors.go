package models

import "errors"

var (
	ErrUserNotFound     = errors.New("user with such id was not found")
	ErrEmptyUserID      = errors.New("user id is empty")
	ErrSessionNotFound  = errors.New("session was not found")
	ErrInvalidToken     = errors.New("token is invalid")
	ErrInvalidTokenType = errors.New("token type is invalid. Required refresh token")
	ErrTokenExpired     = errors.New("token is expired")
	ErrParseToken       = errors.New("failed to parse token")

	ErrSMTPEmptyTo        = errors.New("empty to address")
	ErrSMTPEmptyMail      = errors.New("empty subject or body")
	ErrSMTPInvalidToEmail = errors.New("invalid to email")
	ErrEmailFormat        = errors.New("wrong email format")
)
