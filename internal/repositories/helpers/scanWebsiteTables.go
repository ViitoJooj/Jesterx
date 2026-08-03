package helpers

import (
	"database/sql"
	"errors"

	"github.com/ViitoJooj/Jesterx/internal/domain"
)

func ScanWebsites(rows *sql.Rows) ([]*domain.Website, error) {
	var websites []*domain.Website

	for rows.Next() {
		w := &domain.Website{}
		err := rows.Scan(
			&w.UUID,
			&w.OwnerUUID,
			&w.OwnerType,
			&w.Label,
			&w.URL,
			&w.WriteIn,
			&w.Description,
			&w.UpdatedAt,
			&w.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		websites = append(websites, w)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return websites, nil
}

func ScanWebsite(row *sql.Row) (*domain.Website, error) {
	w := &domain.Website{}

	err := row.Scan(
		&w.UUID,
		&w.OwnerUUID,
		&w.OwnerType,
		&w.Label,
		&w.URL,
		&w.WriteIn,
		&w.Description,
		&w.UpdatedAt,
		&w.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("website not found")
		}
		return nil, err
	}

	return w, nil
}
