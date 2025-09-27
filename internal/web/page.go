package web

import (
	"fmt"
	"net/http"
	"strings"
)

type checkoutPageData struct {
	ProductName        string
	ProductDescription string
	PriceDisplay       string
	Currency           string
}

func (h *Handler) renderCheckoutPage(w http.ResponseWriter, r *http.Request) {
	price := float64(h.cfg.Product.PriceCents) / 100
	data := checkoutPageData{
		ProductName:        h.cfg.Product.Name,
		ProductDescription: h.cfg.Product.Description,
		PriceDisplay:       fmt.Sprintf("$%.2f", price),
		Currency:           strings.ToUpper(h.cfg.Product.Currency),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.page.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
		return
	}
}
