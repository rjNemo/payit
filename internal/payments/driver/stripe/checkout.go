package stripe

import (
	"context"
	"errors"

	"github.com/stripe/stripe-go/v83"

	"github.com/rjNemo/payit/config"
	"github.com/rjNemo/payit/internal/payments"
)

type sessionCreator interface {
	Create(ctx context.Context, params *stripe.CheckoutSessionCreateParams) (*stripe.CheckoutSession, error)
}

// Driver implements the CheckoutDriver interface using the Stripe SDK.
type Driver struct {
	product  config.ProductConfig
	sessions sessionCreator
}

// NewDriver creates a Stripe-backed checkout driver with the provided credentials.
func NewDriver(apiKey string, product config.ProductConfig) *Driver {
	stripeClient := stripe.NewClient(apiKey, nil)

	return &Driver{
		product:  product,
		sessions: stripeClient.V1CheckoutSessions,
	}
}

// CreateSession delegates session creation to Stripe, translating domain values to SDK params.
func (d *Driver) CreateSession(ctx context.Context, req payments.CheckoutSessionRequest) (payments.CheckoutSessionResult, error) {
	quantity := req.Quantity
	if quantity <= 0 {
		quantity = 1
	}

	params := &stripe.CheckoutSessionCreateParams{}
	params.Context = ctx
	params.SuccessURL = stripe.String(d.product.SuccessURL)
	params.CancelURL = stripe.String(d.product.CancelURL)
	params.Mode = stripe.String(string(stripe.CheckoutSessionModePayment))
	params.PaymentMethodTypes = stripe.StringSlice([]string{"card"})

	params.LineItems = append(params.LineItems, &stripe.CheckoutSessionCreateLineItemParams{
		Quantity: stripe.Int64(quantity),
		PriceData: &stripe.CheckoutSessionCreateLineItemPriceDataParams{
			Currency:   stripe.String(d.product.Currency),
			UnitAmount: stripe.Int64(d.product.PriceCents),
			ProductData: &stripe.CheckoutSessionCreateLineItemPriceDataProductDataParams{
				Name:        stripe.String(d.product.Name),
				Description: stripe.String(d.product.Description),
			},
		},
	})

	session, err := d.sessions.Create(ctx, params)
	if err != nil {
		return payments.CheckoutSessionResult{}, err
	}
	if session == nil {
		return payments.CheckoutSessionResult{}, errors.New("stripe returned nil session")
	}

	return payments.CheckoutSessionResult{ID: session.ID, URL: session.URL}, nil
}
