package domain

import (
	"database/sql"
	"time"

	"github.com/ViitoJooj/go-sdk/validate"
	"github.com/google/uuid"
)

type User struct {
	UUID        uuid.UUID
	WebSiteUUID uuid.UUID
	ImageURL    *string
	Name        string
	Email       string
	Role        string
	Password    string
	CPF         *string
	GithubOauth bool
	GoogleOauth bool
	AppleOauth  bool
	UpdatedAt   *time.Time
	CreatedAt   time.Time
}

func NewUser(websiteUUID string, imageURL string, name string, email string, role string, password string, cpf string, github bool, google bool, apple bool, db *sql.DB) (*User, error) {

	if err := validate.UUIDv7(websiteUUID, "users", db); err != nil {
		return nil, err
	}

	if err := validate.FullName(name); err != nil {
		return nil, err
	}

	if err := validate.Email(email); err != nil {
		return nil, err
	}

	if role == "" {
		role = "user"
	}

	if err := validate.Password(password); err != nil {
		return nil, err
	}

	if cpf != "" {
		if err := validate.CPF(cpf); err != nil {
			return nil, err
		}
	}

	websiteUUIDParsed, err := uuid.Parse(websiteUUID)
	if err != nil {
		return nil, err
	}

	var imgPtr *string
	if imageURL != "" {
		imgPtr = &imageURL
	}

	return &User{
		UUID:        uuid.Nil,
		WebSiteUUID: websiteUUIDParsed,
		ImageURL:    imgPtr,
		Name:        name,
		Email:       email,
		Role:        role,
		Password:    password,
		CPF:         &cpf,
		GithubOauth: github,
		GoogleOauth: google,
		AppleOauth:  apple,
	}, nil
}
