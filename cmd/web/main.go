package main

import (
	"log"

	"github.com/mlvieira/store/internal/application"
	"github.com/mlvieira/store/internal/web"
)

const version = "1.0.0"
const cssVersion = "1"

func main() {
	baseApp, cleanup, err := application.NewBaseApplication(version)
	if err != nil {
		log.Fatalf("Error initializing application: %v", err)
	}
	defer cleanup()

	if err := web.Serve(baseApp); err != nil {
		baseApp.ErrorLog.Fatalf("Server error: %v", err)
	}
}
