package stripe

// CheckoutSessionRequest captures optional inputs for creating a checkout session.
type CheckoutSessionRequest struct {
	Quantity int64 `json:"quantity"`
}

// CheckoutSessionResult contains the data returned to callers initiating checkout.
type CheckoutSessionResult struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}
