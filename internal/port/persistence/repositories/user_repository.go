package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ViitoJooj/verkoupe/internal/domain/entities"
	"github.com/ViitoJooj/verkoupe/internal/port/persistence/helpers"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(user *domain.User) (*domain.User, error) {
	if user == nil {
		return nil, errors.New("invalid user")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO users (website_uuid, image_url, name, email, role, password, cpf)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING uuid, created_at, updated_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		user.WebSiteUUID,
		user.ImageURL,
		user.Name,
		user.Email,
		user.Role,
		user.Password,
		user.CPF,
	).Scan(
		&user.UUID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, errors.New("could not create user")
	}

	return user, nil
}

func (r *UserRepository) FindUserByUUID(uuid string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, website_uuid, image_url, name, email, role, password, cpf, created_at, updated_at
	FROM users
	WHERE uuid = $1`

	row := r.db.QueryRowContext(ctx, query, uuid)
	return helpers.ScanUser(row)
}

func (r *UserRepository) FindUserByEmailAndWebsite(email string, websiteUUID string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, website_uuid, image_url, name, email, role, password, cpf, created_at, updated_at
	FROM users
	WHERE email = $1 AND website_uuid = $2`

	row := r.db.QueryRowContext(ctx, query, email, websiteUUID)
	return helpers.ScanUser(row)
}

func (r *UserRepository) UserExists(email string, websiteUUID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND website_uuid = $2)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email, websiteUUID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *UserRepository) GetUsersFromWebsite(websiteUUID string) ([]*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, website_uuid, image_url, name, email, role, password, cpf, created_at, updated_at
	FROM users
	WHERE website_uuid = $1`

	rows, err := r.db.QueryContext(ctx, query, websiteUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return helpers.ScanUsers(rows)
}

func (r *UserRepository) GetUsers() ([]*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, website_uuid, image_url, name, email, role, password, cpf, created_at, updated_at
	FROM users`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return helpers.ScanUsers(rows)
}

func (r *UserRepository) UpdateUserByUUID(uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE users SET updated_at = NOW() WHERE uuid = $1`

	result, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *UserRepository) DeleteUserByUUID(uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM users WHERE uuid = $1`

	result, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *UserRepository) DeleteUsersByUUIDS(uuids []string) error {
	if len(uuids) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	placeholders := make([]string, len(uuids))
	args := make([]interface{}, len(uuids))

	for i, u := range uuids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = u
	}

	query := fmt.Sprintf(`DELETE FROM users WHERE uuid IN (%s)`, strings.Join(placeholders, ", "))

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
