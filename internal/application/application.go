package application

import (
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/mlvieira/store/internal/config"
	"github.com/mlvieira/store/internal/render"
	"github.com/mlvieira/store/internal/repository"
	"github.com/mlvieira/store/internal/services"
)

// Application holds the core application context and dependencies.
type Application struct {
	Config       *config.Config
	InfoLog      *log.Logger
	ErrorLog     *log.Logger
	Version      string
	Repositories *repository.Repositories
	Renderer     *render.Renderer
	Session      *scs.SessionManager
	Services     *services.Services
}
