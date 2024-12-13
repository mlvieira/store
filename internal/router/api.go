package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mlvieira/store/internal/handlers"
	"github.com/mlvieira/store/internal/handlers/api"
)

// InitAPIRoutes sets up the routes and handlers for the API.
func InitAPIRoutes(baseHandlers *handlers.Handlers) http.Handler {
	mux := InitBaseRouter(true)

	apiHandlers := api.NewAPIHandlers(baseHandlers)

	mux.Route("/api", func(r chi.Router) {
		r.Post("/payment-intent", apiHandlers.GetPaymentIntent)
		r.Get("/widget/{id}", apiHandlers.GetWidgetByID)
	})

	return mux
}
