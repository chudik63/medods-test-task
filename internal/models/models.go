package models

import (
	"time"
)

type RefreshSession struct {
	ID        uint
	UserID    string
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type User struct {
	ID string
	IP string
}
