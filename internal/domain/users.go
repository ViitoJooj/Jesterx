package domain

import (
	"time"

	"github.com/ViitoJooj/Jesterx/pkg/validators"
	"github.com/google/uuid"
)

type Auth struct {
	Uid                string
	DisplayName        string
	Email              string
	Phone              *string
	Providers          []string
	ProviderType       string
	AvatarURL          *string
	Email_confirmed_at *time.Time
	Created_at         time.Time
	Last_sign_inAt     time.Time
}

type Profile struct {
	Uid         uuid.UUID
	WebsiteID   uuid.UUID
	Name        string
	Email       string
	Password    string
	Avatar      *string
	Role        string
	Plan        *string
	Cpf         *string
	CountryCode *int
	AreaCode    *int
	Address     Address
	Updated_at  time.Time
	Created_at  time.Time
}

func (u *Profile) NewUser(name, email, password, role, cpf string) (*Profile, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	if err := validators.User(name, email, password, role, cpf); err != nil {
		return nil, err
	}

	cpfVal := cpf

	return &Profile{
		Uid:        id,
		Name:       name,
		Email:      email,
		Password:   password,
		Role:       role,
		Cpf:        &cpfVal,
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}, nil
}
