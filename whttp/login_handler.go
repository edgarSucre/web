package whttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/edgarsucre/web/util"
)

// Consider moving this file to a template project for user and password authentication

type (
	loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	loginResponse struct {
		AccessToken string `json:"access_token"`
	}

	User struct {
		UserName          string
		EncryptedPassword string
	}

	userStore interface {
		GetUser(ctx context.Context, username string) (User, error)
	}
)

const duration = time.Duration(time.Minute * 15)

func HandleLogin(tokenManager TokenManager, userStore userStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest

		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&req); err != nil {
			http.Error(w, "malformed credentials", http.StatusBadRequest)
			return
		}

		user, err := userStore.GetUser(r.Context(), req.Username)
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		if err := util.CheckPassword(req.Password, user.EncryptedPassword); err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		token, err := tokenManager.CreateToken(user.UserName, "user company", duration, nil)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("could not authenticate user: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}

		resp := loginResponse{
			AccessToken: token,
		}

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
