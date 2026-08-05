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

var _ contracts.ProductShippedContract = (*ProductShippedRepository)(nil)

type ProductShippedRepository struct {
	db *sql.DB
}

func NewProductShippedRepository(db *sql.DB) *ProductShippedRepository {
	return &ProductShippedRepository{
		db: db,
	}
}

func (r *ProductShippedRepository) CreateProductShipped(productShipped *domain.ProductShipped) (*domain.ProductShipped, error) {
	if productShipped == nil {
		return nil, errors.New("invalid product shipped")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO products_shipped (product_uuid, address_uuid, status)
	VALUES ($1, $2, $3)
	RETURNING uuid`

	err := r.db.QueryRowContext(
		ctx,
		query,
		productShipped.ProductUUID,
		productShipped.AddressUUID,
		productShipped.Status,
	).Scan(
		&productShipped.UUID,
	)

	if err != nil {
		return nil, errors.New("could not create product shipped")
	}

	return productShipped, nil
}

func (r *ProductShippedRepository) FindProductShippedByUUID(uuid string) (*domain.ProductShipped, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, product_uuid, address_uuid, status
	FROM products_shipped
	WHERE uuid = $1`

	row := r.db.QueryRowContext(ctx, query, uuid)
	return helpers.ScanProductShipped(row)
}

func (r *ProductShippedRepository) FindProductShippedByProductUUID(productUUID string) ([]*domain.ProductShipped, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, product_uuid, address_uuid, status
	FROM products_shipped
	WHERE product_uuid = $1`

	rows, err := r.db.QueryContext(ctx, query, productUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return helpers.ScanProductsShipped(rows)
}

func (r *ProductShippedRepository) FindProductShippedByStatus(status string) ([]*domain.ProductShipped, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, product_uuid, address_uuid, status
	FROM products_shipped
	WHERE status = $1`

	rows, err := r.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return helpers.ScanProductsShipped(rows)
}

func (r *ProductShippedRepository) GetProductsShipped() ([]*domain.ProductShipped, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, product_uuid, address_uuid, status
	FROM products_shipped`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return helpers.ScanProductsShipped(rows)
}

func (r *ProductShippedRepository) UpdateProductShippedByUUID(uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE products_shipped SET uuid = uuid WHERE uuid = $1`

	result, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("product shipped not found")
	}

	return nil
}

func (r *ProductShippedRepository) DeleteProductShippedByUUID(uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM products_shipped WHERE uuid = $1`

	result, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("product shipped not found")
	}

	return nil
}

func (r *ProductShippedRepository) DeleteProductsShippedByUUIDS(uuids []string) error {
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

	query := fmt.Sprintf(`DELETE FROM products_shipped WHERE uuid IN (%s)`, strings.Join(placeholders, ", "))

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
