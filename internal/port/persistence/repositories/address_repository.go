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

var _ contracts.AddressContract = (*AddressRepository)(nil)

type AddressRepository struct {
	db *sql.DB
}

func NewAddressRepository(db *sql.DB) *AddressRepository {
	return &AddressRepository{
		db: db,
	}
}

func (r *AddressRepository) CreateAddress(address *domain.AddressBR) (*domain.AddressBR, error) {
	if address == nil {
		return nil, errors.New("invalid address")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO addresses (website_uuid, owner_uuid, owner_type, label, address_line1, address_line2, neighborhood, city, state, state_code, postal_code, reference_point, delivery_notes, is_default)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	RETURNING uuid, created_at, updated_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		address.WebSiteUUID,
		address.OwnerUUID,
		address.OwnerType,
		address.Label,
		address.AddressLine1,
		address.AddressLine2,
		address.Neighborhood,
		address.City,
		address.State,
		address.StateCode,
		address.PostalCode,
		address.ReferencePoint,
		address.DeliveryNotes,
		address.IsDefault,
	).Scan(
		&address.UUID,
		&address.CreatedAt,
		&address.UpdatedAt,
	)

	if err != nil {
		return nil, errors.New("could not create address")
	}

	return address, nil
}

func (r *AddressRepository) FindAddressByUUID(uuid string) (*domain.AddressBR, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, website_uuid, owner_uuid, owner_type, label, address_line1, address_line2, neighborhood, city, state, state_code, postal_code, reference_point, delivery_notes, is_default, created_at, updated_at
	FROM addresses
	WHERE uuid = $1`

	row := r.db.QueryRowContext(ctx, query, uuid)
	return helpers.ScanAddress(row)
}

func (r *AddressRepository) FindDefaultAddressByOwner(ownerUUID string, websiteUUID string) (*domain.AddressBR, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, website_uuid, owner_uuid, owner_type, label, address_line1, address_line2, neighborhood, city, state, state_code, postal_code, reference_point, delivery_notes, is_default, created_at, updated_at
	FROM addresses
	WHERE owner_uuid = $1 AND website_uuid = $2 AND is_default = TRUE
	LIMIT 1`

	row := r.db.QueryRowContext(ctx, query, ownerUUID, websiteUUID)
	return helpers.ScanAddress(row)
}

func (r *AddressRepository) GetAddressesFromOwner(ownerUUID string, websiteUUID string) ([]*domain.AddressBR, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, website_uuid, owner_uuid, owner_type, label, address_line1, address_line2, neighborhood, city, state, state_code, postal_code, reference_point, delivery_notes, is_default, created_at, updated_at
	FROM addresses
	WHERE owner_uuid = $1 AND website_uuid = $2`

	rows, err := r.db.QueryContext(ctx, query, ownerUUID, websiteUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return helpers.ScanAddresses(rows)
}

func (r *AddressRepository) GetAddressesFromWebsite(websiteUUID string) ([]*domain.AddressBR, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, website_uuid, owner_uuid, owner_type, label, address_line1, address_line2, neighborhood, city, state, state_code, postal_code, reference_point, delivery_notes, is_default, created_at, updated_at
	FROM addresses
	WHERE website_uuid = $1`

	rows, err := r.db.QueryContext(ctx, query, websiteUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return helpers.ScanAddresses(rows)
}

func (r *AddressRepository) GetAddresses() ([]*domain.AddressBR, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, website_uuid, owner_uuid, owner_type, label, address_line1, address_line2, neighborhood, city, state, state_code, postal_code, reference_point, delivery_notes, is_default, created_at, updated_at
	FROM addresses`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return helpers.ScanAddresses(rows)
}

func (r *AddressRepository) UpdateAddressByUUID(uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE addresses SET updated_at = NOW() WHERE uuid = $1`

	result, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("address not found")
	}

	return nil
}

func (r *AddressRepository) DeleteAddressByUUID(uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM addresses WHERE uuid = $1`

	result, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("address not found")
	}

	return nil
}

func (r *AddressRepository) DeleteAddressesByUUIDS(uuids []string) error {
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

	query := fmt.Sprintf(`DELETE FROM addresses WHERE uuid IN (%s)`, strings.Join(placeholders, ", "))

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
