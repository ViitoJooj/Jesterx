package helpers

import (
	"database/sql"
	"errors"

	"github.com/ViitoJooj/Jesterx/internal/domain/entities"
)

func ScanTermsSlice(rows *sql.Rows) ([]*domain.Terms, error) {
	var termsList []*domain.Terms

	for rows.Next() {
		terms := &domain.Terms{}
		err := rows.Scan(
			&terms.UUID,
			&terms.Name,
			&terms.Description,
			&terms.CreatedAt,
			&terms.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		termsList = append(termsList, terms)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return termsList, nil
}

func ScanTerms(row *sql.Row) (*domain.Terms, error) {
	terms := &domain.Terms{}

	err := row.Scan(
		&terms.UUID,
		&terms.Name,
		&terms.Description,
		&terms.CreatedAt,
		&terms.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("terms not found")
		}
		return nil, err
	}

	return terms, nil
}
