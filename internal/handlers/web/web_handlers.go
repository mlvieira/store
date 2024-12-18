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

// PaymentVirtualTerminal processes payment success from virtual terminal and renders a success page.
func (h *WebHandlers) PaymentVirtualTerminal(w http.ResponseWriter, r *http.Request) {
	txnData, err := h.GetTransactionData(r)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	intExpiryMonth, _ := strconv.Atoi(txnData.ExpiryMonth)
	intExpiryYear, _ := strconv.Atoi(txnData.ExpiryYear)

	txn := models.Transaction{
		Amount:              txnData.PaymentAmount,
		Currency:            txnData.PaymentCurrency,
		LastFour:            txnData.LastFour,
		ExpiryMonth:         intExpiryMonth,
		ExpiryYear:          intExpiryYear,
		BankReturnCode:      txnData.BankReturnCode,
		TransactionStatusID: 2,
		PaymentIntent:       txnData.PaymentIntentID,
		PaymentMethod:       txnData.PaymentMethodID,
	}

	_, err = h.App.Services.TransactionService.SaveTransaction(r.Context(), txn)
	if err != nil {
		h.App.InfoLog.Println(err)
		return
	}

	h.App.Session.Put(r.Context(), "receipt", txnData)

	http.Redirect(w, r, "/terminal/receipt", http.StatusSeeOther)
}

// GetTransactionData gets transaction data from post request and stripe
func (h *WebHandlers) GetTransactionData(r *http.Request) (models.TransactionData, error) {
	var txnData models.TransactionData

	err := r.ParseForm()
	if err != nil {
		h.App.ErrorLog.Println(err)
		return txnData, err
	}

	firstName := r.Form.Get("first_name")
	lastName := r.Form.Get("last_name")
	email := r.Form.Get("email")
	paymentIntent := r.Form.Get("payment_intent")
	paymentMethod := r.Form.Get("payment_method")
	paymentAmount := r.Form.Get("payment_amount")
	paymentCurrency := r.Form.Get("payment_currency")
	amount, err := strconv.ParseInt(paymentAmount, 10, 64)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return txnData, err
	}

	card := cards.Card{
		Secret: h.App.Config.Stripe.Secret,
		Key:    h.App.Config.Stripe.Key,
	}

	ci, err := card.RetrieveChargeID(paymentIntent)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return txnData, err
	}

	pm, err := card.GetPaymentMethod(paymentMethod)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return txnData, err
	}

	lastFour := pm.Card.Last4
	expiryMonth := pm.Card.ExpMonth
	expiryYear := pm.Card.ExpYear

	txnData = models.TransactionData{
		FirstName:       firstName,
		LastName:        lastName,
		Email:           email,
		PaymentIntentID: paymentIntent,
		PaymentMethodID: paymentMethod,
		PaymentAmount:   amount,
		PaymentCurrency: paymentCurrency,
		LastFour:        lastFour,
		ExpiryMonth:     strconv.FormatInt(expiryMonth, 10),
		ExpiryYear:      strconv.FormatInt(expiryYear, 10),
		BankReturnCode:  ci,
	}

	return txnData, nil

}

// PaymentSucceeded processes payment success and renders a success page.
func (h *WebHandlers) PaymentSucceeded(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	widgetID, err := strconv.Atoi(r.Form.Get("widget_id"))
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	txnData, err := h.GetTransactionData(r)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	customer := models.Customer{
		FirstName: txnData.FirstName,
		LastName:  txnData.LastName,
		Email:     txnData.Email,
	}

	customerID, err := h.App.Services.CustomerService.SaveCustomer(r.Context(), customer)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	intExpiryMonth, _ := strconv.Atoi(txnData.ExpiryMonth)
	intExpiryYear, _ := strconv.Atoi(txnData.ExpiryYear)

	txn := models.Transaction{
		Amount:              txnData.PaymentAmount,
		Currency:            txnData.PaymentCurrency,
		LastFour:            txnData.LastFour,
		ExpiryMonth:         intExpiryMonth,
		ExpiryYear:          intExpiryYear,
		BankReturnCode:      txnData.BankReturnCode,
		TransactionStatusID: 2,
		PaymentIntent:       txnData.PaymentIntentID,
		PaymentMethod:       txnData.PaymentMethodID,
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
		Amount:        txnData.PaymentAmount,
	}

	_, err = h.App.Services.OrderService.PlaceOrder(r.Context(), order)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	h.App.Session.Put(r.Context(), "receipt", txnData)

	http.Redirect(w, r, "/receipt", http.StatusSeeOther)
}

// ReceiptVirtualTerminal display receipt page for orders from virtual terminal
func (h *WebHandlers) ReceiptVirtualTerminal(w http.ResponseWriter, r *http.Request) {
	txn := h.App.Session.Get(r.Context(), "receipt").(models.TransactionData)

	data := map[string]any{
		"txn": txn,
	}

	h.App.Session.Remove(r.Context(), "receipt")

	if err := h.App.Renderer.RenderTemplate(w, r, "terminal-receipt", &render.TemplateData{
		Data: data,
	}); err != nil {
		h.App.ErrorLog.Println(err)
	}
}

// Receipt display receipt page for orders from user order page
func (h *WebHandlers) Receipt(w http.ResponseWriter, r *http.Request) {
	txn := h.App.Session.Get(r.Context(), "receipt").(models.TransactionData)

	data := map[string]any{
		"txn": txn,
	}

	h.App.Session.Remove(r.Context(), "receipt")

	if err := h.App.Renderer.RenderTemplate(w, r, "receipt", &render.TemplateData{
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
