package domain

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	Id               string
	WebsiteId        string
	Name             string
	Description      string
	ShortDescription *string
	Price            float64
	ComparePrice     *float64
	Stock            int
	Sku              *string
	Category         *string
	Slug             *string
	Brand            *string
	Model            *string
	Barcode          *string
	Condition        *string
	WeightGrams      *int
	WidthCm          *float64
	HeightCm         *float64
	LengthCm         *float64
	Material         *string
	Color            *string
	Size             *string
	WarrantyMonths   *int
	OriginCountry    *string
	Tags             []string
	Attributes       map[string]string
	RequiresShipping bool
	Images           []string
	Active           bool
	SoldCount        int
	CreatedBy        string
	UpdatedAt        time.Time
	CreatedAt        time.Time
}

func NewProduct(
	websiteId, name, description string,
	shortDescription *string,
	price float64, comparePrice *float64,
	stock int,
	sku, category, slug, brand, model, barcode, condition, material, color, size, originCountry *string,
	weightGrams, warrantyMonths *int,
	widthCm, heightCm, lengthCm *float64,
	tags []string,
	attributes map[string]string,
	requiresShipping bool,
	images []string,
	active bool,
	createdBy string,
) *Product {
	id, _ := uuid.NewV7()
	now := time.Now()
	if images == nil {
		images = []string{}
	}
	if tags == nil {
		tags = []string{}
	}
	if attributes == nil {
		attributes = map[string]string{}
	}
	return &Product{
		Id:               id.String(),
		WebsiteId:        websiteId,
		Name:             name,
		Description:      description,
		ShortDescription: shortDescription,
		Price:            price,
		ComparePrice:     comparePrice,
		Stock:            stock,
		Sku:              sku,
		Category:         category,
		Slug:             slug,
		Brand:            brand,
		Model:            model,
		Barcode:          barcode,
		Condition:        condition,
		WeightGrams:      weightGrams,
		WidthCm:          widthCm,
		HeightCm:         heightCm,
		LengthCm:         lengthCm,
		Material:         material,
		Color:            color,
		Size:             size,
		WarrantyMonths:   warrantyMonths,
		OriginCountry:    originCountry,
		Tags:             tags,
		Attributes:       attributes,
		RequiresShipping: requiresShipping,
		Images:           images,
		Active:           active,
		CreatedBy:        createdBy,
		UpdatedAt:        now,
		CreatedAt:        now,
	}
}
