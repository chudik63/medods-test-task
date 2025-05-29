package models

import (
	"time"

	"github.com/google/uuid"
)

type RefreshSession struct {
	ID        uint
	UserID    uuid.UUID
	IP        string
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type User struct {
	ID    uuid.UUID
	Email string
}
