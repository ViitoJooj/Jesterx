package domain

import "time"

type User_phones struct {
	Id         string
	User_id    string
	Phone      string
	Is_primary bool
	Updated_at time.Time
	Created_at time.Time
}
