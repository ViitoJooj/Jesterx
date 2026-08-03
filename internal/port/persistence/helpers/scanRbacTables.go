package helpers

import (
	"database/sql"
	"errors"

	"github.com/ViitoJooj/Jesterx/internal/domain/entities"
)

func ScanRbacSlice(rows *sql.Rows) ([]*domain.Rbac, error) {
	var rbacList []*domain.Rbac

	for rows.Next() {
		rbac := &domain.Rbac{}
		err := rows.Scan(
			&rbac.UUID,
			&rbac.WebSiteUUID,
			&rbac.Label,
			&rbac.CanRead,
			&rbac.CanWrite,
			&rbac.CanUpdate,
			&rbac.CanUpgrade,
			&rbac.CanDelete,
			&rbac.CreatedAt,
			&rbac.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		rbacList = append(rbacList, rbac)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rbacList, nil
}

func ScanRbac(row *sql.Row) (*domain.Rbac, error) {
	rbac := &domain.Rbac{}

	err := row.Scan(
		&rbac.UUID,
		&rbac.WebSiteUUID,
		&rbac.Label,
		&rbac.CanRead,
		&rbac.CanWrite,
		&rbac.CanUpdate,
		&rbac.CanUpgrade,
		&rbac.CanDelete,
		&rbac.CreatedAt,
		&rbac.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("rbac not found")
		}
		return nil, err
	}

	return rbac, nil
}
