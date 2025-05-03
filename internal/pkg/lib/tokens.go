package lib

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	Token     string       `json:"token"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	UserID    uuid.UUID    `json:"user_id"`
	ExpiresAt time.Time    `json:"expires_at"`
	RevokedAt time.Time    `json:"revoked_at"`
}

type RefreshTokenDatabase interface {
	SaveToken(string, uuid.UUID, time.Duration) (RefreshToken, error)
	GetToken(string) (RefreshToken, error)
	RevokeToken(string) error
}