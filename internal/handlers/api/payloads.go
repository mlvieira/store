package api

// stripePayload represents a payment intent request payload.
type stripePayload struct {
	Currency      string `json:"currency"`
	Amount        int64  `json:"amount"`
	PaymentMethod string `json:"payment_method"`
	Email         string `json:"email"`
	LastFour      string `json:"last_four"`
	PlanID        string `json:"plan_id"`
}
