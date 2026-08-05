package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ViitoJooj/verkoupe/internal/domain/entities"
	"github.com/ViitoJooj/verkoupe/internal/domain/repositories/contracts"
	"github.com/ViitoJooj/verkoupe/internal/port/persistence/helpers"
)

var _ contracts.ProductContract = (*ProductRepository)(nil)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) CreateProduct(product *domain.Products) (*domain.Products, error) {
	if product == nil {
		return nil, errors.New("invalid product")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO products (name, description, short_description, active)
	VALUES ($1, $2, $3, $4)
	RETURNING uuid, created_at, updated_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		product.Name,
		product.Description,
		product.ShortDescription,
		product.Active,
	).Scan(
		&product.UUID,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, errors.New("could not create product")
	}

	return product, nil
}

func (r *ProductRepository) FindProductByUUID(uuid string) (*domain.Products, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, name, description, short_description, active, created_at, updated_at
	FROM products
	WHERE uuid = $1`

	row := r.db.QueryRowContext(ctx, query, uuid)
	return helpers.ScanProduct(row)
}

func (r *ProductRepository) FindProductByName(name string) (*domain.Products, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, name, description, short_description, active, created_at, updated_at
	FROM products
	WHERE name = $1`

	row := r.db.QueryRowContext(ctx, query, name)
	return helpers.ScanProduct(row)
}

func (r *ProductRepository) GetProducts() ([]*domain.Products, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, name, description, short_description, active, created_at, updated_at
	FROM products`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return helpers.ScanProducts(rows)
}

func (r *ProductRepository) GetActiveProducts() ([]*domain.Products, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, name, description, short_description, active, created_at, updated_at
	FROM products
	WHERE active = true`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return helpers.ScanProducts(rows)
}

func (r *ProductRepository) UpdateProductByUUID(uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE products SET updated_at = NOW() WHERE uuid = $1`

	result, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("product not found")
	}

	return nil
}

func (r *ProductRepository) DeleteProductByUUID(uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM products WHERE uuid = $1`

	result, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("product not found")
	}

	return nil
}

func (r *ProductRepository) DeleteProductsByUUIDS(uuids []string) error {
	if len(uuids) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	placeholders := make([]string, len(uuids))
	args := make([]interface{}, len(uuids))

	for i, u := range uuids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = u
	}

	query := fmt.Sprintf(`DELETE FROM products WHERE uuid IN (%s)`, strings.Join(placeholders, ", "))

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
