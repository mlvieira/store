package main

import (
	"github.com/mlvieira/store/internal/application"
	"github.com/mlvieira/store/internal/config"
	"github.com/mlvieira/store/internal/driver"
	"github.com/mlvieira/store/internal/render"
	"github.com/mlvieira/store/internal/repository"
	"github.com/mlvieira/store/internal/web"
)

const version = "1.0.0"
const cssVersion = "1"

func main() {
	cfg := config.NewConfig()

	infoLog, errorLog := config.NewLoggers()

	conn, err := driver.OpenDB(cfg.DB.DSN)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer conn.Close()

	baseApp := &application.Application{
		Config:       cfg,
		InfoLog:      infoLog,
		ErrorLog:     errorLog,
		Version:      version,
		Repositories: repository.NewRepositories(conn),
		Renderer:     render.NewRenderer(cfg.Env, cfg.Stripe.Key, cfg.API, errorLog),
	}

	webApp := &web.Application{Application: baseApp}

	if err := webApp.Serve(); err != nil {
		errorLog.Fatalf("Server error: %v", err)
	}
}
