package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ViitoJooj/Jesterx/internal/domain"
)

func NewWebSiteRepository(db *sql.DB) *connection {
	return &connection{db: db}
}

func (r *connection) SaveWebSite(website domain.WebSite) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `INSERT INTO websites (id, website_type, image, name, short_description, description, creator_id, banned, updated_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.ExecContext(ctx, query, website.Id, website.Type, website.Image, website.Name, website.Short_description, website.Description, website.Creator_id, website.Banned, website.Updated_at, website.Created_at)
	return err
}

func (r *connection) FindWebSiteByID(id string) (*domain.WebSite, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `SELECT * FROM websites WHERE id = $1`
	var website domain.WebSite
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&website.Id, &website.Type, &website.Image, &website.Name, &website.Short_description, &website.Description, &website.Creator_id, &website.Banned, &website.Updated_at, &website.Created_at)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &website, nil
}

func (r *connection) FindWebSiteByUserID(id string) (*domain.WebSite, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `SELECT * FROM websites WHERE creator_id = $1`
	var website domain.WebSite
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&website.Id, &website.Type, &website.Image, &website.Name, &website.Short_description, &website.Description, &website.Creator_id, &website.Banned, &website.Updated_at, &website.Created_at)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &website, nil
}

func (r *connection) FindWebSiteByName(name string) (*domain.WebSite, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT * FROM websites WHERE name = $1`

	var website domain.WebSite
	err := r.db.QueryRowContext(ctx, query, name).
		Scan(&website.Id, &website.Type, &website.Image, &website.Name, &website.Short_description, &website.Description, &website.Creator_id, &website.Banned, &website.Updated_at, &website.Created_at)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &website, nil
}

func (r *connection) UpdateWebSiteByID(id string, website domain.WebSite) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `UPDATE websites SET website_type = $1, image = $2, name = $3, short_description = $4, description = $5, creator_id = $6, banned = $7, updated_at = $8 WHERE id = $9`
	_, err := r.db.ExecContext(ctx, query, website.Type, website.Image, website.Name, website.Short_description, website.Description, website.Creator_id, website.Banned, website.Updated_at, id)
	return err
}

func (r *connection) DeleteWebSiteByID(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `DELETE FROM websites WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
