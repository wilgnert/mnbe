package lib

import (
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
)

type TokenStub struct {
	tokens []RefreshToken
}

func (t *TokenStub) SaveToken(token string, user_id uuid.UUID, expires_in time.Duration) (RefreshToken, error) {
	tkn := RefreshToken{
		Token: token,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user_id,
		ExpiresAt: time.Now().Add(expires_in),
		RevokedAt: time.Now().Add(expires_in),
	}
	t.tokens = append(t.tokens, tkn)
	return tkn, nil
}

func (t *TokenStub) GetToken(token string) (RefreshToken, error) {
	idx := slices.IndexFunc(t.tokens, func(r RefreshToken) bool {
		return r.Token == token
	})
	if idx == -1 {
		return RefreshToken{}, fmt.Errorf("token not found")
	}
	return t.tokens[idx], nil
}

func (t *TokenStub) RevokeToken(token string) error {
	idx := slices.IndexFunc(t.tokens, func(r RefreshToken) bool {
		return r.Token == token
	})
	if idx == -1 {
		return fmt.Errorf("token not found")
	}
	t.tokens[idx].RevokedAt = time.Now()
	t.tokens[idx].UpdatedAt = time.Now()
	return nil
}
