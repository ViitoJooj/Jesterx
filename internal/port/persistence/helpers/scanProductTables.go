package helpers

import (
	"database/sql"
	"errors"

	"github.com/ViitoJooj/verkoupe/internal/domain/entities"
)

func ScanProducts(rows *sql.Rows) ([]*domain.Products, error) {
	var products []*domain.Products

	for rows.Next() {
		p := &domain.Products{}
		err := rows.Scan(
			&p.UUID,
			&p.Name,
			&p.Description,
			&p.ShortDescription,
			&p.Active,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func ScanProduct(row *sql.Row) (*domain.Products, error) {
	p := &domain.Products{}

	err := row.Scan(
		&p.UUID,
		&p.Name,
		&p.Description,
		&p.ShortDescription,
		&p.Active,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return p, nil
}
