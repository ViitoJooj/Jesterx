package helpers

import (
	"database/sql"
	"errors"

	"github.com/ViitoJooj/Jesterx/internal/domain/entities"
)

func ScanOrganizations(rows *sql.Rows) ([]*domain.OrganizationBR, error) {
	var orgs []*domain.OrganizationBR

	for rows.Next() {
		org := &domain.OrganizationBR{}
		err := rows.Scan(
			&org.UUID,
			&org.WebSiteUUID,
			&org.OwnerUUID,
			&org.ImageURL,
			&org.Name,
			&org.TradeName,
			&org.CNPJ,
			&org.CreatedAt,
			&org.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		orgs = append(orgs, org)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orgs, nil
}

func ScanOrganization(row *sql.Row) (*domain.OrganizationBR, error) {
	org := &domain.OrganizationBR{}

	err := row.Scan(
		&org.UUID,
		&org.WebSiteUUID,
		&org.OwnerUUID,
		&org.ImageURL,
		&org.Name,
		&org.TradeName,
		&org.CNPJ,
		&org.CreatedAt,
		&org.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("organization not found")
		}
		return nil, err
	}

	return org, nil
}
