package stripe

import (
	"context"
	"errors"

	stripesdk "github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/client"

	"github.com/rjNemo/payit/config"
	"github.com/rjNemo/payit/internal/payments"
)

type sessionCreator interface {
	New(params *stripesdk.CheckoutSessionParams) (*stripesdk.CheckoutSession, error)
}

// Driver implements the CheckoutDriver interface using the Stripe SDK.
type Driver struct {
	product  config.ProductConfig
	sessions sessionCreator
}

// NewDriver creates a Stripe-backed checkout driver with the provided credentials.
func NewDriver(apiKey string, product config.ProductConfig) *Driver {
	stripeClient := client.New(apiKey, nil)

	return &Driver{
		product:  product,
		sessions: stripeClient.CheckoutSessions,
	}
}

// CreateSession delegates session creation to Stripe, translating domain values to SDK params.
func (d *Driver) CreateSession(ctx context.Context, req payments.CheckoutSessionRequest) (payments.CheckoutSessionResult, error) {
	quantity := req.Quantity
	if quantity <= 0 {
		quantity = 1
	}

	params := &stripesdk.CheckoutSessionParams{}
	params.Context = ctx
	params.SuccessURL = stripesdk.String(d.product.SuccessURL)
	params.CancelURL = stripesdk.String(d.product.CancelURL)
	params.Mode = stripesdk.String(string(stripesdk.CheckoutSessionModePayment))
	params.PaymentMethodTypes = stripesdk.StringSlice([]string{"card"})

	params.LineItems = append(params.LineItems, &stripesdk.CheckoutSessionLineItemParams{
		Quantity: stripesdk.Int64(quantity),
		PriceData: &stripesdk.CheckoutSessionLineItemPriceDataParams{
			Currency:   stripesdk.String(d.product.Currency),
			UnitAmount: stripesdk.Int64(d.product.PriceCents),
			ProductData: &stripesdk.CheckoutSessionLineItemPriceDataProductDataParams{
				Name:        stripesdk.String(d.product.Name),
				Description: stripesdk.String(d.product.Description),
			},
		},
	})

	session, err := d.sessions.New(params)
	if err != nil {
		return payments.CheckoutSessionResult{}, err
	}
	if session == nil {
		return payments.CheckoutSessionResult{}, errors.New("stripe returned nil session")
	}

	return payments.CheckoutSessionResult{ID: session.ID, URL: session.URL}, nil
}
