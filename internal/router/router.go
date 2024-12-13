package router

import (
	"fmt"
	"net/http"

	"github.com/mlvieira/store/internal/handlers"
)

// InitRouter initializes the appropriate router based on the mode (API or Web).
func InitRouter(baseHandlers *handlers.Handlers, mode string) (http.Handler, error) {
	switch mode {
	case "api":
		return InitAPIRoutes(baseHandlers), nil
	case "web":
		return InitWebRoutes(baseHandlers), nil
	default:
		return nil, fmt.Errorf("invalid router mode: %s", mode)
	}
}
