package web

import (
	"net/http"
)

func (h *Handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/checkout", h.createCheckoutSession)
	mux.Handle("GET /", http.HandlerFunc(h.renderCheckoutPage))
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(h.fs))))
}
