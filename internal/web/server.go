package web

import (
	"fmt"
	"net/http"
	"time"
)

func (app *Application) Serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.Config.Port),
		Handler:           WebRoutes(app),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	app.InfoLog.Printf("Starting HTTP server in %s mode on port %d", app.Config.Env, app.Config.Port)
	return srv.ListenAndServe()
}
