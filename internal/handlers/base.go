package handlers

import "github.com/mlvieira/store/internal/application"

// Handlers provides methods to handle web and API requests.
type Handlers struct {
	App *application.Application
}

// NewHandlers initializes a new Handlers instance.
func NewHandlers(app *application.Application) *Handlers {
	return &Handlers{App: app}
}
