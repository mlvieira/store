package cards

import (
	"errors"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/charge"
	"github.com/stripe/stripe-go/v81/customer"
	"github.com/stripe/stripe-go/v81/paymentintent"
	"github.com/stripe/stripe-go/v81/paymentmethod"
	"github.com/stripe/stripe-go/v81/setupintent"
	"github.com/stripe/stripe-go/v81/subscription"
)

// Card represents payment card details for transactions.
type Card struct {
	Secret   string
	Key      string
	Currency string
}

// Transaction represents a financial transaction record.
type Transaction struct {
	TransactionStatusID int
	Amount              int
	Currency            string
	LastFour            string
	BankReturnCode      string
}

// Charge creates a payment intent for a specified currency and amount.
func (c *Card) Charge(currency string, amount int64) (*stripe.PaymentIntent, string, error) {
	return c.CreatePaymentIntent(currency, amount)
}

// CreatePaymentIntent generates a Stripe payment intent for a given currency and amount.
func (c *Card) CreatePaymentIntent(currency string, amount int64) (*stripe.PaymentIntent, string, error) {
	stripe.Key = c.Secret

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(currency),
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		msg := ""
		if stripeErr, ok := err.(*stripe.Error); ok {
			msg = cardErrorMessage(stripeErr.Code)
		}

		return nil, msg, err
	}

	return pi, "", nil
}

// CreateSetupIntent generates a Stripe SetupIntent to save payment details for future use.
func (c *Card) CreateSetupIntent(customerID string, paymentMethodID string) (*stripe.SetupIntent, string, error) {
	stripe.Key = c.Secret

	params := &stripe.SetupIntentParams{
		Customer: stripe.String(customerID),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		Usage: stripe.String(string(stripe.SetupIntentUsageOffSession)),
	}

	if paymentMethodID != "" {
		params.PaymentMethod = stripe.String(paymentMethodID)
	}

	si, err := setupintent.New(params)
	if err != nil {
		msg := ""
		if stripeErr, ok := err.(*stripe.Error); ok {
			msg = cardErrorMessage(stripeErr.Code)
		}
		return nil, msg, err
	}

	return si, "", nil
}

// GetPaymentMethod gets the payment method by payment intent id
func (c *Card) GetPaymentMethod(s string) (*stripe.PaymentMethod, error) {
	stripe.Key = c.Secret

	pm, err := paymentmethod.Get(s, nil)
	if err != nil {
		return nil, err
	}

	return pm, nil
}

// RetrievePaymentIntent gets an existing payment intent by id
func (c *Card) RetrievePaymentIntent(id string) (*stripe.PaymentIntent, error) {
	stripe.Key = c.Secret

	pi, err := paymentintent.Get(id, nil)
	if err != nil {
		return nil, err
	}

	return pi, nil
}

// RetrieveChargeID retrieves the charge ID associated with a PaymentIntent
func (c *Card) RetrieveChargeID(paymentIntentID string) (string, error) {
	stripe.Key = c.Secret

	params := &stripe.ChargeListParams{
		PaymentIntent: stripe.String(paymentIntentID),
	}
	params.Filters.AddFilter("limit", "", "1")

	iter := charge.List(params)
	if iter.Next() {
		ch := iter.Charge()
		return ch.ID, nil
	}

	if iter.Err() != nil {
		return "", iter.Err()
	}

	return "", errors.New("no charges found for this PaymentIntent")
}

// CreateCustomer creates a customer in Stripe
func (c *Card) CreateCustomer(pm, email string) (*stripe.Customer, string, error) {
	stripe.Key = c.Secret

	params := &stripe.CustomerParams{
		PaymentMethod: stripe.String(pm),
		Email:         stripe.String(email),
		InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
			DefaultPaymentMethod: stripe.String(pm),
		},
	}

	cust, err := customer.New(params)
	if err != nil {
		msg := ""
		if stripeErr, ok := err.(*stripe.Error); ok {
			msg = cardErrorMessage(stripeErr.Code)
		}

		return nil, msg, err
	}

	return cust, "", nil
}

// SubscribeToPlan subscribes a customer to a Stripe plan
func (c *Card) SubscribeToPlan(cust *stripe.Customer, plan, email, last4, cardType string) (*stripe.Subscription, error) {
	stripeCustomerID := cust.ID
	items := []*stripe.SubscriptionItemsParams{
		{Plan: stripe.String(plan)},
	}

	params := &stripe.SubscriptionParams{
		Customer: stripe.String(stripeCustomerID),
		Items:    items,
	}

	params.AddMetadata("last_four", last4)
	params.AddMetadata("card_type", cardType)
	params.AddExpand("latest_invoice.payment_intent")

	sub, err := subscription.New(params)
	if err != nil {
		return nil, err
	}
	return sub, nil
}

// cardErrorMessage maps Stripe error codes to user-friendly error messages.
func cardErrorMessage(code stripe.ErrorCode) string {
	var msg = ""
	switch code {
	case stripe.ErrorCodeExpiredCard:
		msg = "Your card is expired"
	case stripe.ErrorCodeIncorrectCVC:
		msg = "Incorrect CVC code"
	case stripe.ErrorCodeIncorrectZip:
		msg = "Incorrect zip/postal code"
	case stripe.ErrorCodeAmountTooLarge:
		msg = "The amount is too much to charge to your card"
	case stripe.ErrorCodeAmountTooSmall:
		msg = "The amount is too small to charge to your card"
	case stripe.ErrorCodeBalanceInsufficient:
		msg = "Insufficient balance"
	case stripe.ErrorCodePostalCodeInvalid:
		msg = "Your postal code is invalid"
	case stripe.ErrorCodeCardDeclined:
	default:
		msg = "Your card was declined"
	}

	return msg
}
