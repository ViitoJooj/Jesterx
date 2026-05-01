package domain

import (
	"time"

	"github.com/ViitoJooj/Jesterx/pkg/validators"
	"github.com/ViitoJooj/Jesterx/pkg/validators/users_validations"
	"github.com/google/uuid"
)

type Users struct {
	Uid                uuid.UUID
	WebsiteID          uuid.UUID
	AvatarURL          *string
	Name               string
	Email              string
	Email_confirmed_at time.Time
	Phone              string
	Password           string
	Avatar             *string
	Role               string
	Plan               *string
	Cpf                *string
	CountryCode        *int
	AreaCode           *int
	Address            Address
	Updated_at         time.Time
	Created_at         time.Time
}

func (u *Users) NewUser(websiteID uuid.UUID, name, email, password, role, cpf string) (*Users, error) {

	uid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	// validate data
	if err := validators.Uuid(uid); err != nil {
		return nil, err
	}

	if err := validators.Uuid(websiteID); err != nil {
		return nil, err
	}

	if err := users_validations.Name(name); err != nil {
		return nil, err
	}

	if err := users_validations.Email(email); err != nil {
		return nil, err
	}

	if err := users_validations.Password(password); err != nil {
		return nil, err
	}

	if err := users_validations.Role(role); err != nil {
		return nil, err
	}

	cpfVal := cpf

	return &Users{
		Uid:        uid,
		WebsiteID:  websiteID,
		Name:       name,
		Email:      email,
		Password:   password,
		Role:       role,
		Cpf:        &cpfVal,
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}, nil
}
