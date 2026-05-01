package domain

import (
	"time"

	"github.com/google/uuid"
)

type WebsiteID struct {
	Uuid              uuid.UUID
	Owner_id          uuid.UUID
	Name              string
	Website_type      string
	Logo              *string
	Short_description string
	Description       string
	Banned            bool
	Mature_content    bool
	Rating_count      int
	Updated_at        time.Time
	Created_at        time.Time
}
