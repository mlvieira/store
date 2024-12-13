package application

import (
	"github.com/mlvieira/store/internal/config"
	"github.com/mlvieira/store/internal/driver"
	"github.com/mlvieira/store/internal/render"
	"github.com/mlvieira/store/internal/repository"
)

// NewBaseApplication initializes the application with configuration, logging, and resources.
func NewBaseApplication(version string) (*Application, func(), error) {
	cfg := config.NewConfig()

	infoLog, errorLog := config.NewLoggers()

	conn, err := driver.OpenDB(cfg.DB.DSN)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		conn.Close()
	}

	repositories := repository.NewRepositories(conn)
	renderer := render.NewRenderer(cfg.Env, cfg.Stripe.Key, cfg.API, errorLog)

	baseApp := &Application{
		Config:       cfg,
		InfoLog:      infoLog,
		ErrorLog:     errorLog,
		Version:      version,
		Repositories: repositories,
		Renderer:     renderer,
	}

	return baseApp, cleanup, nil
}
