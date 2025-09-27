package web

import (
	"context"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/rjNemo/payit/config"
	"github.com/rjNemo/payit/internal/payments"
	"github.com/rjNemo/payit/internal/payments/driver/stripe"
	"github.com/rjNemo/payit/internal/payments/service"
	webassets "github.com/rjNemo/payit/web"
)

type checkoutService interface {
	CreateSession(context.Context, payments.CheckoutSessionRequest) (payments.CheckoutSessionResult, error)
}

// Handler aggregates dependencies required by HTTP handlers.
type Handler struct {
	cfg      config.Config
	checkout checkoutService
	page     *template.Template
	fs       fs.FS
}

// NewServer constructs the root HTTP handler, wiring Stripe-backed endpoints as they are implemented.
func NewServer(cfg config.Config) http.Handler {
	driver := stripe.NewDriver(cfg.StripeSecretKey, cfg.Product)
	checkoutSvc := service.NewCheckoutService(driver)
	tmpl := template.Must(template.ParseFS(webassets.Assets, "templates/index.html"))
	staticFS, err := fs.Sub(webassets.Assets, "static")
	if err != nil {
		panic(fmt.Errorf("failed to load static assets: %w", err))
	}

	h := &Handler{cfg: cfg, checkout: checkoutSvc, page: tmpl, fs: staticFS}

	mux := http.NewServeMux()
	h.registerRoutes(mux)

	return mux
}
