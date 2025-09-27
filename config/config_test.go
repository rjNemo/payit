package config

import (
	"strings"
	"testing"
)

func TestLoadSuccess(t *testing.T) {
	t.Setenv("PAYIT_STRIPE_SECRET_KEY", "sk_test")
	t.Setenv("PAYIT_STRIPE_PUBLISHABLE_KEY", "pk_test")
	t.Setenv("PAYIT_PRODUCT_NAME", "Demo product")
	t.Setenv("PAYIT_PRODUCT_DESCRIPTION", "Great product")
	t.Setenv("PAYIT_PRODUCT_PRICE_CENTS", "2500")
	t.Setenv("PAYIT_PRODUCT_CURRENCY", "usd")
	t.Setenv("PAYIT_PRODUCT_SUCCESS_URL", "https://example.com/success")
	t.Setenv("PAYIT_PRODUCT_CANCEL_URL", "https://example.com/cancel")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Product.PriceCents != 2500 {
		t.Fatalf("expected price 2500, got %d", cfg.Product.PriceCents)
	}
	if cfg.Product.Name != "Demo product" {
		t.Fatalf("unexpected product name: %s", cfg.Product.Name)
	}
}

func TestLoadMissingMandatoryVariables(t *testing.T) {
	clearAllEnv(t)

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when required variables are missing")
	}
	if !strings.Contains(err.Error(), "missing required environment variables") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadInvalidPrice(t *testing.T) {
	t.Setenv("PAYIT_STRIPE_SECRET_KEY", "sk_test")
	t.Setenv("PAYIT_STRIPE_PUBLISHABLE_KEY", "pk_test")
	t.Setenv("PAYIT_PRODUCT_NAME", "Demo product")
	t.Setenv("PAYIT_PRODUCT_DESCRIPTION", "Great product")
	t.Setenv("PAYIT_PRODUCT_PRICE_CENTS", "-1")
	t.Setenv("PAYIT_PRODUCT_CURRENCY", "usd")
	t.Setenv("PAYIT_PRODUCT_SUCCESS_URL", "https://example.com/success")
	t.Setenv("PAYIT_PRODUCT_CANCEL_URL", "https://example.com/cancel")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid price")
	}
	if !strings.Contains(err.Error(), "PAYIT_PRODUCT_PRICE_CENTS") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func clearAllEnv(t *testing.T) {
	t.Helper()
	envs := []string{
		"PAYIT_STRIPE_SECRET_KEY",
		"PAYIT_STRIPE_PUBLISHABLE_KEY",
		"PAYIT_PRODUCT_NAME",
		"PAYIT_PRODUCT_DESCRIPTION",
		"PAYIT_PRODUCT_PRICE_CENTS",
		"PAYIT_PRODUCT_CURRENCY",
		"PAYIT_PRODUCT_SUCCESS_URL",
		"PAYIT_PRODUCT_CANCEL_URL",
	}
	for _, env := range envs {
		t.Setenv(env, "")
	}
}
