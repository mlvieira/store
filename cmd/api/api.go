package main

import (
	"log"

	"github.com/mlvieira/store/internal/application"
	"github.com/mlvieira/store/internal/handlers"
	"github.com/mlvieira/store/internal/router"
)

const version = "1.0.0"

func main() {
	baseApp, cleanup, err := application.NewBaseApplication(version)
	if err != nil {
		log.Fatalf("Error initializing application: %v", err)
	}
	defer cleanup()

	baseHandlers := handlers.NewHandlers(baseApp)

	apiRouter, err := router.InitRouter(baseHandlers, "api")
	if err != nil {
		log.Fatalf("Error initializing API router: %v", err)
	}

	if err := router.Serve(baseApp, apiRouter); err != nil {
		baseApp.ErrorLog.Fatalf("API server error: %v", err)
	}
}
