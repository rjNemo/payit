package web

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	stripe "github.com/rjNemo/payit/internal/stripe"
)

func (h *Handler) createCheckoutSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var req stripe.CheckoutSessionRequest
	if r.Body != nil {
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		if err := dec.Decode(&req); err != nil {
			if errors.Is(err, io.EOF) {
				// Empty body is acceptable; default quantity applies.
			} else {
				http.Error(w, "invalid request payload", http.StatusBadRequest)
				return
			}
		}

		if dec.More() {
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
