package web

import (
	"net/http"

	"github.com/rjNemo/payit/config"
)

// Handler aggregates dependencies required by HTTP handlers.
type Handler struct {
	cfg config.Config
}

// NewServer constructs the root HTTP handler. The initial implementation only
// exposes placeholder routes; later phases will wire Stripe-backed handlers and
// templates.
func NewServer(cfg config.Config) http.Handler {
	h := &Handler{cfg: cfg}

	mux := http.NewServeMux()
	mux.HandleFunc("/", h.notImplemented)

	return mux
}

func (h *Handler) notImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "PayIt demo coming soon", http.StatusNotImplemented)
}
