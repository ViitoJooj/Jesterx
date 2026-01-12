package domain

import "time"

type Sites struct {
	Id                string
	Owner_id          string
	Name              string
	Description       string
	Short_description string
	Domain            string
	Updated_at        time.Time
	Created_at        time.Time
}
