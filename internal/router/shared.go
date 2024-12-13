package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/mlvieira/store/internal/application"
	"github.com/mlvieira/store/internal/shared"
)

// InitBaseRouter initializes a base router with common middleware.
func InitBaseRouter(enableCORS bool) *chi.Mux {
	mux := chi.NewRouter()

	if enableCORS {
		mux.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			AllowCredentials: false,
			MaxAge:           300,
		}))
	}

	return mux
}

// Serve initializes and starts the HTTP server using the shared Serve logic.
func Serve(app *application.Application, router http.Handler) error {
	return shared.Serve(
		app.Config.Port,
		app.Config.Env,
		router,
		app.InfoLog,
	)
}
