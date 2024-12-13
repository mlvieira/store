package web

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mlvieira/store/internal/render"
)

func (app *Application) VirtualTerminal(w http.ResponseWriter, r *http.Request) {
	if err := app.Renderer.RenderTemplate(w, r, "terminal", nil); err != nil {
		app.ErrorLog.Println(err)
	}
}

func (app *Application) PaymentSucceeded(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Println(err)
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

	if err = app.Renderer.RenderTemplate(w, r, "succeeded", &render.TemplateData{
		Data: data,
	}); err != nil {
		app.ErrorLog.Println(err)
	}

}

func (app *Application) ChargeOnce(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	widgetID, _ := strconv.Atoi(id)

	widget, err := app.Repositories.Widget.GetWidgetByID(r.Context(), widgetID)
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	data := make(map[string]any)
	data["widget"] = widget

	if err := app.Renderer.RenderTemplate(w, r, "buy-once", &render.TemplateData{
		Data: data,
	}); err != nil {
		app.ErrorLog.Println(err)
	}
}
