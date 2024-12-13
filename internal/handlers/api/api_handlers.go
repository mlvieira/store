package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mlvieira/store/internal/application"
	"github.com/mlvieira/store/internal/cards"
)

type Handlers struct {
	App *application.Application
}

func (h *Handlers) GetPaymentIntent(w http.ResponseWriter, r *http.Request) {
	var payload stripePayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	h.App.InfoLog.Printf("Raw amount: %s", payload.Amount)

	amount, err := strconv.ParseFloat(payload.Amount, 64)
	if err != nil {
		h.App.ErrorLog.Println("Failed to convert to float:", err)
		return
	}

	h.App.InfoLog.Printf("Calling card.Charge with Currency: %s, Amount: %f", payload.Currency, amount)
	card := cards.Card{
		Secret:   h.App.Config.Stripe.Secret,
		Key:      h.App.Config.Stripe.Key,
		Currency: payload.Currency,
	}

	pi, msg, err := card.Charge(payload.Currency, amount)
	if err != nil {
		h.App.ErrorLog.Printf("card.Charge failed: %v", err)

		writeJSON(w, http.StatusInternalServerError, jsonResponse{
			OK:      false,
			Message: msg,
		}, h.App.ErrorLog)

		return
	}

	writeJSON(w, http.StatusOK, pi, h.App.ErrorLog)
}

func (h *Handlers) GetWidgetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	widgetID, _ := strconv.Atoi(id)

	widget, err := h.App.Repositories.Widget.GetWidgetByID(r.Context(), widgetID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	writeJSON(w, http.StatusOK, widget, h.App.ErrorLog)
}
