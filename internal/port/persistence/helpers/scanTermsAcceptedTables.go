package helpers

import (
	"database/sql"
	"errors"

	"github.com/ViitoJooj/Jesterx/internal/domain/entities"
)

func ScanTermsAcceptedSlice(rows *sql.Rows) ([]*domain.TermsAcceptedBy, error) {
	var list []*domain.TermsAcceptedBy

	for rows.Next() {
		t := &domain.TermsAcceptedBy{}
		err := rows.Scan(
			&t.UUID,
			&t.WebSiteUUID,
			&t.OwnerUUID,
			&t.OwnerType,
			&t.AcceptedWhen,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func ScanTermsAccepted(row *sql.Row) (*domain.TermsAcceptedBy, error) {
	t := &domain.TermsAcceptedBy{}

	err := row.Scan(
		&t.UUID,
		&t.WebSiteUUID,
		&t.OwnerUUID,
		&t.OwnerType,
		&t.AcceptedWhen,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("terms accepted not found")
		}
		return nil, err
	}

	return t, nil
}
