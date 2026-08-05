package helpers

import (
	"database/sql"
	"errors"

	"github.com/ViitoJooj/verkoupe/internal/domain/entities"
)

func ScanWebsiteComponents(rows *sql.Rows) ([]*domain.ComponentWebsites, error) {
	var components []*domain.ComponentWebsites

	for rows.Next() {
		c := &domain.ComponentWebsites{}
		err := rows.Scan(
			&c.UUID,
			&c.WebsiteUUID,
			&c.LogoURL,
			&c.Tittle,
			&c.Description,
			&c.Path,
			&c.Content,
			&c.Visists,
			&c.UpdatedBy,
			&c.UpdatedAt,
			&c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		components = append(components, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return components, nil
}

func ScanWebsiteComponent(row *sql.Row) (*domain.ComponentWebsites, error) {
	c := &domain.ComponentWebsites{}

	err := row.Scan(
		&c.UUID,
		&c.WebsiteUUID,
		&c.LogoURL,
		&c.Tittle,
		&c.Description,
		&c.Path,
		&c.Content,
		&c.Visists,
		&c.UpdatedBy,
		&c.UpdatedAt,
		&c.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("website component not found")
		}
		return nil, err
	}

	return c, nil
}
