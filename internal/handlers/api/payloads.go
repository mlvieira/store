package api

// stripePayload represents a payment intent request payload.
type stripePayload struct {
	Currency string `json:"currency"`
	Amount   int64  `json:"amount"`
}
