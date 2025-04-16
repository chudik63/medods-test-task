package models

import (
	"time"
)

type RefreshToken struct {
	ID        uint
	UserID    uint64
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time
}
