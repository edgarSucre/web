package token

import (
	"errors"

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

func (manager *JWTManager) CreateToken(claims Claims) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return jwtToken.SignedString([]byte(manager.secretKey))
}

func (manager *JWTManager) VerifyToken(token string) (Claims, error) {
	keyFunc := func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", ErrInvalidToken
		}

		return []byte(manager.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, keyFunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return Claims{}, ErrExpiredToken
		}

		return Claims{}, ErrInvalidToken
	}

	claims, ok := jwtToken.Claims.(Claims)
	if !ok {
		return Claims{}, ErrInvalidToken
	}

	return claims, nil
}
