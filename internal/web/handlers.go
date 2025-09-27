package web

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/rjNemo/payit/internal/payments"
)

func (h *Handler) createCheckoutSession(w http.ResponseWriter, r *http.Request) {
	var req payments.CheckoutSessionRequest

	if r.Body != nil {
		defer func(body io.ReadCloser) {
			_ = body.Close()
		}(r.Body)
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		if err := dec.Decode(&req); err != nil {
			if errors.Is(err, io.EOF) {
				// Empty body is acceptable; default quantity applies.
			} else {
				http.Error(w, "invalid request payload", http.StatusBadRequest)
				return
			}
		} else if dec.More() {
			http.Error(w, "unexpected data in request body", http.StatusBadRequest)
			return
		}
	}

	session, err := h.checkout.CreateSession(r.Context(), req)
	if err != nil {
		http.Error(w, "checkout session failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(session); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
