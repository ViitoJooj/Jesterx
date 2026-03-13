package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/ViitoJooj/Jesterx/internal/domain"
)

func NewProductRepository(db *sql.DB) *connection {
	return &connection{db: db}
}

func scanProduct(row interface {
	Scan(...any) error
}) (*domain.Product, error) {
	var p domain.Product
	var imagesJSON []byte
	var tagsJSON []byte
	var attrsJSON []byte
	err := row.Scan(
		&p.Id,
		&p.WebsiteId,
		&p.Name,
		&p.Description,
		&p.ShortDescription,
		&p.Price,
		&p.ComparePrice,
		&p.Stock,
		&p.Sku,
		&p.Category,
		&p.Slug,
		&p.Brand,
		&p.Model,
		&p.Barcode,
		&p.Condition,
		&p.WeightGrams,
		&p.WidthCm,
		&p.HeightCm,
		&p.LengthCm,
		&p.Material,
		&p.Color,
		&p.Size,
		&p.WarrantyMonths,
		&p.OriginCountry,
		&tagsJSON,
		&attrsJSON,
		&p.RequiresShipping,
		&imagesJSON,
		&p.Active,
		&p.SoldCount,
		&p.CreatedBy,
		&p.UpdatedAt,
		&p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(imagesJSON, &p.Images); err != nil {
		p.Images = []string{}
	}
	if err := json.Unmarshal(tagsJSON, &p.Tags); err != nil {
		p.Tags = []string{}
	}
	if err := json.Unmarshal(attrsJSON, &p.Attributes); err != nil {
		p.Attributes = map[string]string{}
	}
	return &p, nil
}

const productCols = `id, website_id, name, description, short_description, price, compare_price, stock, sku, category, slug, brand, model, barcode, condition, weight_grams, width_cm, height_cm, length_cm, material, color, size, warranty_months, origin_country, tags, attributes, requires_shipping, images, active, sold_count, created_by, updated_at, created_at`

func (r *connection) CreateProduct(p domain.Product) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	imagesJSON, _ := json.Marshal(p.Images)
	tagsJSON, _ := json.Marshal(p.Tags)
	attrsJSON, _ := json.Marshal(p.Attributes)
	row := r.db.QueryRowContext(ctx, `
		INSERT INTO products (
			id, website_id, name, description, short_description, price, compare_price, stock, sku, category,
			slug, brand, model, barcode, condition, weight_grams, width_cm, height_cm, length_cm, material,
			color, size, warranty_months, origin_country, tags, attributes, requires_shipping,
			images, active, sold_count, created_by, updated_at, created_at
		)
		VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,
			$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,
			$21,$22,$23,$24,$25,$26,$27,
			$28,$29,$30,$31,$32,$33
		)
		RETURNING `+productCols,
		p.Id, p.WebsiteId, p.Name, p.Description, p.ShortDescription, p.Price, p.ComparePrice,
		p.Stock, p.Sku, p.Category, p.Slug, p.Brand, p.Model, p.Barcode, p.Condition,
		p.WeightGrams, p.WidthCm, p.HeightCm, p.LengthCm, p.Material,
		p.Color, p.Size, p.WarrantyMonths, p.OriginCountry, tagsJSON, attrsJSON, p.RequiresShipping,
		imagesJSON, p.Active, p.SoldCount, p.CreatedBy, p.UpdatedAt, p.CreatedAt,
	)
	return scanProduct(row)
}

func (r *connection) FindProductByID(id, websiteId string) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(ctx,
		`SELECT `+productCols+` FROM products WHERE id=$1 AND website_id=$2`,
		id, websiteId,
	)
	p, err := scanProduct(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return p, err
}

func (r *connection) ListProductsByWebsiteID(websiteId string, onlyActive bool) ([]domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `SELECT ` + productCols + ` FROM products WHERE website_id=$1`
	if onlyActive {
		query += ` AND active=true`
	}
	query += ` ORDER BY sold_count DESC, created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, websiteId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		p, err := scanProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *p)
	}
	if products == nil {
		products = []domain.Product{}
	}
	return products, rows.Err()
}

func (r *connection) UpdateProduct(p domain.Product) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	imagesJSON, _ := json.Marshal(p.Images)
	tagsJSON, _ := json.Marshal(p.Tags)
	attrsJSON, _ := json.Marshal(p.Attributes)
	row := r.db.QueryRowContext(ctx, `
		UPDATE products SET
			name=$1, description=$2, short_description=$3, price=$4, compare_price=$5,
			stock=$6, sku=$7, category=$8, slug=$9, brand=$10, model=$11, barcode=$12,
			condition=$13, weight_grams=$14, width_cm=$15, height_cm=$16, length_cm=$17,
			material=$18, color=$19, size=$20, warranty_months=$21, origin_country=$22,
			tags=$23, attributes=$24, requires_shipping=$25, images=$26, active=$27, updated_at=NOW()
		WHERE id=$28 AND website_id=$29
		RETURNING `+productCols,
		p.Name, p.Description, p.ShortDescription, p.Price, p.ComparePrice,
		p.Stock, p.Sku, p.Category, p.Slug, p.Brand, p.Model, p.Barcode,
		p.Condition, p.WeightGrams, p.WidthCm, p.HeightCm, p.LengthCm,
		p.Material, p.Color, p.Size, p.WarrantyMonths, p.OriginCountry,
		tagsJSON, attrsJSON, p.RequiresShipping, imagesJSON, p.Active,
		p.Id, p.WebsiteId,
	)
	updated, err := scanProduct(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("product not found")
	}
	return updated, err
}

func (r *connection) DeleteProduct(id, websiteId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := r.db.ExecContext(ctx,
		`DELETE FROM products WHERE id=$1 AND website_id=$2`, id, websiteId,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("product not found")
	}
	return nil
}
