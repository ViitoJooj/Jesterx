package domain

import "time"

type Plan struct {
	ID           int64
	Name         string
	Description  string
	DescriptionM string
	Price        float64
	BillingCycle string
	Active       bool
	MaxSites     int
	MaxRoutes    int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Payment struct {
	ID          int64
	UserID      string
	WebsiteID   string
	ReferenceID string
	Type        string
	Quantity    int
	Amount      float64
	Currency    string
	Status      string
	PurchasedIn time.Time
}
