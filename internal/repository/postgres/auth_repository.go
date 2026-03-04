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
	query := `INSERT INTO users (id, website_id, first_name, last_name, email, verified_email, password, role, updated_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.ExecContext(ctx, query, user.Id, user.WebsiteId, user.First_name, user.Last_name, user.Email, user.Verified_email, user.Password, user.Role, user.Updated_at, user.Created_at)
	return err
}

func (r *connection) FindUserByID(id string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT 
		u.id,
		u.website_id,
		u.first_name,
		u.last_name,
		u.email,
		u.verified_email,
		u.password,
		u.role,
		u.updated_at,
		u.created_at,
		p.name AS plan_name
	FROM users u
	LEFT JOIN LATERAL (
		SELECT *
		FROM payments pay
		WHERE pay.user_id = u.id
		  AND pay.website_id = u.website_id
		  AND pay.status = 'completed'
		ORDER BY pay.purchased_in DESC
		LIMIT 1
	) pay ON TRUE
	LEFT JOIN plans p
		ON p.price = pay.amount
	   AND p.active = true
	WHERE u.id = $1
	`

	var user domain.User
	var planName sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).
		Scan(
			&user.Id,
			&user.WebsiteId,
			&user.First_name,
			&user.Last_name,
			&user.Email,
			&user.Verified_email,
			&user.Password,
			&user.Role,
			&user.Updated_at,
			&user.Created_at,
			&planName,
		)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if planName.Valid {
		user.Plan = &planName.String
	}

	return &user, nil
}

func (r *connection) FindUserByEmail(email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT 
		u.id,
		u.website_id,
		u.first_name,
		u.last_name,
		u.email,
		u.verified_email,
		u.password,
		u.role,
		u.updated_at,
		u.created_at,
		p.name AS plan_name
	FROM users u
	LEFT JOIN LATERAL (
		SELECT *
		FROM payments pay
		WHERE pay.user_id = u.id
		  AND pay.website_id = u.website_id
		  AND pay.status = 'completed'
		ORDER BY pay.purchased_in DESC
		LIMIT 1
	) pay ON TRUE
	LEFT JOIN plans p
		ON p.price = pay.amount
	   AND p.active = true
	WHERE u.email = $1
	`

	var user domain.User
	var planName sql.NullString

	err := r.db.QueryRowContext(ctx, query, email).
		Scan(
			&user.Id,
			&user.WebsiteId,
			&user.First_name,
			&user.Last_name,
			&user.Email,
			&user.Verified_email,
			&user.Password,
			&user.Role,
			&user.Updated_at,
			&user.Created_at,
			&planName,
		)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if planName.Valid {
		user.Plan = &planName.String
	}

	return &user, nil
}

func (r *connection) FindUserByEmailAndWebsite(email string, websiteId string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT 
		u.id,
		u.website_id,
		u.first_name,
		u.last_name,
		u.email,
		u.verified_email,
		u.password,
		u.role,
		u.updated_at,
		u.created_at,
		p.name AS plan_name
	FROM users u
	LEFT JOIN LATERAL (
		SELECT *
		FROM payments pay
		WHERE pay.user_id = u.id
		  AND pay.website_id = u.website_id
		  AND pay.status = 'completed'
		ORDER BY pay.purchased_in DESC
		LIMIT 1
	) pay ON TRUE
	LEFT JOIN plans p
		ON p.price = pay.amount
	   AND p.active = true
	WHERE u.email = $1
	  AND u.website_id = $2
	`

	var user domain.User
	var planName sql.NullString

	err := r.db.QueryRowContext(ctx, query, email, websiteId).
		Scan(
			&user.Id,
			&user.WebsiteId,
			&user.First_name,
			&user.Last_name,
			&user.Email,
			&user.Verified_email,
			&user.Password,
			&user.Role,
			&user.Updated_at,
			&user.Created_at,
			&planName,
		)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if planName.Valid {
		user.Plan = &planName.String
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

func (r *connection) DeleteExpiredUnverifiedUsers() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, `
		DELETE FROM users
		WHERE verified_email = false
		AND created_at < NOW() - INTERVAL '10 minutes'
	`)
	return err
}

func (r *connection) UpdateVerifiedEmailToTrue(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `UPDATE users SET verified_email = true WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
