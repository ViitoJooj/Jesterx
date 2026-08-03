package helpers

import (
	"database/sql"
	"errors"

	"github.com/ViitoJooj/Jesterx/internal/domain"
)

func ScanPreparingShippingProducts(rows *sql.Rows) ([]*domain.PreparingShippingProducts, error) {
	var psps []*domain.PreparingShippingProducts

	for rows.Next() {
		psp := &domain.PreparingShippingProducts{}
		err := rows.Scan(
			&psp.UUID,
			&psp.ProductUUID,
			&psp.AddressUUID,
		)
		if err != nil {
			return nil, err
		}
		psps = append(psps, psp)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return psps, nil
}

func ScanPreparingShippingProduct(row *sql.Row) (*domain.PreparingShippingProducts, error) {
	psp := &domain.PreparingShippingProducts{}

	err := row.Scan(
		&psp.UUID,
		&psp.ProductUUID,
		&psp.AddressUUID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("preparing shipping product not found")
		}
		return nil, err
	}

	return psp, nil
}
