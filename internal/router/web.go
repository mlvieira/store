package router

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/mlvieira/store/internal/handlers"
	"github.com/mlvieira/store/internal/handlers/web"
	"github.com/mlvieira/store/internal/middleware"
)

// InitWebRoutes sets up the routes and handlers for the web application.
func InitWebRoutes(baseHandlers *handlers.Handlers, scs *scs.SessionManager) http.Handler {
	mux := InitBaseRouter(false)

	mux.Use(middleware.MiddlewareSession(scs))

	webHandlers := web.NewWebHandlers(baseHandlers)

	mux.Get("/", webHandlers.Homepage)
	mux.Get("/widget/{id}", webHandlers.ChargeOnce)

	mux.Route("/payment", func(r chi.Router) {
		r.Post("/", webHandlers.PaymentSucceeded)
		r.Get("/receipt", webHandlers.Receipt)
	})

	mux.Route("/terminal", func(r chi.Router) {
		r.Get("/", webHandlers.VirtualTerminal)
		r.Post("/payment", webHandlers.PaymentVirtualTerminal)
		r.Get("/receipt", webHandlers.ReceiptVirtualTerminal)
	})

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	return mux
}
