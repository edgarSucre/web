package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (manager JWTManager) newClaims(username, audience string, duration time.Duration) (jwt.RegisteredClaims, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return jwt.RegisteredClaims{}, fmt.Errorf("unable to generate uuid: %w", err)
	}

	return jwt.RegisteredClaims{
		Issuer:    manager.issuer,
		Subject:   username,
		Audience:  []string{audience},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        tokenID.String(),
	}, nil
}
