package helpers

import (
	"database/sql"
	"errors"

	"github.com/ViitoJooj/Jesterx/internal/domain"
)

func ScanStorageProducts(rows *sql.Rows) ([]*domain.StorageProducts, error) {
	var sps []*domain.StorageProducts

	for rows.Next() {
		sp := &domain.StorageProducts{}
		err := rows.Scan(
			&sp.UUID,
			&sp.ProductUUID,
		)
		if err != nil {
			return nil, err
		}
		sps = append(sps, sp)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sps, nil
}

func ScanStorageProduct(row *sql.Row) (*domain.StorageProducts, error) {
	sp := &domain.StorageProducts{}

	err := row.Scan(
		&sp.UUID,
		&sp.ProductUUID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("storage product not found")
		}
		return nil, err
	}

	return sp, nil
}
