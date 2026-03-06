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
	err := row.Scan(
		&p.Id,
		&p.WebsiteId,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.ComparePrice,
		&p.Stock,
		&p.Sku,
		&p.Category,
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
	return &p, nil
}

const productCols = `id, website_id, name, description, price, compare_price, stock, sku, category, images, active, sold_count, created_by, updated_at, created_at`

func (r *connection) CreateProduct(p domain.Product) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	imagesJSON, _ := json.Marshal(p.Images)
	row := r.db.QueryRowContext(ctx, `
		INSERT INTO products (id, website_id, name, description, price, compare_price, stock, sku, category, images, active, sold_count, created_by, updated_at, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
		RETURNING `+productCols,
		p.Id, p.WebsiteId, p.Name, p.Description, p.Price, p.ComparePrice,
		p.Stock, p.Sku, p.Category, imagesJSON, p.Active, p.SoldCount, p.CreatedBy, p.UpdatedAt, p.CreatedAt,
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
	row := r.db.QueryRowContext(ctx, `
		UPDATE products SET
			name=$1, description=$2, price=$3, compare_price=$4,
			stock=$5, sku=$6, category=$7, images=$8, active=$9, updated_at=NOW()
		WHERE id=$10 AND website_id=$11
		RETURNING `+productCols,
		p.Name, p.Description, p.Price, p.ComparePrice,
		p.Stock, p.Sku, p.Category, imagesJSON, p.Active,
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
