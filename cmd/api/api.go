package main

import (
	"log"

	"github.com/mlvieira/store/internal/api"
	"github.com/mlvieira/store/internal/application"
)

const version = "1.0.0"

func main() {
	baseApp, cleanup, err := application.NewBaseApplication(version)
	if err != nil {
		log.Fatalf("Error initializing application: %v", err)
	}
	defer cleanup()

	if err := api.Serve(baseApp); err != nil {
		baseApp.ErrorLog.Fatalf("API server error: %v", err)
	}
}
