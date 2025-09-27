package service

import (
	"context"
	"errors"
	"testing"

	"github.com/rjNemo/payit/internal/payments"
)

type fakeDriver struct {
	lastReq payments.CheckoutSessionRequest
	result  payments.CheckoutSessionResult
	err     error
}

func (f *fakeDriver) CreateSession(ctx context.Context, req payments.CheckoutSessionRequest) (payments.CheckoutSessionResult, error) {
	f.lastReq = req
	if f.err != nil {
		return payments.CheckoutSessionResult{}, f.err
	}
	return f.result, nil
}

func TestCheckoutService_DefaultQuantity(t *testing.T) {
	drv := &fakeDriver{}
	svc := NewCheckoutService(drv)

	_, err := svc.CreateSession(context.Background(), payments.CheckoutSessionRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if drv.lastReq.Quantity != 1 {
		t.Fatalf("expected default quantity 1, got %d", drv.lastReq.Quantity)
	}
}

func TestCheckoutService_PreservesQuantity(t *testing.T) {
	drv := &fakeDriver{}
	svc := NewCheckoutService(drv)

	_, err := svc.CreateSession(context.Background(), payments.CheckoutSessionRequest{Quantity: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if drv.lastReq.Quantity != 5 {
		t.Fatalf("expected quantity 5, got %d", drv.lastReq.Quantity)
	}
}

func TestCheckoutService_PropagatesError(t *testing.T) {
	drv := &fakeDriver{err: errors.New("driver failed")}
	svc := NewCheckoutService(drv)

	_, err := svc.CreateSession(context.Background(), payments.CheckoutSessionRequest{Quantity: 2})
	if err == nil {
		t.Fatal("expected error from driver")
	}
}
