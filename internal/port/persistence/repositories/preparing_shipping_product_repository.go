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

var _ contracts.PreparingShippingProductContract = (*PreparingShippingProductRepository)(nil)

type PreparingShippingProductRepository struct {
	db *sql.DB
}

func NewPreparingShippingProductRepository(db *sql.DB) *PreparingShippingProductRepository {
	return &PreparingShippingProductRepository{
		db: db,
	}
}

func (r *PreparingShippingProductRepository) CreatePreparingShippingProduct(psp *domain.PreparingShippingProducts) (*domain.PreparingShippingProducts, error) {
	if psp == nil {
		return nil, errors.New("invalid preparing shipping product")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO preparing_shipping_products (product_uuid, address_uuid)
	VALUES ($1, $2)`

	_, err := r.db.ExecContext(
		ctx,
		query,
		psp.ProductUUID,
		psp.AddressUUID,
	)

	if err != nil {
		return nil, errors.New("could not create preparing shipping product")
	}

	return psp, nil
}

func (r *PreparingShippingProductRepository) FindPreparingShippingProductByUUID(uuid string) (*domain.PreparingShippingProducts, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, product_uuid, address_uuid
	FROM preparing_shipping_products
	WHERE uuid = $1`

	row := r.db.QueryRowContext(ctx, query, uuid)
	return helpers.ScanPreparingShippingProduct(row)
}

func (r *PreparingShippingProductRepository) FindPreparingShippingProductByProductUUID(productUUID string) (*domain.PreparingShippingProducts, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, product_uuid, address_uuid
	FROM preparing_shipping_products
	WHERE product_uuid = $1`

	row := r.db.QueryRowContext(ctx, query, productUUID)
	return helpers.ScanPreparingShippingProduct(row)
}

func (r *PreparingShippingProductRepository) GetPreparingShippingProducts() ([]*domain.PreparingShippingProducts, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, product_uuid, address_uuid
	FROM preparing_shipping_products`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return helpers.ScanPreparingShippingProducts(rows)
}

func (r *PreparingShippingProductRepository) UpdatePreparingShippingProductByUUID(uuid string) error {
	return nil
}

func (r *PreparingShippingProductRepository) DeletePreparingShippingProductByUUID(uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM preparing_shipping_products WHERE uuid = $1`

	result, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("preparing shipping product not found")
	}

	return nil
}

func (r *PreparingShippingProductRepository) DeletePreparingShippingProductsByUUIDS(uuids []string) error {
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

	query := fmt.Sprintf(`DELETE FROM preparing_shipping_products WHERE uuid IN (%s)`, strings.Join(placeholders, ", "))

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
