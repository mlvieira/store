package main

import (
	"html/template"

	"github.com/mlvieira/store/internal/config"
	"github.com/mlvieira/store/internal/driver"
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

	repositories := repository.Repositories{
		Widget:      repository.NewWidgetRepository(conn),
		Transaction: repository.NewTransactionRepository(conn),
	}

	tc := make(map[string]*template.Template)

	app := &web.Application{
		Config:        cfg,
		InfoLog:       infoLog,
		ErrorLog:      errorLog,
		TemplateCache: tc,
		Version:       version,
		Repositories:  repositories,
	}

	if err := app.Serve(); err != nil {
		errorLog.Fatalf("Server error: %v", err)
	}
}
