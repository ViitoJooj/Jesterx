package models

import "time"

type UserPlan string

const (
	PlanFree       UserPlan = "free"
	PlanBusiness   UserPlan = "business"
	PlanPro        UserPlan = "pro"
	PlanEnterprise UserPlan = "enterprise"
)

type User struct {
	ID        string    `db:"id"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Plan      UserPlan  `db:"plan"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
