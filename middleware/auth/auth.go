//go:build !solution

package auth

import (
	"strings"
	"context"
	"errors"
	"net/http"
)

type User struct {
	Name  string
	Email string
}

func ContextUser(ctx context.Context) (*User, bool) {
	user := ctx.Value("user")
	if user == nil {
		return nil, false
	}
	return user.(*User), true
}

var ErrInvalidToken = errors.New("invalid token")

type TokenChecker interface {
	CheckToken(ctx context.Context, token string) (*User, error)
}

func CheckAuth(checker TokenChecker) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func (w http.ResponseWriter, r *http.Request){
				token, ok :=  strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
				if !ok {
					unauthorized(w)
					return 
				}
				
				user, err := checker.CheckToken(r.Context(), token)
				switch {
				case errors.Is(err, ErrInvalidToken):
					unauthorized(w)
				case err != nil:
					http.Error(w, err.Error(), http.StatusInternalServerError)
				default:
					ctx := context.WithValue(r.Context(), "user", user)
					next.ServeHTTP(w, r.WithContext(ctx))
				}
			})
	}
}

func unauthorized(w http.ResponseWriter) {
	http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)	
}
