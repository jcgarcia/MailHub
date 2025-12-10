package middleware

import (
	"context"
	"net/http"

	"github.com/Ingasti/mailhub-admin/internal/config"
)

// AuthUser represents the authenticated user
type AuthUser struct {
	Email string
	Name  string
}

// contextKey for storing auth user in context
type contextKey string

const AuthUserKey contextKey = "authUser"

// Auth middleware extracts authentication headers from Caddy
func Auth(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var user AuthUser

			if cfg.DevMode && cfg.DevAuthEmail != "" {
				// Dev mode - use simulated user
				user = AuthUser{
					Email: cfg.DevAuthEmail,
					Name:  "Dev User",
				}
			} else {
				// Production - extract from Caddy headers
				email := r.Header.Get("X-Auth-User")
				name := r.Header.Get("X-Auth-Name")

				if email == "" {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}

				user = AuthUser{
					Email: email,
					Name:  name,
				}
			}

			// Add user to request context
			ctx := context.WithValue(r.Context(), AuthUserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetAuthUser extracts authenticated user from context
func GetAuthUser(r *http.Request) AuthUser {
	user, ok := r.Context().Value(AuthUserKey).(AuthUser)
	if !ok {
		return AuthUser{}
	}
	return user
}
