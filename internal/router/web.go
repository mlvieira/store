package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mlvieira/store/internal/application"
	"github.com/mlvieira/store/internal/handlers/web"
)

// InitWebRoutes sets up the routes and handlers for the web application.
func InitWebRoutes(app *application.Application) http.Handler {
	mux := chi.NewRouter()

	handlers := &web.Handlers{App: app}

	mux.Get("/terminal", handlers.VirtualTerminal)
	mux.Post("/payment", handlers.PaymentSucceeded)
	mux.Get("/widget/{id}", handlers.ChargeOnce)

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	return mux
}
