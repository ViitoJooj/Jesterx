package helpers

import (
	"database/sql"
	"errors"

	"github.com/ViitoJooj/verkoupe/internal/domain/entities"
)

func ScanAddresses(rows *sql.Rows) ([]*domain.AddressBR, error) {
	var addresses []*domain.AddressBR

	for rows.Next() {
		addr := &domain.AddressBR{}
		err := rows.Scan(
			&addr.UUID,
			&addr.WebSiteUUID,
			&addr.OwnerUUID,
			&addr.OwnerType,
			&addr.Label,
			&addr.AddressLine1,
			&addr.AddressLine2,
			&addr.Neighborhood,
			&addr.City,
			&addr.State,
			&addr.StateCode,
			&addr.PostalCode,
			&addr.ReferencePoint,
			&addr.DeliveryNotes,
			&addr.IsDefault,
			&addr.CreatedAt,
			&addr.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, addr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return addresses, nil
}

func ScanAddress(row *sql.Row) (*domain.AddressBR, error) {
	addr := &domain.AddressBR{}

	err := row.Scan(
		&addr.UUID,
		&addr.WebSiteUUID,
		&addr.OwnerUUID,
		&addr.OwnerType,
		&addr.Label,
		&addr.AddressLine1,
		&addr.AddressLine2,
		&addr.Neighborhood,
		&addr.City,
		&addr.State,
		&addr.StateCode,
		&addr.PostalCode,
		&addr.ReferencePoint,
		&addr.DeliveryNotes,
		&addr.IsDefault,
		&addr.CreatedAt,
		&addr.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("address not found")
		}
		return nil, err
	}

	return addr, nil
}
