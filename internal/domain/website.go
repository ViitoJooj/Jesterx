package domain

import (
	"time"

	"github.com/google/uuid"
)

type WebSite struct {
	Id                string
	Type              string
	Image             []byte
	Name              string
	Short_description string
	Description       string
	Creator_id        string
	Banned            bool
	Updated_at        time.Time
	Created_at        time.Time
}

func NewWebSite(Type string, Image []byte, Name string, Short_description string, Description string, Creator_id string) *WebSite {
	id, _ := uuid.NewV7()

	return &WebSite{
		Id:                id.String(),
		Type:              Type,
		Image:             Image,
		Name:              Name,
		Short_description: Short_description,
		Description:       Description,
		Creator_id:        Creator_id,
		Banned:            false,
		Updated_at:        time.Now(),
		Created_at:        time.Now(),
	}
}
