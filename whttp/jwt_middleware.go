package whttp

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/edgarsucre/web/token"
)

type CtxKey int

const (
	ClaimsKey CtxKey = iota
)

type (
	HttpSkipper func(*http.Request) bool

	TokenManager interface {
		VerifyToken(token string) (token.Claims, error)
		CreateToken(claims token.Claims) (string, error)
	}
)

const (
	authHeader   = "Authorization"
	bearerPrefix = "Bearer"
)

func JwtMiddlewareHandler(
	next http.Handler,
	verifier TokenManager,
	skipper HttpSkipper,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if skipper(r) {
			next.ServeHTTP(w, r)
			return
		}

		auth := r.Header.Get(authHeader)

		if len(auth) == 0 {
			err := fmt.Sprintf("missing %s header", authHeader)
			http.Error(w, err, http.StatusUnauthorized)

			return
		}

		if !strings.HasPrefix(auth, bearerPrefix) {
			err := fmt.Sprintf("missing %s token", bearerPrefix)
			http.Error(w, err, http.StatusUnauthorized)
			return
		}

		claims, err := verifier.VerifyToken(auth[len(bearerPrefix)+1:])
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ClaimsKey, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
