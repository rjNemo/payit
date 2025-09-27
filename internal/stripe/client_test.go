package stripe

import (
	"context"
	"errors"
	"testing"

	"github.com/rjNemo/payit/config"
	stripeapi "github.com/stripe/stripe-go/v78"
)

type fakeSessionCreator struct {
	lastParams *stripeapi.CheckoutSessionParams
	result     *stripeapi.CheckoutSession
	err        error
}

func (f *fakeSessionCreator) New(params *stripeapi.CheckoutSessionParams) (*stripeapi.CheckoutSession, error) {
	f.lastParams = params
	return f.result, f.err
}

func TestService_CreateSessionSuccess(t *testing.T) {
	product := testProductConfig()
	fake := &fakeSessionCreator{
		result: &stripeapi.CheckoutSession{
			ID:  "cs_test_123",
			URL: "https://stripe.test/checkout",
		},
	}

	svc := &Service{product: product, sessions: fake}

	res, err := svc.CreateSession(context.Background(), CheckoutSessionRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.ID != "cs_test_123" || res.URL != "https://stripe.test/checkout" {
		t.Fatalf("unexpected result: %#v", res)
	}

	params := fake.lastParams
	if params == nil {
		t.Fatal("expected params to be captured")
	}
	if params.Context == nil {
		t.Fatal("expected context to propagate")
	}
	if params.Mode == nil || *params.Mode != string(stripeapi.CheckoutSessionModePayment) {
		t.Fatalf("unexpected mode: %v", params.Mode)
	}
	if len(params.PaymentMethodTypes) != 1 || params.PaymentMethodTypes[0] == nil || *params.PaymentMethodTypes[0] != "card" {
		t.Fatalf("unexpected payment methods: %#v", params.PaymentMethodTypes)
	}

	if len(params.LineItems) != 1 {
		t.Fatalf("expected one line item, got %d", len(params.LineItems))
	}

	item := params.LineItems[0]
	if item.Quantity == nil || *item.Quantity != 1 {
		t.Fatalf("unexpected quantity: %v", item.Quantity)
	}
	if item.PriceData == nil {
		t.Fatal("expected price data to be set")
	}
	if item.PriceData.UnitAmount == nil || *item.PriceData.UnitAmount != product.PriceCents {
		t.Fatalf("unexpected unit amount: %v", item.PriceData.UnitAmount)
	}
	if item.PriceData.Currency == nil || *item.PriceData.Currency != product.Currency {
		t.Fatalf("unexpected currency: %v", item.PriceData.Currency)
	}
	if item.PriceData.ProductData == nil {
		t.Fatal("expected product data")
	}
	if item.PriceData.ProductData.Name == nil || *item.PriceData.ProductData.Name != product.Name {
		t.Fatalf("unexpected product name: %v", item.PriceData.ProductData.Name)
	}
	if item.PriceData.ProductData.Description == nil || *item.PriceData.ProductData.Description != product.Description {
		t.Fatalf("unexpected product description: %v", item.PriceData.ProductData.Description)
	}
}

func TestService_CreateSessionWithCustomQuantity(t *testing.T) {
	product := testProductConfig()
	fake := &fakeSessionCreator{
		result: &stripeapi.CheckoutSession{},
	}

	svc := &Service{product: product, sessions: fake}

	_, err := svc.CreateSession(context.Background(), CheckoutSessionRequest{Quantity: 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	params := fake.lastParams
	if params == nil || len(params.LineItems) != 1 {
		t.Fatalf("expected line item to be set: %#v", params)
	}

	qty := params.LineItems[0].Quantity
	if qty == nil || *qty != 3 {
		t.Fatalf("unexpected quantity: %v", qty)
	}
}

func TestService_CreateSessionError(t *testing.T) {
	product := testProductConfig()
	fake := &fakeSessionCreator{err: errors.New("boom")}

	svc := &Service{product: product, sessions: fake}

	_, err := svc.CreateSession(context.Background(), CheckoutSessionRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestService_CreateSessionNilSession(t *testing.T) {
	product := testProductConfig()
	fake := &fakeSessionCreator{}

	svc := &Service{product: product, sessions: fake}

	_, err := svc.CreateSession(context.Background(), CheckoutSessionRequest{})
	if err == nil {
		t.Fatal("expected error for nil session")
	}
}

func testProductConfig() config.ProductConfig {
	return config.ProductConfig{
		Name:        "Demo Widget",
		Description: "A very cool widget",
		PriceCents:  1999,
		Currency:    "usd",
		SuccessURL:  "https://example.com/success",
		CancelURL:   "https://example.com/cancel",
	}
}
