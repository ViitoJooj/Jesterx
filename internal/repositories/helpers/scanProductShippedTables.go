package helpers

import (
	"database/sql"
	"errors"

	"github.com/ViitoJooj/Jesterx/internal/domain"
)

func ScanProductsShipped(rows *sql.Rows) ([]*domain.ProductShipped, error) {
	var productsShipped []*domain.ProductShipped

	for rows.Next() {
		ps := &domain.ProductShipped{}
		err := rows.Scan(
			&ps.UUID,
			&ps.ProductUUID,
			&ps.AddressUUID,
			&ps.Status,
		)
		if err != nil {
			return nil, err
		}
		productsShipped = append(productsShipped, ps)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return productsShipped, nil
}

func ScanProductShipped(row *sql.Row) (*domain.ProductShipped, error) {
	ps := &domain.ProductShipped{}

	err := row.Scan(
		&ps.UUID,
		&ps.ProductUUID,
		&ps.AddressUUID,
		&ps.Status,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("product shipped not found")
		}
		return nil, err
	}

	return ps, nil
}
