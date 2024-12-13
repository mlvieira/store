package main

import (
	"github.com/mlvieira/store/internal/api"
	"github.com/mlvieira/store/internal/config"
	"github.com/mlvieira/store/internal/driver"
	"github.com/mlvieira/store/internal/repository"
)

const version = "1.0.0"

func main() {
	cfg := config.NewConfig()

	infoLog, errorLog := config.NewLoggers()

	conn, err := driver.OpenDB(cfg.db.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer conn.Close()

	repositories := repository.Repositories{
		Widget:      repository.NewWidgetRepository(conn),
		Transaction: repository.NewTransactionRepository(conn),
	}

	app := &api.Application{
		Config:       cfg,
		InfoLog:      infoLog,
		ErrorLog:     errorLog,
		Repositories: repositories,
	}

	if err := app.Serve(); err != nil {
		errorLog.Fatalf("Server error: %v", err)
	}
}
