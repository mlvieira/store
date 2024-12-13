package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mlvieira/store/internal/cards"
	"github.com/mlvieira/store/internal/handlers"
)

// APIHandlers embeds the shared Handlers to provide API-specific handlers.
type APIHandlers struct {
	*handlers.Handlers
}

// NewAPIHandlers initializes and returns an APIHandlers instance.
func NewAPIHandlers(h *handlers.Handlers) *APIHandlers {
	return &APIHandlers{Handlers: h}
}

// GetPaymentIntent creates a Stripe payment intent and returns it as JSON.
func (h *APIHandlers) GetPaymentIntent(w http.ResponseWriter, r *http.Request) {
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

// GetWidgetByID fetches a widget by its ID and returns it as JSON.
func (h *APIHandlers) GetWidgetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	widgetID, _ := strconv.Atoi(id)

	widget, err := h.App.Repositories.Widget.GetWidgetByID(r.Context(), widgetID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	writeJSON(w, http.StatusOK, widget, h.App.ErrorLog)
}
