package main

import (
	"log"

	"github.com/mlvieira/store/internal/application"
	"github.com/mlvieira/store/internal/router"
)

const version = "1.0.0"

func main() {
	baseApp, cleanup, err := application.NewBaseApplication(version)
	if err != nil {
		log.Fatalf("Error initializing application: %v", err)
	}
	defer cleanup()

	webRouter, err := router.InitRouter(baseApp, "web")
	if err != nil {
		log.Fatalf("Error initializing web router: %v", err)
	}

	if err := router.Serve(baseApp, webRouter); err != nil {
		baseApp.ErrorLog.Fatalf("Server error: %v", err)
	}
}
