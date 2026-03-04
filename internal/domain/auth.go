package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id             string
	WebsiteId      string
	First_name     string
	Last_name      string
	Email          string
	Verified_email bool
	Password       string
	Role           string
	Plan           *string
	Updated_at     time.Time
	Created_at     time.Time
}

func NewUser(WebsiteId string, first_name string, last_name string, email string, password_hash string) *User {
	id, _ := uuid.NewV7()

	return &User{
		Id:             id.String(),
		WebsiteId:      WebsiteId,
		First_name:     first_name,
		Last_name:      last_name,
		Email:          email,
		Verified_email: false,
		Password:       password_hash,
		Role:           "user",
		Updated_at:     time.Now(),
		Created_at:     time.Now(),
	}
}
