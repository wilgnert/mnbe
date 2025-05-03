package lib

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wilgnert/mnbe/internal/pkg/auth"
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

func (u *UserStub) RetrieveUsers(page, limit, offset int, sort string) ([]User, error) {
	lower := offset + (page - 1) * limit
	if lower > len(u.users) {
		return nil, nil
	}
	higher := lower + limit
	if higher > len(u.users) {
		higher = len(u.users)
	}
	var users []User
	for _, user := range u.users[lower:higher] {
		users = append(users, user.User)
	}
	switch sort {
	case "newest":
		slices.SortFunc(users, func(u1, u2 User) int {
			if u1.CreatedAt.Before(u2.CreatedAt) {
				return -1
			}
			if u1.CreatedAt.After(u2.CreatedAt) {
				return 1
			}
			return 0
		})
	case "oldest":
		slices.SortFunc(users, func(u1, u2 User) int {
			if u1.CreatedAt.After(u2.CreatedAt) {
				return -1
			}
			if u1.CreatedAt.Before(u2.CreatedAt) {
				return 1
			}
			return 0
		})
	case "lexical_asc":
		slices.SortFunc(users, func(u1, u2 User) int {
			return strings.Compare(u1.Username, u2.Username)
		})
	case "lexical_desc":
		slices.SortFunc(users, func(u1, u2 User) int {
			return strings.Compare(u2.Username, u1.Username)
		})
	default:
		return []User{}, fmt.Errorf("unknown sort type")
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

func (u *UserStub) RetrieveUserByEmailAndPassword(email, password string) (User, error) {
	idx := slices.IndexFunc(u.users, func(user UserWithHashedPassword) bool {
		err := auth.CheckPasswordHash(user.HashedPassword, password)
		return user.Email == email && err == nil
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

func (u *UserStub) RetrieveUserByUsernameAndPassword(username, password string) (User, error) {
	idx := slices.IndexFunc(u.users, func(user UserWithHashedPassword) bool {
		err := auth.CheckPasswordHash(user.HashedPassword, password)
		return user.Username == username && err == nil
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