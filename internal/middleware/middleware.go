package middleware

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
)

// MiddlewareSession wraps the session manager around requests.
func MiddlewareSession(sessionManager *scs.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return sessionManager.LoadAndSave(next)
	}
}
