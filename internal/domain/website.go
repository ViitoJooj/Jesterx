package domain

import (
	"errors"
	"log"
	"time"

	"github.com/ViitoJooj/Jesterx/pkg/validators"
	"github.com/ViitoJooj/Jesterx/pkg/validators/website_validations"
	"github.com/google/uuid"
)

type WebsiteID struct {
	Uuid              uuid.UUID
	Owner_id          uuid.UUID
	Name              string
	Website_type      string
	Logo              *string
	Short_description *string
	Description       *string
	Banned            bool
	Mature_content    bool
	Rating_count      int
	Updated_at        time.Time
	Created_at        time.Time
}

func (u *WebsiteID) NewWebsite(owner_id uuid.UUID, name string, website_type string, logo string, short_description string, description string, banned bool, mature_content bool, Rating_count int) (*WebsiteID, error) {
	uid, err := uuid.NewV7()
	if err != nil {
		log.Println(err)
		return nil, errors.New("Internal error.")
	}

	if err := validators.Uuid(uid); err != nil {
		log.Println(err)
		return nil, errors.New("Internal error.")
	}

	if err := website_validations.Name(name); err != nil {
		log.Println(err)
		return nil, err
	}

}
