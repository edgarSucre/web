package whttp

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/edgarsucre/web"
)

type (
	HttpSkipper func(*http.Request) bool
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

		ctx, err := jwtCheck(r, verifier)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func JwtMiddlewareHandlerFunc(
	next http.HandlerFunc,
	verifier TokenManager,
	skipper HttpSkipper,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if skipper(r) {
			next(w, r)
		}

		ctx, err := jwtCheck(r, verifier)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}

		next(w, r.WithContext(ctx))
	}
}

var (
	ErrNoAuthHeader   = fmt.Errorf("%s header is required", authHeader)
	ErrNoBearerPrefix = fmt.Errorf("missing %s token", bearerPrefix)
)

func jwtCheck(
	r *http.Request,
	verifier TokenManager,
) (context.Context, error) {
	auth := r.Header.Get(authHeader)

	if len(auth) == 0 {
		return nil, ErrNoAuthHeader
	}

	if !strings.HasPrefix(auth, bearerPrefix) {
		return nil, ErrNoBearerPrefix
	}

	claims, err := verifier.VerifyToken(auth[len(bearerPrefix)+1:])
	if err != nil {
		return nil, err
	}

	return context.WithValue(r.Context(), web.ClaimsKey, claims), nil
}
