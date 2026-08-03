package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ViitoJooj/Jesterx/internal/domain"
	"github.com/ViitoJooj/Jesterx/internal/repositories/contracts"
	"github.com/ViitoJooj/Jesterx/internal/repositories/helpers"
)

var _ contracts.ProductTagContract = (*ProductTagRepository)(nil)

type ProductTagRepository struct {
	db *sql.DB
}

func NewProductTagRepository(db *sql.DB) *ProductTagRepository {
	return &ProductTagRepository{
		db: db,
	}
}

func (r *ProductTagRepository) CreateProductTag(tag *domain.ProductsTags) (*domain.ProductsTags, error) {
	if tag == nil {
		return nil, errors.New("invalid product tag")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO products_tags (product_uuid, label)
	VALUES ($1, $2)
	RETURNING uuid, created_at, updated_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		tag.ProductUUID,
		tag.Label,
	).Scan(
		&tag.UUID,
		&tag.CreatedAt,
		&tag.UpdatedAt,
	)

	if err != nil {
		return nil, errors.New("could not create product tag")
	}

	return tag, nil
}

func (r *ProductTagRepository) FindProductTagByUUID(uuid string) (*domain.ProductsTags, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, product_uuid, label, created_at, updated_at
	FROM products_tags
	WHERE uuid = $1`

	row := r.db.QueryRowContext(ctx, query, uuid)
	return helpers.ScanProductTag(row)
}

func (r *ProductTagRepository) FindProductTagsByLabel(label string) ([]*domain.ProductsTags, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, product_uuid, label, created_at, updated_at
	FROM products_tags
	WHERE label = $1`

	rows, err := r.db.QueryContext(ctx, query, label)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return helpers.ScanProductTags(rows)
}

func (r *ProductTagRepository) GetProductTagsFromProduct(productUUID string) ([]*domain.ProductsTags, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, product_uuid, label, created_at, updated_at
	FROM products_tags
	WHERE product_uuid = $1`

	rows, err := r.db.QueryContext(ctx, query, productUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return helpers.ScanProductTags(rows)
}

func (r *ProductTagRepository) GetProductTags() ([]*domain.ProductsTags, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, product_uuid, label, created_at, updated_at
	FROM products_tags`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return helpers.ScanProductTags(rows)
}

func (r *ProductTagRepository) UpdateProductTagByUUID(uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE products_tags SET updated_at = NOW() WHERE uuid = $1`

	result, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("product tag not found")
	}

	return nil
}

func (r *ProductTagRepository) DeleteProductTagByUUID(uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM products_tags WHERE uuid = $1`

	result, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("product tag not found")
	}

	return nil
}

func (r *ProductTagRepository) DeleteProductTagsByUUIDS(uuids []string) error {
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

	query := fmt.Sprintf(`DELETE FROM products_tags WHERE uuid IN (%s)`, strings.Join(placeholders, ", "))

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
