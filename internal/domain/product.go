package domain

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	Id           string
	WebsiteId    string
	Name         string
	Description  string
	Price        float64
	ComparePrice *float64
	Stock        int
	Sku          *string
	Category     *string
	Images       []string
	Active       bool
	SoldCount    int
	CreatedBy    string
	UpdatedAt    time.Time
	CreatedAt    time.Time
}

func NewProduct(
	websiteId, name, description string,
	price float64, comparePrice *float64,
	stock int,
	sku, category *string,
	images []string,
	active bool,
	createdBy string,
) *Product {
	id, _ := uuid.NewV7()
	now := time.Now()
	if images == nil {
		images = []string{}
	}
	return &Product{
		Id:           id.String(),
		WebsiteId:    websiteId,
		Name:         name,
		Description:  description,
		Price:        price,
		ComparePrice: comparePrice,
		Stock:        stock,
		Sku:          sku,
		Category:     category,
		Images:       images,
		Active:       active,
		CreatedBy:    createdBy,
		UpdatedAt:    now,
		CreatedAt:    now,
	}
}
