package helpers

import (
	"database/sql"
	"errors"

	"github.com/ViitoJooj/verkoupe/internal/domain/entities"
)

func ScanPhones(rows *sql.Rows) ([]*domain.Phone, error) {
	var phones []*domain.Phone

	for rows.Next() {
		phone := &domain.Phone{}
		err := rows.Scan(
			&phone.UUID,
			&phone.WebSiteUUID,
			&phone.OwnerUUID,
			&phone.OwnerType,
			&phone.Label,
			&phone.Number,
			&phone.IsDefault,
			&phone.CreatedAt,
			&phone.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		phones = append(phones, phone)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return phones, nil
}

func ScanPhone(row *sql.Row) (*domain.Phone, error) {
	phone := &domain.Phone{}

	err := row.Scan(
		&phone.UUID,
		&phone.WebSiteUUID,
		&phone.OwnerUUID,
		&phone.OwnerType,
		&phone.Label,
		&phone.Number,
		&phone.IsDefault,
		&phone.CreatedAt,
		&phone.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("phone not found")
		}
		return nil, err
	}

	return phone, nil
}
