package api

type stripePayload struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}
