package router

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
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
	mux.Get("/terminal", webHandlers.VirtualTerminal)
	mux.Post("/payment", webHandlers.PaymentSucceeded)
	mux.Get("/widget/{id}", webHandlers.ChargeOnce)

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	return mux
}
