package domain

import (
	"log"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id         string
	First_name string
	Last_name  string
	Email      string
	Password   string
	Role       string
	Updated_at time.Time
	Created_at time.Time
}

func NewUser(first_name string, last_name string, email string, password_hash string) *User {
	id, err := uuid.NewV7()
	if err != nil {
		log.Fatal(err)
	}

	return &User{
		Id:         id.String(),
		First_name: first_name,
		Last_name:  last_name,
		Email:      email,
		Password:   password_hash,
		Role:       "user",
		Updated_at: time.Now(),
		Created_at: time.Now(),
	}
}
