package whttp

import (
	"time"

	"github.com/edgarsucre/web/token"
)

type TokenManager interface {
	VerifyToken(token string) (token.Claims, error)
	CreateToken(
		username string,
		audience string,
		duration time.Duration,
		content map[string]any,
	) (string, error)
}
