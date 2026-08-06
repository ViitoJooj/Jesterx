package domain

import (
	"errors"
	"time"

	"github.com/ViitoJooj/go-sdk/validate"
	"github.com/google/uuid"
)

type OrganizationBR struct {
	UUID        uuid.UUID
	WebSiteUUID uuid.UUID
	OwnerUUID   uuid.UUID
	ImageURL    string
	Name        string
	TradeName   string
	CNPJ        string
	UpdatedAt   *time.Time
	CreatedAt   time.Time
}

func NewOrganization(websiteUUID string, ownerUUID string, imageURL string, name string, tradeName string, cnpj string) (*OrganizationBR, error) {
	if err := validate.FullName(name); err != nil {
		return nil, err
	}

	if tradeName == "" {
		return nil, errors.New("TradeName cannot be null.")
	}

	if err := validate.CNPJ(cnpj); err != nil {
		return nil, err
	}

	websiteUUIDParsed, err := uuid.Parse(websiteUUID)
	if err != nil {
		return nil, err
	}

	ownerUUIDParsed, err := uuid.Parse(ownerUUID)
	if err != nil {
		return nil, err
	}

	return &OrganizationBR{
		UUID:        uuid.Nil,
		WebSiteUUID: websiteUUIDParsed,
		OwnerUUID:   ownerUUIDParsed,
		ImageURL:    imageURL,
		Name:        name,
		TradeName:   tradeName,
		CNPJ:        cnpj,
	}, nil
}
