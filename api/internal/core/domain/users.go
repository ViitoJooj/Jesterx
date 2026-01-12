package domain

import "time"

type User struct {
	ID         string
	FirstName  string
	LastName   string
	Email      string
	Password   string
	Role       string
	Updated_at time.Time
	Created_at time.Time
}
