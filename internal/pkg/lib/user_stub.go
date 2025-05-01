package lib

import (
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
)

type UserWithHashedPassword struct {
	User
	HashedPassword string `json:"hashed_password"`
}

type UserStub struct {
	users []UserWithHashedPassword
}

func (u *UserStub) CreateUser(params CreateUserParams) (User, error) {
	u.users = append(u.users, UserWithHashedPassword{
		User: User{
			ID: uuid.New(),
			Email: params.Email,
			DisplayName: params.Email,
			Username: params.Email,
			ProfilePictureURL: "",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			IsPremium: false,
			PremiumExpiresAt: time.Now(),
		},
		HashedPassword: params.HashedPassword,
	})
	return u.users[len(u.users) - 1].User, nil
}

func (u *UserStub) RetrieveUsers() ([]User, error) {
	var users []User
	for _, user := range u.users {
		users = append(users, user.User)
	}
	return users, nil
}

func (u *UserStub) RetrieveUserById(id uuid.UUID) (User, error) {
	idx := slices.IndexFunc(u.users, func(user UserWithHashedPassword) bool {
		return user.ID == id
	})
	if idx == -1 {
		return User{}, fmt.Errorf("user not found")
	}
	return u.users[idx].User, nil
}

func (u *UserStub) RetrieveUserByDisplayName(name string) (User, error) {
	idx := slices.IndexFunc(u.users, func(user UserWithHashedPassword) bool {
		return user.DisplayName == name
	})
	if idx == -1 {
		return User{}, fmt.Errorf("user not found")
	}
	return u.users[idx].User, nil
}

func (u *UserStub) RetrieveUserByEmail(email string) (User, error) {
	idx := slices.IndexFunc(u.users, func(user UserWithHashedPassword) bool {
		return user.Email == email
	})
	if idx == -1 {
		return User{}, fmt.Errorf("user not found")
	}
	return u.users[idx].User, nil
}

func (u *UserStub) RetrieveUserByUsername(username string) (User, error) {
	idx := slices.IndexFunc(u.users, func(user UserWithHashedPassword) bool {
		return user.Username == username
	})
	if idx == -1 {
		return User{}, fmt.Errorf("user not found")
	}
	return u.users[idx].User, nil
}

func (u *UserStub) UpdateUser(params UpdateUserParams) (User, error) {
	idx := slices.IndexFunc(u.users, func(user UserWithHashedPassword) bool {
		return user.ID == params.ID
	})
	if idx == -1 {
		return User{}, fmt.Errorf("user not found")
	}
	if params.DisplayName == "" && params.ProfilePictureURL == "" && params.Username == "" {
		return User{}, fmt.Errorf("bad request, please provide a field to update")
	}
	if params.DisplayName != "" {
		u.users[idx].DisplayName = params.DisplayName
	}
	if params.ProfilePictureURL != "" {
		u.users[idx].ProfilePictureURL = params.ProfilePictureURL
	}
	if params.Username != "" {
		u.users[idx].Username = params.Username
	}
	u.users[idx].UpdatedAt = time.Now()
	return u.users[idx].User, nil
}

func (u *UserStub) DeleteUserByID(id uuid.UUID) error {
	l := len(u.users)
	u.users = slices.DeleteFunc(u.users, func(user UserWithHashedPassword) bool {
		return user.ID == id
	})
	if l == len(u.users) {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (u *UserStub) DeleteUserByEmail(email string) error {
	l := len(u.users)
	u.users = slices.DeleteFunc(u.users, func(user UserWithHashedPassword) bool {
		return user.Email == email
	})
	if l == len(u.users) {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (u *UserStub) DeleteUserByUsername(username string) error {
	l := len(u.users)
	u.users = slices.DeleteFunc(u.users, func(user UserWithHashedPassword) bool {
		return user.Username == username
	})
	if l == len(u.users) {
		return fmt.Errorf("user not found")
	}
	return nil
}