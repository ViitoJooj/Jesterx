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

var _ contracts.OrganizationContract = (*OrganizationRepository)(nil)

type OrganizationRepository struct {
	db *sql.DB
}

func NewOrganizationRepository(db *sql.DB) *OrganizationRepository {
	return &OrganizationRepository{
		db: db,
	}
}

func (r *OrganizationRepository) CreateOrganization(org *domain.OrganizationBR) (*domain.OrganizationBR, error) {
	if org == nil {
		return nil, errors.New("invalid organization")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO organizations (website_uuid, owner_uuid, image_url, name, trade_name, cnpj)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING uuid, created_at, updated_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		org.WebSiteUUID,
		org.OwnerUUID,
		org.ImageURL,
		org.Name,
		org.TradeName,
		org.CNPJ,
	).Scan(
		&org.UUID,
		&org.CreatedAt,
		&org.UpdatedAt,
	)

	if err != nil {
		return nil, errors.New("could not create organization")
	}

	return org, nil
}

func (r *OrganizationRepository) FindOrganizationByUUID(uuid string) (*domain.OrganizationBR, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, website_uuid, owner_uuid, image_url, name, trade_name, cnpj, created_at, updated_at
	FROM organizations
	WHERE uuid = $1`

	row := r.db.QueryRowContext(ctx, query, uuid)
	return helpers.ScanOrganization(row)
}

func (r *OrganizationRepository) FindOrganizationByCNPJAndWebsite(cnpj string, websiteUUID string) (*domain.OrganizationBR, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, website_uuid, owner_uuid, image_url, name, trade_name, cnpj, created_at, updated_at
	FROM organizations
	WHERE cnpj = $1 AND website_uuid = $2`

	row := r.db.QueryRowContext(ctx, query, cnpj, websiteUUID)
	return helpers.ScanOrganization(row)
}

func (r *OrganizationRepository) GetOrganizationsFromWebsite(websiteUUID string) ([]*domain.OrganizationBR, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, website_uuid, owner_uuid, image_url, name, trade_name, cnpj, created_at, updated_at
	FROM organizations
	WHERE website_uuid = $1`

	rows, err := r.db.QueryContext(ctx, query, websiteUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return helpers.ScanOrganizations(rows)
}

func (r *OrganizationRepository) GetOrganizations() ([]*domain.OrganizationBR, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, website_uuid, owner_uuid, image_url, name, trade_name, cnpj, created_at, updated_at
	FROM organizations`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return helpers.ScanOrganizations(rows)
}

func (r *OrganizationRepository) UpdateOrganizationByUUID(uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE organizations SET updated_at = NOW() WHERE uuid = $1`

	result, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("organization not found")
	}

	return nil
}

func (r *OrganizationRepository) DeleteOrganizationByUUID(uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM organizations WHERE uuid = $1`

	result, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("organization not found")
	}

	return nil
}

func (r *OrganizationRepository) DeleteOrganizationsByUUIDS(uuids []string) error {
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

	query := fmt.Sprintf(`DELETE FROM organizations WHERE uuid IN (%s)`, strings.Join(placeholders, ", "))

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
