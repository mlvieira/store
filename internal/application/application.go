package application

import (
	"log"

	"github.com/mlvieira/store/internal/config"
	"github.com/mlvieira/store/internal/render"
	"github.com/mlvieira/store/internal/repository"
)

type Application struct {
	Config       *config.Config
	InfoLog      *log.Logger
	ErrorLog     *log.Logger
	Version      string
	Repositories *repository.Repositories
	Renderer     *render.Renderer
}
