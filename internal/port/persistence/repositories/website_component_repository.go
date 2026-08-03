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

var _ contracts.WebsiteComponentContract = (*WebsiteComponentRepository)(nil)

type WebsiteComponentRepository struct {
	db *sql.DB
}

func NewWebsiteComponentRepository(db *sql.DB) *WebsiteComponentRepository {
	return &WebsiteComponentRepository{
		db: db,
	}
}

func (r *WebsiteComponentRepository) CreateWebsiteComponent(component *domain.ComponentWebsites) (*domain.ComponentWebsites, error) {
	if component == nil {
		return nil, errors.New("invalid website component")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO websites_components (website_uuid, logo_url, tittle, description, path, content, visits)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING uuid, created_at, updated_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		component.WebsiteUUID,
		component.LogoURL,
		component.Tittle,
		component.Description,
		component.Path,
		component.Content,
		component.Visists,
	).Scan(
		&component.UUID,
		&component.CreatedAt,
		&component.UpdatedAt,
	)

	if err != nil {
		return nil, errors.New("could not create website component")
	}

	return component, nil
}

func (r *WebsiteComponentRepository) FindWebsiteComponentByUUID(uuid string) (*domain.ComponentWebsites, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, website_uuid, logo_url, tittle, description, path, content, visits, updated_by, updated_at, created_at
	FROM websites_components
	WHERE uuid = $1`

	row := r.db.QueryRowContext(ctx, query, uuid)
	return helpers.ScanWebsiteComponent(row)
}

func (r *WebsiteComponentRepository) FindWebsiteComponentByPath(path string) (*domain.ComponentWebsites, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, website_uuid, logo_url, tittle, description, path, content, visits, updated_by, updated_at, created_at
	FROM websites_components
	WHERE path = $1`

	row := r.db.QueryRowContext(ctx, query, path)
	return helpers.ScanWebsiteComponent(row)
}

func (r *WebsiteComponentRepository) GetWebsiteComponentsFromWebsite(websiteUUID string) ([]*domain.ComponentWebsites, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, website_uuid, logo_url, tittle, description, path, content, visits, updated_by, updated_at, created_at
	FROM websites_components
	WHERE website_uuid = $1`

	rows, err := r.db.QueryContext(ctx, query, websiteUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return helpers.ScanWebsiteComponents(rows)
}

func (r *WebsiteComponentRepository) GetWebsiteComponents() ([]*domain.ComponentWebsites, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, website_uuid, logo_url, tittle, description, path, content, visits, updated_by, updated_at, created_at
	FROM websites_components`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return helpers.ScanWebsiteComponents(rows)
}

func (r *WebsiteComponentRepository) UpdateWebsiteComponentByUUID(uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE websites_components SET updated_at = NOW() WHERE uuid = $1`

	result, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("website component not found")
	}

	return nil
}

func (r *WebsiteComponentRepository) DeleteWebsiteComponentByUUID(uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM websites_components WHERE uuid = $1`

	result, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("website component not found")
	}

	return nil
}

func (r *WebsiteComponentRepository) DeleteWebsiteComponentsByUUIDS(uuids []string) error {
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

	query := fmt.Sprintf(`DELETE FROM websites_components WHERE uuid IN (%s)`, strings.Join(placeholders, ", "))

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
