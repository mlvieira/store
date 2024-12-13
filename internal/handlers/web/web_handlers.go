package web

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mlvieira/store/internal/application"
	"github.com/mlvieira/store/internal/render"
)

// Handlers provides methods to handle web and API requests.
type Handlers struct {
	App *application.Application
}

// VirtualTerminal renders the virtual terminal page.
func (h *Handlers) VirtualTerminal(w http.ResponseWriter, r *http.Request) {
	if err := h.App.Renderer.RenderTemplate(w, r, "terminal", nil); err != nil {
		h.App.ErrorLog.Println(err)
	}
}

// PaymentSucceeded processes payment success and renders a success page.
func (h *Handlers) PaymentSucceeded(w http.ResponseWriter, r *http.Request) {
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

	data := make(map[string]any)
	data["cardholder"] = cardHolder
	data["email"] = email
	data["pi"] = paymentIntent
	data["pm"] = paymentMethod
	data["pa"] = paymentAmount
	data["pc"] = paymentCurrency

	if err = h.App.Renderer.RenderTemplate(w, r, "succeeded", &render.TemplateData{
		Data: data,
	}); err != nil {
		h.App.ErrorLog.Println(err)
	}

}

// ChargeOnce renders the widget purchase page for a single charge.
func (h *Handlers) ChargeOnce(w http.ResponseWriter, r *http.Request) {
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
