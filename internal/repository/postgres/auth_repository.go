package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ViitoJooj/Jesterx/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Save(user domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO users (id, first_name, last_name,  email, password, banned, updated_at, created_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(
		ctx,
		query,
		user.Id,
		user.First_name,
		user.Last_name,
		user.Email,
		user.Password,
		user.Banned,
		user.Updated_at,
		user.Created_at,
	)

	return err
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, email, password, banned, updated_at, created_at
		FROM users
		WHERE email = $1
	`

	var user domain.User

	err := r.db.QueryRowContext(ctx, query, email).
		Scan(&user.Id, &user.Email, &user.Password, &user.Banned, &user.Updated_at, &user.Created_at)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByID(id string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, email, password, banned, updated_at, created_at,
		FROM users
		WHERE id = $1
	`

	var user domain.User

	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&user.Id, &user.Email, &user.Password, &user.Banned, &user.Updated_at, &user.Created_at)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) DeleteByID(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM users WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
