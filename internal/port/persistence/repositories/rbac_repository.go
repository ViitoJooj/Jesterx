package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ViitoJooj/verkoupe/internal/domain/entities"
	"github.com/ViitoJooj/verkoupe/internal/domain/repositories/contracts"
	"github.com/ViitoJooj/verkoupe/internal/port/persistence/helpers"
)

var _ contracts.RbacContract = (*RbacRepository)(nil)

type RbacRepository struct {
	db *sql.DB
}

func NewRbacRepository(db *sql.DB) *RbacRepository {
	return &RbacRepository{
		db: db,
	}
}

func (r *RbacRepository) CreateRbac(rbac *domain.Rbac) (*domain.Rbac, error) {
	if rbac == nil {
		return nil, errors.New("invalid rbac")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO rbac (website_uuid, label, can_read, can_write, can_update, can_upgrade, can_delete)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING uuid, created_at, updated_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		rbac.WebSiteUUID,
		rbac.Label,
		rbac.CanRead,
		rbac.CanWrite,
		rbac.CanUpdate,
		rbac.CanUpgrade,
		rbac.CanDelete,
	).Scan(
		&rbac.UUID,
		&rbac.CreatedAt,
		&rbac.UpdatedAt,
	)

	if err != nil {
		return nil, errors.New("could not create rbac")
	}

	return rbac, nil
}

func (r *RbacRepository) FindRbacByUUID(uuid string) (*domain.Rbac, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, website_uuid, label, can_read, can_write, can_update, can_upgrade, can_delete, created_at, updated_at
	FROM rbac
	WHERE uuid = $1`

	row := r.db.QueryRowContext(ctx, query, uuid)
	return helpers.ScanRbac(row)
}

func (r *RbacRepository) FindRbacByLabelAndWebsite(label string, websiteUUID string) (*domain.Rbac, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, website_uuid, label, can_read, can_write, can_update, can_upgrade, can_delete, created_at, updated_at
	FROM rbac
	WHERE label = $1 AND website_uuid = $2`

	row := r.db.QueryRowContext(ctx, query, label, websiteUUID)
	return helpers.ScanRbac(row)
}

func (r *RbacRepository) GetRbacFromWebsite(websiteUUID string) ([]*domain.Rbac, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, website_uuid, label, can_read, can_write, can_update, can_upgrade, can_delete, created_at, updated_at
	FROM rbac
	WHERE website_uuid = $1`

	rows, err := r.db.QueryContext(ctx, query, websiteUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return helpers.ScanRbacSlice(rows)
}

func (r *RbacRepository) GetRbac() ([]*domain.Rbac, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, website_uuid, label, can_read, can_write, can_update, can_upgrade, can_delete, created_at, updated_at
	FROM rbac`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return helpers.ScanRbacSlice(rows)
}

func (r *RbacRepository) UpdateRbacByUUID(uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE rbac SET updated_at = NOW() WHERE uuid = $1`

	result, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("rbac not found")
	}

	return nil
}

func (r *RbacRepository) DeleteRbacByUUID(uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM rbac WHERE uuid = $1`

	result, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("rbac not found")
	}

	return nil
}

func (r *RbacRepository) DeleteRbacByUUIDS(uuids []string) error {
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

	query := fmt.Sprintf(`DELETE FROM rbac WHERE uuid IN (%s)`, strings.Join(placeholders, ", "))

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
