package services

import (
	"github.com/google/uuid"
	"github.com/wilgnert/mnbe/internal/pkg/auth"
	"github.com/wilgnert/mnbe/internal/pkg/lib"
)

const (
	Nanosecond  = 1
	Microsecond = 1000 * Nanosecond
	Milisecond  = 1000 * Microsecond
	Second      = 1000 * Milisecond
	Minute      = 60 * Second
	Hour        = 60 * Minute
)
const (
	JWTDuration          = 2 * Hour
	RefreshTokenDuration = 90 * 24 * Hour
)

type ServiceError struct {
	msg string
}

func (s *ServiceError) Error() string {
	return s.msg
}

func LoginWithEmail(email, password, secret string, db lib.Database) (statusCode int, payload any, err error) {
	var user lib.User
	var jwt string
	var refresh string

	if user, err = db.RetrieveUserByEmailAndPassword(email, password); err != nil {
		return 401, nil, &ServiceError{msg: "incorrect login information"}
	}

	if jwt, err = auth.MakeJWT(user.ID, secret, JWTDuration); err != nil {
		return 500, nil, &ServiceError{msg: "failed to generate access token"}
	}

	if refresh, err = auth.MakeRefreshToken(); err != nil {
		return 500, nil, &ServiceError{msg: "failed to generate refresh token"}
	}

	if _, err = db.SaveToken(refresh, user.ID, RefreshTokenDuration); err != nil {
		return 500, nil, &ServiceError{msg: "failed to save refresh token"}
	}

	type UserWithTokens struct {
		lib.User
		AccessToken          string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	payload = UserWithTokens{
		User:         user,
		AccessToken:  jwt,
		RefreshToken: refresh,
	}
	return 200, payload, nil

}

func GetUserById(id string, db lib.Database) (statusCode int, payload any, err error) {
	parsed_id, err := uuid.Parse(id)
	if err != nil {
		return 404, nil, &ServiceError{msg: "user not found with given id"}
	}
	user, err := db.RetrieveUserById(parsed_id)
	if err != nil {
		return 404, nil, &ServiceError{msg: "user not found with given id"}
	}

	return 200, user, nil
}

func GetAllUsers(page, limit, offset int, sortMethod string, db lib.Database) (statusCode int, payload any, err error) {
	allUsers, err := db.RetrieveUsers(page, limit, offset, sortMethod)
	if err != nil {
		return 500, nil, err
	}
	payload = map[string]any{
		"data":     allUsers,
		"rowCount": len(allUsers),
	}
	return 200, payload, nil
}

func CreateUserWithEmail(email, password string, db lib.Database) (statusCode int, payload any, err error) {
	hashed_password, _ := auth.HashPassword(password)
	user, err := db.CreateUser(lib.CreateUserParams{
		Email:          email,
		HashedPassword: hashed_password,
	})
	if err != nil {
		return 400, nil, err
	}
	return 201, user, nil
}
