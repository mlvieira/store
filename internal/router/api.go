package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/mlvieira/store/internal/application"
	"github.com/mlvieira/store/internal/handlers/api"
)

// InitAPIRoutes sets up the routes and handlers for the API.
func InitAPIRoutes(app *application.Application) http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	handlers := &api.Handlers{App: app}

	mux.Route("/api", func(r chi.Router) {
		r.Post("/payment-intent", handlers.GetPaymentIntent)
		r.Get("/widget/{id}", handlers.GetWidgetByID)
	})

	return mux
}
