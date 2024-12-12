package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mlvieira/store/internal/cards"
)

type stripePayload struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
	Content string `json:"content,omitempty"`
	ID      int    `json:"id,omitempty"`
}

func (app *application) GetPaymentIntent(w http.ResponseWriter, r *http.Request) {
	var payload stripePayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		app.errorLog.Println(err)
		return
	}

	app.infoLog.Printf("Raw amount: %s", payload.Amount)

	amount, err := strconv.ParseFloat(payload.Amount, 64)
	if err != nil {
		app.errorLog.Println("Failed to convert to float:", err)
		return
	}

	app.infoLog.Printf("Calling card.Charge with Currency: %s, Amount: %f", payload.Currency, amount)
	card := cards.Card{
		Secret:   app.config.stripe.secret,
		Key:      app.config.stripe.key,
		Currency: payload.Currency,
	}

	pi, msg, err := card.Charge(payload.Currency, amount)
	if err != nil {
		app.errorLog.Printf("card.Charge failed: %v", err)

		j := jsonResponse{
			OK:      false,
			Message: msg,
		}

		if err := writeJson(w, http.StatusInternalServerError, j); err != nil {
			app.errorLog.Println("Failed to write JSON response:", err)
		}

		return
	}

	if err := writeJson(w, http.StatusOK, pi); err != nil {
		app.errorLog.Println("Failed to write JSON response:", err)
		return
	}
}

func (app *application) GetWidgetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	widgetID, _ := strconv.Atoi(id)

	widget, err := app.repositories.Widget.GetWidgetByID(r.Context(), widgetID)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	if err := writeJson(w, http.StatusOK, widget); err != nil {
		app.errorLog.Println("Failed to write JSON response:", err)
		return
	}
}

func writeJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
