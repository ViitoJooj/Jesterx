package domain

import "time"

type Refresh_tokens struct {
	Id           string
	User_id      string
	Token        string
	Expires_at   time.Time
	Revoke_at    time.Time
	Created_at   time.Time
	Ip_address   string
	User_agent   string
	Rotated_from string
}
