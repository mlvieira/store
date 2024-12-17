package web

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mlvieira/store/internal/cards"
	"github.com/mlvieira/store/internal/handlers"
	"github.com/mlvieira/store/internal/models"
	"github.com/mlvieira/store/internal/render"
)

// WebHandlers embeds the shared Handlers to provide Web-specific handlers.
type WebHandlers struct {
	*handlers.Handlers
}

// NewWebHandlers initializes and returns a WebHandlers instance.
func NewWebHandlers(h *handlers.Handlers) *WebHandlers {
	return &WebHandlers{Handlers: h}
}

// Homepage renders the home page.
func (h *WebHandlers) Homepage(w http.ResponseWriter, r *http.Request) {
	if err := h.App.Renderer.RenderTemplate(w, r, "home", nil); err != nil {
		h.App.ErrorLog.Println(err)
	}
}

// VirtualTerminal renders the virtual terminal page.
func (h *WebHandlers) VirtualTerminal(w http.ResponseWriter, r *http.Request) {
	if err := h.App.Renderer.RenderTemplate(w, r, "terminal", nil); err != nil {
		h.App.ErrorLog.Println(err)
	}
}

// PaymentSucceeded processes payment success and renders a success page.
func (h *WebHandlers) PaymentSucceeded(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	firstName := r.Form.Get("first_name")
	lastName := r.Form.Get("last_name")
	email := r.Form.Get("email")
	paymentIntent := r.Form.Get("payment_intent")
	paymentMethod := r.Form.Get("payment_method")
	paymentAmount := r.Form.Get("payment_amount")
	paymentCurrency := r.Form.Get("payment_currency")
	widgetID, err := strconv.Atoi(r.Form.Get("widget_id"))
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	card := cards.Card{
		Secret: h.App.Config.Stripe.Secret,
		Key:    h.App.Config.Stripe.Key,
	}

	ci, err := card.RetrieveChargeID(paymentIntent)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	pm, err := card.GetPaymentMethod(paymentMethod)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	lastFour := pm.Card.Last4
	expiryMonth := pm.Card.ExpMonth
	expiryYear := pm.Card.ExpYear

	customer := models.Customer{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	customerID, err := h.App.Services.CustomerService.SaveCustomer(r.Context(), customer)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	amount, err := strconv.ParseInt(paymentAmount, 10, 64)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	txn := models.Transaction{
		Amount:              amount,
		Currency:            paymentCurrency,
		LastFour:            lastFour,
		ExpiryMonth:         int(expiryMonth),
		ExpiryYear:          int(expiryYear),
		BankReturnCode:      ci,
		TransactionStatusID: 2,
	}

	txnID, err := h.App.Services.TransactionService.SaveTransaction(r.Context(), txn)
	if err != nil {
		h.App.InfoLog.Println(err)
		return
	}

	order := models.Order{
		WidgetID:      widgetID,
		TransactionID: txnID,
		CustomerID:    customerID,
		StatusID:      1,
		Quantity:      1,
		Amount:        amount,
	}

	_, err = h.App.Services.OrderService.PlaceOrder(r.Context(), order)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	data := map[string]any{
		"email":            email,
		"pi":               paymentIntent,
		"pm":               paymentMethod,
		"pa":               amount,
		"pc":               paymentCurrency,
		"last_four":        lastFour,
		"expiry_month":     strconv.FormatInt(pm.Card.ExpMonth, 10),
		"expiry_year":      strconv.FormatInt(pm.Card.ExpYear, 10),
		"bank_return_code": ci,
		"first_name":       firstName,
		"last_name":        lastName,
	}

	if err = h.App.Renderer.RenderTemplate(w, r, "succeeded", &render.TemplateData{
		Data: data,
	}); err != nil {
		h.App.ErrorLog.Println(err)
	}

}

// ChargeOnce renders the widget purchase page for a single charge.
func (h *WebHandlers) ChargeOnce(w http.ResponseWriter, r *http.Request) {
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

	data := map[string]any{
		"widget": widget,
	}

	if err := h.App.Renderer.RenderTemplate(w, r, "buy-once", &render.TemplateData{
		Data: data,
	}); err != nil {
		h.App.ErrorLog.Println(err)
	}
}
