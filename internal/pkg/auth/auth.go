package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func CheckPasswordHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password))
}

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	return GetFromAuthorization(authHeader, "ApiKey ")

}

func GetFromAuthorization(authorization, prefix string) (string, error) {
	if authorization == "" {
		return "", fmt.Errorf("missing Authorization header")
	}
	if len(authorization) < len(prefix) || authorization[:len(prefix)] != prefix {
		return "", fmt.Errorf("invalid Authorization header format")
	}
	token := authorization[len(prefix):]
	if token == "" {
		return "", fmt.Errorf("missing token in Authorization header")
	}
	return token, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "mnbe", 
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()), 
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject: userID.String(),
	})
	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	var c jwt.RegisteredClaims
	tkn, err := jwt.ParseWithClaims(tokenString, &c, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil { return uuid.UUID{}, err }
	id, err := tkn.Claims.GetSubject()
	if err != nil { return uuid.UUID{}, err }
	return uuid.Parse(id)
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	return GetFromAuthorization(authHeader, "Bearer ")
}

func MakeRefreshToken() (string, error) {
	var b [32]byte
	_, err := rand.Read(b[:])
	if err != nil {
		return "", fmt.Errorf("could not generate random bytes: %w", err)
	}
	return hex.EncodeToString(b[:]), nil
}