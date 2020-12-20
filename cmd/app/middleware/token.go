package middleware

import (
	"context"
	"errors"
	"net/http"
)

// ErrNoAuthentication ...
var ErrNoAuthentication = errors.New("no authentication")

var authenticationContextKey = &contextKey{"authentication context"}

type contextKey struct {
	name string
}

func (c *contextKey) string() string {
	return c.name
}

var r *http.Request

//IDFunc ...
type IDFunc func(ctx context.Context, token string) (int64, error)

//Authenticate ...
func Authenticate(idFunc IDFunc) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			id, err := idFunc(r.Context(), token)
			if err != nil && r.URL.Path != "/api/managers/token" {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), authenticationContextKey, id)
			r = r.WithContext(ctx)

			handler.ServeHTTP(w, r)
		})
	}
}

//Authentication ...
func Authentication(ctx context.Context) (int64, error) {
	if value, ok := ctx.Value(authenticationContextKey).(int64); ok {
		return value, nil
	}
	return 0, ErrNoAuthentication
}
