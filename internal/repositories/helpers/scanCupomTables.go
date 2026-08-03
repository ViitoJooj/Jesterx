package helpers

import (
	"database/sql"
	"errors"

	"github.com/ViitoJooj/Jesterx/internal/domain"
)

func ScanCupoms(rows *sql.Rows) ([]*domain.Cupons, error) {
	var cupoms []*domain.Cupons

	for rows.Next() {
		c := &domain.Cupons{}
		err := rows.Scan(
			&c.UUID,
			&c.TagUUID,
			&c.Label,
			&c.Description,
			&c.Value,
			&c.ValueType,
		)
		if err != nil {
			return nil, err
		}
		cupoms = append(cupoms, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cupoms, nil
}

func ScanCupom(row *sql.Row) (*domain.Cupons, error) {
	c := &domain.Cupons{}

	err := row.Scan(
		&c.UUID,
		&c.TagUUID,
		&c.Label,
		&c.Description,
		&c.Value,
		&c.ValueType,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("cupom not found")
		}
		return nil, err
	}

	return c, nil
}
