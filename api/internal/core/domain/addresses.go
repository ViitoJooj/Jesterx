package domain

import "time"

type Addresses struct {
	Id           string
	User_id      string
	Country_code string
	Stade        string
	City         string
	Postal_code  string
	Street       string
	Number       string
	Complement   string
	Updated_at   time.Time
	Created_at   time.Time
}
