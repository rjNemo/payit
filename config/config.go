package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ProductConfig holds metadata for the single demo product.
type ProductConfig struct {
	Name        string
	Description string
	PriceCents  int64
	Currency    string
	SuccessURL  string
	CancelURL   string
}

// Config aggregates all runtime configuration required by the server.
type Config struct {
	StripeSecretKey      string
	StripePublishableKey string
	Product              ProductConfig
}

// Load reads configuration from environment variables, optionally sourcing
// a .env.local file when present.
func Load() (Config, error) {
	_ = loadDotEnv()

	priceRaw := strings.TrimSpace(os.Getenv("PAYIT_PRODUCT_PRICE_CENTS"))
	cfg := Config{
		StripeSecretKey:      os.Getenv("PAYIT_STRIPE_SECRET_KEY"),
		StripePublishableKey: os.Getenv("PAYIT_STRIPE_PUBLISHABLE_KEY"),
		Product: ProductConfig{
			Name:        os.Getenv("PAYIT_PRODUCT_NAME"),
			Description: os.Getenv("PAYIT_PRODUCT_DESCRIPTION"),
			Currency:    os.Getenv("PAYIT_PRODUCT_CURRENCY"),
			SuccessURL:  os.Getenv("PAYIT_PRODUCT_SUCCESS_URL"),
			CancelURL:   os.Getenv("PAYIT_PRODUCT_CANCEL_URL"),
		},
	}

	if missing := validate(cfg, priceRaw); len(missing) > 0 {
		return Config{}, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	price, err := parsePrice(priceRaw)
	if err != nil {
		return Config{}, err
	}
	cfg.Product.PriceCents = price

	return cfg, nil
}

func loadDotEnv() error {
	filename := ".env.local"
	relPaths := []string{
		filename,
		filepath.Join("..", filename),
	}

	for _, path := range relPaths {
		if err := applyEnvFile(path); err != nil {
			return err
		}
	}
	return nil
}

func applyEnvFile(path string) (retErr error) {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer func() {
		if cerr := file.Close(); retErr == nil && cerr != nil {
			retErr = cerr
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func parsePrice(value string) (int64, error) {
	price, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("PAYIT_PRODUCT_PRICE_CENTS must be a positive integer: %w", err)
	}
	if price <= 0 {
		return 0, fmt.Errorf("PAYIT_PRODUCT_PRICE_CENTS must be a positive integer")
	}
	return price, nil
}

func validate(cfg Config, priceRaw string) []string {
	missing := make([]string, 0)
	if cfg.StripeSecretKey == "" {
		missing = append(missing, "PAYIT_STRIPE_SECRET_KEY")
	}
	if cfg.StripePublishableKey == "" {
		missing = append(missing, "PAYIT_STRIPE_PUBLISHABLE_KEY")
	}
	if cfg.Product.Name == "" {
		missing = append(missing, "PAYIT_PRODUCT_NAME")
	}
	if cfg.Product.Description == "" {
		missing = append(missing, "PAYIT_PRODUCT_DESCRIPTION")
	}
	if priceRaw == "" {
		missing = append(missing, "PAYIT_PRODUCT_PRICE_CENTS")
	}
	if cfg.Product.Currency == "" {
		missing = append(missing, "PAYIT_PRODUCT_CURRENCY")
	}
	if cfg.Product.SuccessURL == "" {
		missing = append(missing, "PAYIT_PRODUCT_SUCCESS_URL")
	}
	if cfg.Product.CancelURL == "" {
		missing = append(missing, "PAYIT_PRODUCT_CANCEL_URL")
	}

	return missing
}
