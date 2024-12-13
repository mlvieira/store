package main

import (
	"github.com/mlvieira/store/internal/api"
	"github.com/mlvieira/store/internal/application"
	"github.com/mlvieira/store/internal/config"
	"github.com/mlvieira/store/internal/driver"
	"github.com/mlvieira/store/internal/repository"
)

const version = "1.0.0"

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
		Repositories: repository.NewRepositories(conn),
	}

	apiApp := &api.Application{Application: baseApp}

	if err := apiApp.Serve(); err != nil {
		errorLog.Fatalf("Server error: %v", err)
	}
}
