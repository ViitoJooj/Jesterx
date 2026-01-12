package domain

import "time"

type Pages struct {
	id          string
	site_id     string
	tittle      string
	slug        string
	content     string
	user_header bool
	user_footer bool
	Updated_at  time.Time
	Created_at  time.Time
}
