package application

import (
	"time"

	"github.com/alexedwards/scs/v2"
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

	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.Secure = cfg.Env == "production"

	repositories := repository.NewRepositories(conn)
	renderer := render.NewRenderer(cfg.Env, cfg.Stripe.Key, cfg.API, errorLog)

	baseApp := &Application{
		Config:       cfg,
		InfoLog:      infoLog,
		ErrorLog:     errorLog,
		Version:      version,
		Repositories: repositories,
		Renderer:     renderer,
		Session:      sessionManager,
	}

	return baseApp, cleanup, nil
}
