package shared

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func Serve(port int, env string, handler http.Handler, infoLog *log.Logger) error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           handler,
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	infoLog.Printf("Starting HTTP server in %s mode on port %d", env, port)
	return srv.ListenAndServe()
}
