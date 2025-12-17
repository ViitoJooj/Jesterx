package config

import (
	"os"

	"github.com/stripe/stripe-go/v82"
)

var (
	StripeSuccessURL string
	StripeCancelURL  string
)

func InitStripe() {
	stripe.Key = StripeSecretKey

	StripeSuccessURL = os.Getenv("STRIPE_SUCCESS_URL")
	if StripeSuccessURL == "" {
		StripeSuccessURL = "http://localhost:3000/payment-success"
	}

	StripeCancelURL = os.Getenv("STRIPE_CANCEL_URL")
	if StripeCancelURL == "" {
		StripeCancelURL = "http://localhost:3000/payment-cancel"
	}
}
