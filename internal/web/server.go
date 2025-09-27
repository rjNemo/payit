package web

import (
	"context"
	"net/http"

	stripe "github.com/rjNemo/payit/internal/stripe"

	"github.com/rjNemo/payit/config"
)

type checkoutService interface {
	CreateSession(context.Context, stripe.CheckoutSessionRequest) (stripe.CheckoutSessionResult, error)
}

// Handler aggregates dependencies required by HTTP handlers.
type Handler struct {
	cfg      config.Config
	checkout checkoutService
}

// NewServer constructs the root HTTP handler, wiring Stripe-backed endpoints as they are implemented.
func NewServer(cfg config.Config) http.Handler {
	checkoutSvc := stripe.NewService(cfg.StripeSecretKey, cfg.Product)
	h := &Handler{cfg: cfg, checkout: checkoutSvc}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/checkout", h.createCheckoutSession)

	return mux
}
