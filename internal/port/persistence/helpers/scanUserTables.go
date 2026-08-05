package helpers

import (
	"database/sql"
	"errors"

	"github.com/ViitoJooj/verkoupe/internal/domain/entities"
)

func ScanUsers(rows *sql.Rows) ([]*domain.User, error) {
	var users []*domain.User

	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(
			&user.UUID,
			&user.WebSiteUUID,
			&user.ImageURL,
			&user.Name,
			&user.Email,
			&user.Role,
			&user.Password,
			&user.CPF,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func ScanUser(row *sql.Row) (*domain.User, error) {
	user := &domain.User{}

	err := row.Scan(
		&user.UUID,
		&user.WebSiteUUID,
		&user.ImageURL,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.Password,
		&user.CPF,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}
