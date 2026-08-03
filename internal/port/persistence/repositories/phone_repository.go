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

var _ contracts.PhoneContract = (*PhoneRepository)(nil)

type PhoneRepository struct {
	db *sql.DB
}

func NewPhoneRepository(db *sql.DB) *PhoneRepository {
	return &PhoneRepository{
		db: db,
	}
}

func (r *PhoneRepository) CreatePhone(phone *domain.Phone) (*domain.Phone, error) {
	if phone == nil {
		return nil, errors.New("invalid phone")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO phones (website_uuid, owner_uuid, owner_type, label, number, is_default)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING uuid, created_at, updated_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		phone.WebSiteUUID,
		phone.OwnerUUID,
		phone.OwnerType,
		phone.Label,
		phone.Number,
		phone.IsDefault,
	).Scan(
		&phone.UUID,
		&phone.CreatedAt,
		&phone.UpdatedAt,
	)

	if err != nil {
		return nil, errors.New("could not create phone")
	}

	return phone, nil
}

func (r *PhoneRepository) FindPhoneByUUID(uuid string) (*domain.Phone, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, website_uuid, owner_uuid, owner_type, label, number, is_default, created_at, updated_at
	FROM phones
	WHERE uuid = $1`

	row := r.db.QueryRowContext(ctx, query, uuid)
	return helpers.ScanPhone(row)
}

func (r *PhoneRepository) FindDefaultPhoneByOwner(ownerUUID string, websiteUUID string) (*domain.Phone, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, website_uuid, owner_uuid, owner_type, label, number, is_default, created_at, updated_at
	FROM phones
	WHERE owner_uuid = $1 AND website_uuid = $2 AND is_default = TRUE
	LIMIT 1`

	row := r.db.QueryRowContext(ctx, query, ownerUUID, websiteUUID)
	return helpers.ScanPhone(row)
}

func (r *PhoneRepository) GetPhonesFromOwner(ownerUUID string, websiteUUID string) ([]*domain.Phone, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, website_uuid, owner_uuid, owner_type, label, number, is_default, created_at, updated_at
	FROM phones
	WHERE owner_uuid = $1 AND website_uuid = $2`

	rows, err := r.db.QueryContext(ctx, query, ownerUUID, websiteUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return helpers.ScanPhones(rows)
}

func (r *PhoneRepository) GetPhonesFromWebsite(websiteUUID string) ([]*domain.Phone, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, website_uuid, owner_uuid, owner_type, label, number, is_default, created_at, updated_at
	FROM phones
	WHERE website_uuid = $1`

	rows, err := r.db.QueryContext(ctx, query, websiteUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return helpers.ScanPhones(rows)
}

func (r *PhoneRepository) GetPhones() ([]*domain.Phone, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, website_uuid, owner_uuid, owner_type, label, number, is_default, created_at, updated_at
	FROM phones`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return helpers.ScanPhones(rows)
}

func (r *PhoneRepository) UpdatePhoneByUUID(uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE phones SET updated_at = NOW() WHERE uuid = $1`

	result, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("phone not found")
	}

	return nil
}

func (r *PhoneRepository) DeletePhoneByUUID(uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM phones WHERE uuid = $1`

	result, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("phone not found")
	}

	return nil
}

func (r *PhoneRepository) DeletePhonesByUUIDS(uuids []string) error {
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

	query := fmt.Sprintf(`DELETE FROM phones WHERE uuid IN (%s)`, strings.Join(placeholders, ", "))

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
