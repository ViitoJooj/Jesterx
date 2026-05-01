package validators

import (
	"errors"
	"strings"
)

func NewPlan(name, currency string, price int) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("plan name is required.")
	}

	if price < 0 {
		return errors.New("price must be non-negative.")
	}

	if strings.TrimSpace(currency) == "" {
		return errors.New("currency is required.")
	}

	return nil
}
