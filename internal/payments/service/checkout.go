package service

import (
	"context"

	"github.com/rjNemo/payit/internal/payments"
)

// CheckoutDriver represents a payment provider capable of creating checkout sessions.
type CheckoutDriver interface {
	CreateSession(ctx context.Context, req payments.CheckoutSessionRequest) (payments.CheckoutSessionResult, error)
}

// CheckoutService contains provider-agnostic business rules for initiating checkout flows.
type CheckoutService struct {
	driver CheckoutDriver
}

// NewCheckoutService wires the given driver into a reusable checkout service.
func NewCheckoutService(driver CheckoutDriver) *CheckoutService {
	return &CheckoutService{driver: driver}
}

// CreateSession applies domain defaults before delegating to the configured driver.
func (s *CheckoutService) CreateSession(ctx context.Context, req payments.CheckoutSessionRequest) (payments.CheckoutSessionResult, error) {
	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	return s.driver.CreateSession(ctx, req)
}
