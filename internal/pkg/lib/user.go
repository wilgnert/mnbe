package lib

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `json:"ID"`
	Email string `json:"email"`
	DisplayName string `json:"display_name"`
	Username string `json:"username"`
	ProfilePictureURL string `json:"profile_picture_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsPremium bool `json:"is_premium"`
	PremiumExpiresAt time.Time `json:"premium_expires_at"`
}

type CreateUserParams struct {
	Email string `json:"email"`
	HashedPassword string `json:"hashed_password"`
}

type UpdateUserParams struct {
	ID uuid.UUID `json:"ID"`
	DisplayName string `json:"display_name,omitempty"`
	ProfilePictureURL string `json:"profile_picture_url,omitempty"`
	Username string `json:"username,omitempty"`
}

type UserCrud interface {
	CreateUser(CreateUserParams) (User, error)
	RetrieveUsers() ([]User, error)
	RetrieveUserById(uuid.UUID) (User, error)
	RetrieveUserByDisplayName(string) (User, error)
	RetrieveUserByEmail(string) (User, error)
	RetrieveUserByUsername(string) (User, error)
	UpdateUser(UpdateUserParams) (User, error)
	DeleteUserByID(uuid.UUID) error
	DeleteUserByEmail(string) error
	DeleteUserByUsername(string) error
}