package helpers

import (
	"database/sql"
	"errors"

	"github.com/ViitoJooj/Jesterx/internal/domain/entities"
)

func ScanProductTags(rows *sql.Rows) ([]*domain.ProductsTags, error) {
	var tags []*domain.ProductsTags

	for rows.Next() {
		t := &domain.ProductsTags{}
		err := rows.Scan(
			&t.UUID,
			&t.ProductUUID,
			&t.Label,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func ScanProductTag(row *sql.Row) (*domain.ProductsTags, error) {
	t := &domain.ProductsTags{}

	err := row.Scan(
		&t.UUID,
		&t.ProductUUID,
		&t.Label,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("product tag not found")
		}
		return nil, err
	}

	return t, nil
}
