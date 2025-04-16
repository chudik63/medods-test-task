package models

import (
	"time"
)

type RefreshSession struct {
	ID        uint
	UserID    uint64
	IP        string
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time
}
