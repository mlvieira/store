package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/terminal", app.VirtualTerminal)
	mux.Post("/payment", app.PaymentSucceeded)
	mux.Get("/widget/{id}", app.ChargeOnce)

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	return mux
}
