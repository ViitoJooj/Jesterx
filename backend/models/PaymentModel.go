package models

import "time"

type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusPaid     PaymentStatus = "paid"
	PaymentStatusFailed   PaymentStatus = "failed"
	PaymentStatusRefunded PaymentStatus = "refunded"
)

type Payment struct {
	ID                string        `db:"id"`
	UserID            string        `db:"user_id"`
	Plan              string        `db:"plan"`
	Provider          string        `db:"provider"`
	ProviderPaymentID string        `db:"provider_payment_id"`
	Status            PaymentStatus `db:"status"`
	AmountCents       int           `db:"amount_cents"`
	Currency          string        `db:"currency"`
	CreatedAt         time.Time     `db:"created_at"`
	UpdatedAt         time.Time     `db:"updated_at"`
}
