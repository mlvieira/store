package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mlvieira/store/internal/cards"
	"github.com/mlvieira/store/internal/handlers"
	"github.com/mlvieira/store/internal/models"
	"github.com/stripe/stripe-go/v81"
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

	h.App.InfoLog.Printf("Calling card.Charge with Currency: %s, Amount: %d", payload.Currency, payload.Amount)
	card := cards.Card{
		Secret:   h.App.Config.Stripe.Secret,
		Key:      h.App.Config.Stripe.Key,
		Currency: payload.Currency,
	}

	pi, msg, err := card.Charge(payload.Currency, payload.Amount)
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
	widgetID, err := strconv.Atoi(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	widget, err := h.App.Repositories.Widget.GetWidgetByID(r.Context(), widgetID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	writeJSON(w, http.StatusOK, widget, h.App.ErrorLog)
}

// CreateSubscription creates a subscription for a product
func (h *APIHandlers) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var payload stripePayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	h.App.InfoLog.Println(payload.Email, payload.LastFour, payload.PaymentMethod, payload.PlanID)

	card := cards.Card{
		Secret:   h.App.Config.Stripe.Secret,
		Key:      h.App.Config.Stripe.Key,
		Currency: payload.Currency,
	}

	var subscription *stripe.Subscription

	stripeCustomer, msg, err := card.CreateCustomer(payload.PaymentMethod, payload.Email)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, jsonResponse{
			OK:      false,
			Message: msg,
		}, h.App.ErrorLog)
		return
	}

	sp, msg, err := card.CreateSetupIntent(stripeCustomer.ID, payload.PaymentMethod)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, jsonResponse{
			OK:      false,
			Message: msg,
		}, h.App.ErrorLog)
		return
	}

	subscription, err = card.SubscribeToPlan(stripeCustomer, payload.PlanID, payload.Email, payload.LastFour, "")
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, jsonResponse{
			OK:      false,
			Message: "Error while subscribing to plan",
		}, h.App.ErrorLog)
		return
	}

	h.App.InfoLog.Println(subscription.ID)

	productID, err := strconv.Atoi(payload.ProductID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, jsonResponse{
			OK:      false,
			Message: "Error converting product ID",
		}, h.App.ErrorLog)
		return
	}

	cust := models.Customer{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
	}

	customerID, err := h.App.Services.CustomerService.SaveCustomer(r.Context(), cust)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, jsonResponse{
			OK:      false,
			Message: "Error saving customer",
		}, h.App.ErrorLog)
		return
	}

	widget, err := h.App.Repositories.Widget.GetWidgetByID(r.Context(), productID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, jsonResponse{
			OK:      false,
			Message: "Error getting product info",
		}, h.App.ErrorLog)
		return
	}

	txn := models.Transaction{
		Amount:              widget.Price,
		Currency:            payload.Currency,
		LastFour:            payload.LastFour,
		ExpiryMonth:         payload.ExpiryMonth,
		ExpiryYear:          payload.ExpiryYear,
		TransactionStatusID: 2,
	}

	txnID, err := h.App.Services.TransactionService.SaveTransaction(r.Context(), txn)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, jsonResponse{
			OK:      false,
			Message: "Error saving transaction",
		}, h.App.ErrorLog)
		return
	}

	order := models.Order{
		WidgetID:      widget.ID,
		TransactionID: txnID,
		CustomerID:    customerID,
		StatusID:      1,
		Quantity:      1,
		Amount:        widget.Price,
	}

	_, err = h.App.Services.OrderService.PlaceOrder(r.Context(), order)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, jsonResponse{
			OK:      false,
			Message: "Error saving order",
		}, h.App.ErrorLog)
		return
	}

	writeJSON(w, http.StatusOK, jsonResponse{
		OK:      true,
		Message: "Transaction successful",
		Content: sp.ClientSecret,
	}, h.App.InfoLog)
}
