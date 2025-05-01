package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	jwt.RegisteredClaims
	content   map[string]any
	tokenType string
}

func (c Claims) valid(tokenType string) error {
	if c.tokenType != tokenType {
		return ErrInvalidToken
	}

	if time.Now().After(c.ExpiresAt.Time) {
		return ErrExpiredToken
	}

	return nil
}

func (manager JWTManager) createClaims(
	username string,
	audience string,
	duration time.Duration,
	content map[string]any,
	tokenType string,
) (Claims, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return Claims{}, fmt.Errorf("could not generate token id: %w", err)
	}

	return Claims{
		content: content,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    manager.issuer,
			Subject:   username,
			Audience:  []string{audience},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        tokenID.String(),
		},
	}, nil
}
