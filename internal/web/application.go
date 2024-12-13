package web

import (
	"html/template"
	"log"

	"github.com/mlvieira/store/internal/config"
	"github.com/mlvieira/store/internal/repository"
)

type Application struct {
	Config        *config.Config
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	TemplateCache map[string]*template.Template
	Version       string
	Repositories  repository.Repositories
}
