package responses

import "time"

type PlanConfigResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	PriceCents  int64    `json:"price_cents"`
	Description string   `json:"description"`
	Features    []string `json:"features"`
	SiteLimit   int      `json:"site_limit"`
}

type AdminUserResponse struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	ProfileImg string    `json:"profile_img"`
	Plan       string    `json:"plan"`
	Role       string    `json:"role"`
	Banned     bool      `json:"banned"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type AdminMetricPoint struct {
	Label string `json:"label"`
	Value int64  `json:"value"`
}

type AdminOverviewResponse struct {
	TotalUsers           int64              `json:"total_users"`
	ActiveUsers          int64              `json:"active_users"`
	BannedUsers          int64              `json:"banned_users"`
	PaidTotalCents       int64              `json:"paid_total_cents"`
	PaidLast30DaysCents  int64              `json:"paid_last_30_days_cents"`
	NewUsersSeries       []AdminMetricPoint `json:"new_users_series"`
	PaymentsSeries       []AdminMetricPoint `json:"payments_series"`
	RecentPayments       []AdminMetricPoint `json:"recent_payments"`
	NewUsersLast30Days   int64              `json:"new_users_last_30_days"`
	CreatedLast24h       int64              `json:"created_last_24h"`
	PaymentsLast24hCents int64              `json:"payments_last_24h_cents"`
	PlansByUsage         []AdminMetricPoint `json:"plans_by_usage"`
	AverageTicketCents   int64              `json:"average_ticket_cents"`
	PayingUsers          int64              `json:"paying_users"`
}
