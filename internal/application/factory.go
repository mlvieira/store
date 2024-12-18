package application

import (
	"encoding/gob"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/mlvieira/store/internal/config"
	"github.com/mlvieira/store/internal/driver"
	"github.com/mlvieira/store/internal/models"
	"github.com/mlvieira/store/internal/render"
	"github.com/mlvieira/store/internal/repository"
	"github.com/mlvieira/store/internal/services"
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
	services := services.NewServices(repositories)
	renderer := render.NewRenderer(cfg.Env, cfg.Stripe.Key, cfg.API, errorLog)

	baseApp := &Application{
		Config:       cfg,
		InfoLog:      infoLog,
		ErrorLog:     errorLog,
		Version:      version,
		Repositories: repositories,
		Renderer:     renderer,
		Session:      sessionManager,
		Services:     services,
	}

	gob.Register(models.TransactionData{})

	return baseApp, cleanup, nil
}
