package web

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	stripe "github.com/rjNemo/payit/internal/stripe"
)

type fakeCheckoutService struct {
	result stripe.CheckoutSessionResult
	err    error
	req    stripe.CheckoutSessionRequest
}

func (f *fakeCheckoutService) CreateSession(ctx context.Context, req stripe.CheckoutSessionRequest) (stripe.CheckoutSessionResult, error) {
	f.req = req
	if f.err != nil {
		return stripe.CheckoutSessionResult{}, f.err
	}
	return f.result, nil
}

func TestCreateCheckoutSessionSuccess(t *testing.T) {
	handler := &Handler{
		checkout: &fakeCheckoutService{
			result: stripe.CheckoutSessionResult{ID: "cs_test_1", URL: "https://stripe.test/checkout"},
		},
	}

	body, _ := json.Marshal(map[string]int{"quantity": 2})
	req := httptest.NewRequest(http.MethodPost, "/api/checkout", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.createCheckoutSession(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected json content type, got %s", ct)
	}

	var payload stripe.CheckoutSessionResult
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("expected valid json response: %v", err)
	}
	if payload.ID != "cs_test_1" || payload.URL == "" {
		t.Fatalf("unexpected payload: %#v", payload)
	}

	svc := handler.checkout.(*fakeCheckoutService)
	if svc.req.Quantity != 2 {
		t.Fatalf("expected quantity 2, got %d", svc.req.Quantity)
	}
}

func TestCreateCheckoutSessionDefaultsQuantity(t *testing.T) {
	fakeSvc := &fakeCheckoutService{
		result: stripe.CheckoutSessionResult{ID: "cs_test_1"},
	}
	handler := &Handler{checkout: fakeSvc}

	req := httptest.NewRequest(http.MethodPost, "/api/checkout", http.NoBody)
	rec := httptest.NewRecorder()

	handler.createCheckoutSession(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if fakeSvc.req.Quantity != 0 {
		t.Fatalf("expected zero quantity in request, got %d", fakeSvc.req.Quantity)
	}
}

func TestCreateCheckoutSessionRejectsInvalidJSON(t *testing.T) {
	handler := &Handler{checkout: &fakeCheckoutService{}}

	req := httptest.NewRequest(http.MethodPost, "/api/checkout", bytes.NewBufferString("{"))
	rec := httptest.NewRecorder()

	handler.createCheckoutSession(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestCreateCheckoutSessionStripeFailure(t *testing.T) {
	handler := &Handler{
		checkout: &fakeCheckoutService{err: errors.New("stripe failure")},
	}

	req := httptest.NewRequest(http.MethodPost, "/api/checkout", http.NoBody)
	rec := httptest.NewRecorder()

	handler.createCheckoutSession(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}
}

func TestCreateCheckoutSessionMethodNotAllowed(t *testing.T) {
	handler := &Handler{checkout: &fakeCheckoutService{}}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/checkout", handler.createCheckoutSession)

	req := httptest.NewRequest(http.MethodGet, "/api/checkout", http.NoBody)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rec.Code)
	}
	if allow := rec.Header().Get("Allow"); allow != http.MethodPost {
		t.Fatalf("expected Allow header to be POST, got %s", allow)
	}
}
