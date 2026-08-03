package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ViitoJooj/Jesterx/internal/domain/entities"
	"github.com/ViitoJooj/Jesterx/internal/domain/repositories/contracts"
	"github.com/ViitoJooj/Jesterx/internal/port/persistence/helpers"
)

var _ contracts.WebsiteContract = (*WebsiteRepository)(nil)

type WebsiteRepository struct {
	db *sql.DB
}

func NewWebsiteRepository(db *sql.DB) *WebsiteRepository {
	return &WebsiteRepository{
		db: db,
	}
}

func (r *WebsiteRepository) CreateWebsite(website *domain.Website) (*domain.Website, error) {
	if website == nil {
		return nil, errors.New("invalid website")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO websites (owner_uuid, owner_type, label, url, write_in, description)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING uuid, created_at, updated_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		website.OwnerUUID,
		website.OwnerType,
		website.Label,
		website.URL,
		website.WriteIn,
		website.Description,
	).Scan(
		&website.UUID,
		&website.CreatedAt,
		&website.UpdatedAt,
	)

	if err != nil {
		return nil, errors.New("could not create website")
	}

	return website, nil
}

func (r *WebsiteRepository) FindWebsiteByUUID(uuid string) (*domain.Website, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, owner_uuid, owner_type, label, url, write_in, description, updated_at, created_at
	FROM websites
	WHERE uuid = $1`

	row := r.db.QueryRowContext(ctx, query, uuid)
	return helpers.ScanWebsite(row)
}

func (r *WebsiteRepository) FindWebsiteByLabel(label string) (*domain.Website, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, owner_uuid, owner_type, label, url, write_in, description, updated_at, created_at
	FROM websites
	WHERE label = $1`

	row := r.db.QueryRowContext(ctx, query, label)
	return helpers.ScanWebsite(row)
}

func (r *WebsiteRepository) FindWebsitesByOwner(ownerUUID string) ([]*domain.Website, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, owner_uuid, owner_type, label, url, write_in, description, updated_at, created_at
	FROM websites
	WHERE owner_uuid = $1`

	rows, err := r.db.QueryContext(ctx, query, ownerUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return helpers.ScanWebsites(rows)
}

func (r *WebsiteRepository) GetWebsites() ([]*domain.Website, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, owner_uuid, owner_type, label, url, write_in, description, updated_at, created_at
	FROM websites`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return helpers.ScanWebsites(rows)
}

func (r *WebsiteRepository) UpdateWebsiteByUUID(uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE websites SET updated_at = NOW() WHERE uuid = $1`

	result, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("website not found")
	}

	return nil
}

func (r *WebsiteRepository) DeleteWebsiteByUUID(uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM websites WHERE uuid = $1`

	result, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("website not found")
	}

	return nil
}

func (r *WebsiteRepository) DeleteWebsitesByUUIDS(uuids []string) error {
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

	query := fmt.Sprintf(`DELETE FROM websites WHERE uuid IN (%s)`, strings.Join(placeholders, ", "))

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
