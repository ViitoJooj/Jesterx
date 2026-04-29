package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ViitoJooj/Jesterx/internal/domain"
)

func NewUserRepository(db *sql.DB) UsersRepository {
	return &userRepository{db: db}
}

type UsersRepository interface {
	InsertUser(*domain.User) (*domain.User, error)
	FindUserById(id string) (*domain.User, error)
	FindUserByEmail(email string) (*domain.User, error)
	UpdateUser(*domain.User) error
	DeleteUserById(id string) error
}

type userRepository struct {
	db *sql.DB
}

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailAlreadyInUse = errors.New("email already in use")
)

func (r *userRepository) InsertUser(user *domain.User) (*domain.User, error) {
	ctx := context.Background()
	exists, err := r.UserExistsByEmail(user.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailAlreadyInUse
	}

	query := `
		INSERT INTO users (name, email, password, role, cpf)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, email, password, role, cpf, created_at, updated_at
	`

	var newUser domain.User

	err = r.db.QueryRowContext(ctx, query,
		user.Name,
		user.Email,
		user.Password,
		user.Role,
		user.Cpf,
	).Scan(
		&newUser.Uuid,
		&newUser.Name,
		&newUser.Email,
		&newUser.Password,
		&newUser.Role,
		&newUser.Cpf,
		&newUser.Created_at,
		&newUser.Updated_at,
	)

	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (r *userRepository) FindUserById(id string) (*domain.User, error) {
	ctx := context.Background()

	query := `
		SELECT id, name, email, password, role, cpf, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user domain.User

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.Uuid,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.Cpf,
		&user.Created_at,
		&user.Updated_at,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindUserByEmail(email string) (*domain.User, error) {
	ctx := context.Background()

	query := `
		SELECT id, name, email, password, role, cpf, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user domain.User

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.Uuid,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.Cpf,
		&user.Created_at,
		&user.Updated_at,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) UpdateUser(user *domain.User) error {
	ctx := context.Background()

	exists, err := r.UserExists(user.Uuid)
	if err != nil {
		return err
	}
	if !exists {
		return ErrUserNotFound
	}

	query := `
		UPDATE users
		SET name = $1,
			email = $2,
			password = $3,
			role = $4,
			cpf = $5,
			updated_at = NOW()
		WHERE id = $6
	`

	result, err := r.db.ExecContext(ctx, query,
		user.Name,
		user.Email,
		user.Password,
		user.Role,
		user.Cpf,
		user.Uuid,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *userRepository) DeleteUserById(id string) error {
	ctx := context.Background()

	exists, err := r.UserExists(id)
	if err != nil {
		return err
	}
	if !exists {
		return ErrUserNotFound
	}

	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *userRepository) UserExists(id string) (bool, error) {
	ctx := context.Background()

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *userRepository) UserExistsByEmail(email string) (bool, error) {
	ctx := context.Background()

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
