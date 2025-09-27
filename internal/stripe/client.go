package stripe

import (
	"context"
	"errors"

	"github.com/rjNemo/payit/config"
	stripe "github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/client"
)

type sessionCreator interface {
	New(params *stripe.CheckoutSessionParams) (*stripe.CheckoutSession, error)
}

// Service coordinates Stripe Checkout session creation for the demo product.
type Service struct {
	product  config.ProductConfig
	sessions sessionCreator
}

// NewService instantiates a Stripe-backed checkout service using the provided API key.
func NewService(apiKey string, product config.ProductConfig) *Service {
	stripeClient := client.New(apiKey, nil)

	return &Service{
		product:  product,
		sessions: stripeClient.CheckoutSessions,
	}
}

// CreateSession creates a Stripe Checkout session for the configured demo product.
func (s *Service) CreateSession(ctx context.Context, req CheckoutSessionRequest) (CheckoutSessionResult, error) {
	quantity := req.Quantity
	if quantity <= 0 {
		quantity = 1
	}

    params := &stripe.CheckoutSessionParams{}
    params.Context = ctx
	params.SuccessURL = stripe.String(s.product.SuccessURL)
	params.CancelURL = stripe.String(s.product.CancelURL)
	params.Mode = stripe.String(string(stripe.CheckoutSessionModePayment))
	params.PaymentMethodTypes = stripe.StringSlice([]string{"card"})

	params.LineItems = append(params.LineItems, &stripe.CheckoutSessionLineItemParams{
		Quantity: stripe.Int64(quantity),
		PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
			Currency:   stripe.String(s.product.Currency),
			UnitAmount: stripe.Int64(s.product.PriceCents),
			ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
				Name:        stripe.String(s.product.Name),
				Description: stripe.String(s.product.Description),
			},
		},
	})

	session, err := s.sessions.New(params)
	if err != nil {
		return CheckoutSessionResult{}, err
	}
	if session == nil {
		return CheckoutSessionResult{}, errors.New("stripe returned nil session")
	}

	return CheckoutSessionResult{ID: session.ID, URL: session.URL}, nil
}
