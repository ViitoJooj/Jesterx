package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id         string
	First_name string
	Last_name  string
	Email      string
	Password   string
	Updated_at string
	Created_at string
}

func NewUser(first_name string, last_name string, email string, password_hash string) *User {
	layout := "00/00/0000 00:00"

	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}

	return &User{
		Id:         id.String(),
		First_name: first_name,
		Last_name:  last_name,
		Email:      email,
		Password:   password_hash,
		Updated_at: time.Now().Format(layout),
		Created_at: time.Now().Format(layout),
	}
}
