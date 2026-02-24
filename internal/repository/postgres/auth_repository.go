package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ViitoJooj/Jesterx/internal/domain"
)

func NewAuthRepository(db *sql.DB) *connection {
	return &connection{db: db}
}

func (r *connection) UserRegister(user domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `INSERT INTO users (id, website_id, first_name, last_name, email, password, role, updated_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.ExecContext(ctx, query, user.Id, user.WebsiteId, user.First_name, user.Last_name, user.Email, user.Password, user.Role, user.Updated_at, user.Created_at)
	return err
}

func (r *connection) FindUserByEmail(email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `SELECT id, website_id, first_name, last_name, email, password, role, updated_at, created_at FROM users WHERE email = $1`
	var user domain.User
	err := r.db.QueryRowContext(ctx, query, email).
		Scan(&user.Id, &user.WebsiteId, &user.First_name, &user.Last_name, &user.Email, &user.Password, &user.Role, &user.Updated_at, &user.Created_at)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *connection) FindUserByID(id string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `SELECT id, website_id, first_name, last_name, email, password, role, updated_at, created_at FROM users WHERE id = $1`
	var user domain.User
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&user.Id, &user.WebsiteId, &user.First_name, &user.Last_name, &user.Email, &user.Password, &user.Role, &user.Updated_at, &user.Created_at)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *connection) DeleteUserByID(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *connection) FindUserByEmailAndWebsite(email string, websiteId string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, website_id, first_name, last_name, email, password, role, updated_at, created_at
		FROM users
		WHERE email = $1 AND website_id = $2
	`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, email, websiteId).
		Scan(
			&user.Id,
			&user.WebsiteId,
			&user.First_name,
			&user.Last_name,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.Updated_at,
			&user.Created_at,
		)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}
