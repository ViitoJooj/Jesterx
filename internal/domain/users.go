package domain

import (
	"time"

	"github.com/ViitoJooj/Jesterx/pkg/validators"
	"github.com/google/uuid"
)

type User struct {
	Uuid         string
	Avatar       *string
	Name         string
	Email        string
	Password     string
	Role         string
	Plan         *string
	Cpf          string
	Country_code *int
	Area_code    *int
	Phone        *int
	Address      Addresses
	Updated_at   time.Time
	Created_at   time.Time
}

type Addresses struct {
	Address_country    *string
	Zip_code           *string
	Address_street     *string
	Address_number     *string
	Address_complement *string
	Address_district   *string
	Address_city       *string
	Address_state      *string
}

func NewUser(name string, email string, password, role string, cpf string) (*User, error) {

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	uuidString := id.String()

	//validate data
	if err := validators.Uuid(uuidString); err != nil {
		return nil, err
	}

	if err := validators.Name(name); err != nil {
		return nil, err
	}

	if err := validators.Email(email); err != nil {
		return nil, err
	}

	if err := validators.Password(password); err != nil {
		return nil, err
	}

	if err := validators.Role(role); err != nil {
		return nil, err
	}

	return &User{
		Uuid:       uuidString,
		Name:       name,
		Email:      email,
		Password:   password,
		Role:       role,
		Cpf:        cpf,
		Updated_at: time.Now(),
		Created_at: time.Now(),
	}, nil
}
