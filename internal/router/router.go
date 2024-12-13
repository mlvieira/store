package router

import (
	"fmt"
	"net/http"

	"github.com/mlvieira/store/internal/application"
	"github.com/mlvieira/store/internal/shared"
)

// InitRouter initializes the appropriate router based on the mode (API or Web).
func InitRouter(app *application.Application, mode string) (http.Handler, error) {
	switch mode {
	case "api":
		return InitAPIRoutes(app), nil
	case "web":
		return InitWebRoutes(app), nil
	default:
		return nil, ErrInvalidMode
	}
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

var ErrInvalidMode = fmt.Errorf("invalid router mode")
