package web

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mlvieira/store/internal/cards"
	"github.com/mlvieira/store/internal/handlers"
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

	cardHolder := r.Form.Get("cardholder_name")
	email := r.Form.Get("email")
	paymentIntent := r.Form.Get("payment_intent")
	paymentMethod := r.Form.Get("payment_method")
	paymentAmount := r.Form.Get("payment_amount")
	paymentCurrency := r.Form.Get("payment_currency")

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
	expiryMonth := strconv.FormatInt(pm.Card.ExpMonth, 10)
	expiryYear := strconv.FormatInt(pm.Card.ExpYear, 10)

	data := make(map[string]any)
	data["cardholder"] = cardHolder
	data["email"] = email
	data["pi"] = paymentIntent
	data["pm"] = paymentMethod
	data["pa"] = paymentAmount
	data["pc"] = paymentCurrency
	data["last_four"] = lastFour
	data["expiry_month"] = expiryMonth
	data["expiry_year"] = expiryYear
	data["bank_return_code"] = ci

	if err = h.App.Renderer.RenderTemplate(w, r, "succeeded", &render.TemplateData{
		Data: data,
	}); err != nil {
		h.App.ErrorLog.Println(err)
	}

}

// ChargeOnce renders the widget purchase page for a single charge.
func (h *WebHandlers) ChargeOnce(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	widgetID, _ := strconv.Atoi(id)

	widget, err := h.App.Repositories.Widget.GetWidgetByID(r.Context(), widgetID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	data := make(map[string]any)
	data["widget"] = widget

	if err := h.App.Renderer.RenderTemplate(w, r, "buy-once", &render.TemplateData{
		Data: data,
	}); err != nil {
		h.App.ErrorLog.Println(err)
	}
}
