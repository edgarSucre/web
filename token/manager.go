package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrExpiredToken     = errors.New("token has expired")
	ErrInvalidToken     = errors.New("token is invalid")
	ErrInvalidSecretKey = errors.New("secret key is too short")
)

const (
	minSecretKeySize = 32
)

type JWTManager struct {
	issuer    string
	secretKey string
}

func NewJWTMaker(secretKey, issuer string) (*JWTManager, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, ErrInvalidSecretKey
	}

	return &JWTManager{
		issuer:    issuer,
		secretKey: secretKey,
	}, nil
}

func (manager *JWTManager) CreateToken(username, audience string, duration time.Duration) (string, error) {
	claims, err := manager.newClaims(username, audience, duration)
	if err != nil {
		return "", fmt.Errorf("unable to create token claims: %w", err)
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return jwtToken.SignedString([]byte(manager.secretKey))
}

func (manager *JWTManager) VerifyToken(token string) (string, error) {
	keyFunc := func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", ErrInvalidToken
		}

		return []byte(manager.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, keyFunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", ErrExpiredToken
		}

		return "", ErrInvalidToken
	}

	claims, ok := jwtToken.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return "", ErrInvalidToken
	}

	return claims.Subject, nil
}
