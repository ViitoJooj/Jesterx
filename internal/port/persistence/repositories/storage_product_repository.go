package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ViitoJooj/Jesterx/internal/domain/entities"
	"github.com/ViitoJooj/Jesterx/internal/domain/repositories/contracts"
	"github.com/ViitoJooj/Jesterx/internal/port/persistence/helpers"
)

var _ contracts.StorageProductContract = (*StorageProductRepository)(nil)

type StorageProductRepository struct {
	db *sql.DB
}

func NewStorageProductRepository(db *sql.DB) *StorageProductRepository {
	return &StorageProductRepository{
		db: db,
	}
}

func (r *StorageProductRepository) CreateStorageProduct(sp *domain.StorageProducts) (*domain.StorageProducts, error) {
	if sp == nil {
		return nil, errors.New("invalid storage product")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO storage_products (product_uuid)
	VALUES ($1)`

	_, err := r.db.ExecContext(
		ctx,
		query,
		sp.ProductUUID,
	)

	if err != nil {
		return nil, errors.New("could not create storage product")
	}

	return sp, nil
}

func (r *StorageProductRepository) FindStorageProductByUUID(uuid string) (*domain.StorageProducts, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, product_uuid
	FROM storage_products
	WHERE uuid = $1`

	row := r.db.QueryRowContext(ctx, query, uuid)
	return helpers.ScanStorageProduct(row)
}

func (r *StorageProductRepository) FindStorageProductByProductUUID(productUUID string) (*domain.StorageProducts, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, product_uuid
	FROM storage_products
	WHERE product_uuid = $1`

	row := r.db.QueryRowContext(ctx, query, productUUID)
	return helpers.ScanStorageProduct(row)
}

func (r *StorageProductRepository) GetStorageProducts() ([]*domain.StorageProducts, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, product_uuid
	FROM storage_products`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return helpers.ScanStorageProducts(rows)
}

func (r *StorageProductRepository) UpdateStorageProductByUUID(uuid string) error {
	return nil
}

func (r *StorageProductRepository) DeleteStorageProductByUUID(uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM storage_products WHERE uuid = $1`

	result, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("storage product not found")
	}

	return nil
}

func (r *StorageProductRepository) DeleteStorageProductsByUUIDS(uuids []string) error {
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

	query := fmt.Sprintf(`DELETE FROM storage_products WHERE uuid IN (%s)`, strings.Join(placeholders, ", "))

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
