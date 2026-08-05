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

var _ contracts.TermsAcceptedContract = (*TermsAcceptedRepository)(nil)

type TermsAcceptedRepository struct {
	db *sql.DB
}

func NewTermsAcceptedRepository(db *sql.DB) *TermsAcceptedRepository {
	return &TermsAcceptedRepository{
		db: db,
	}
}

func (r *TermsAcceptedRepository) CreateTermsAccepted(termsAccepted *domain.TermsAcceptedBy) (*domain.TermsAcceptedBy, error) {
	if termsAccepted == nil {
		return nil, errors.New("invalid terms accepted")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO terms_accepted (website_uuid, owner_uuid, owner_type, accepted_when)
	VALUES ($1, $2, $3, $4)
	RETURNING uuid`

	err := r.db.QueryRowContext(
		ctx,
		query,
		termsAccepted.WebSiteUUID,
		termsAccepted.OwnerUUID,
		termsAccepted.OwnerType,
		termsAccepted.AcceptedWhen,
	).Scan(
		&termsAccepted.UUID,
	)

	if err != nil {
		return nil, errors.New("could not create terms accepted")
	}

	return termsAccepted, nil
}

func (r *TermsAcceptedRepository) FindTermsAcceptedByUUID(uuid string) (*domain.TermsAcceptedBy, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, website_uuid, owner_uuid, owner_type, accepted_when
	FROM terms_accepted
	WHERE uuid = $1`

	row := r.db.QueryRowContext(ctx, query, uuid)
	return helpers.ScanTermsAccepted(row)
}

func (r *TermsAcceptedRepository) FindTermsAcceptedByOwner(ownerUUID string, websiteUUID string) (*domain.TermsAcceptedBy, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, website_uuid, owner_uuid, owner_type, accepted_when
	FROM terms_accepted
	WHERE owner_uuid = $1 AND website_uuid = $2
	LIMIT 1`

	row := r.db.QueryRowContext(ctx, query, ownerUUID, websiteUUID)
	return helpers.ScanTermsAccepted(row)
}

func (r *TermsAcceptedRepository) GetTermsAcceptedFromWebsite(websiteUUID string) ([]*domain.TermsAcceptedBy, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, website_uuid, owner_uuid, owner_type, accepted_when
	FROM terms_accepted
	WHERE website_uuid = $1`

	rows, err := r.db.QueryContext(ctx, query, websiteUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return helpers.ScanTermsAcceptedSlice(rows)
}

func (r *TermsAcceptedRepository) GetTermsAccepted() ([]*domain.TermsAcceptedBy, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, website_uuid, owner_uuid, owner_type, accepted_when
	FROM terms_accepted`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return helpers.ScanTermsAcceptedSlice(rows)
}

func (r *TermsAcceptedRepository) UpdateTermsAcceptedByUUID(uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE terms_accepted SET accepted_when = NOW() WHERE uuid = $1`

	result, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("terms accepted not found")
	}

	return nil
}

func (r *TermsAcceptedRepository) DeleteTermsAcceptedByUUID(uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM terms_accepted WHERE uuid = $1`

	result, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("terms accepted not found")
	}

	return nil
}

func (r *TermsAcceptedRepository) DeleteTermsAcceptedByUUIDS(uuids []string) error {
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

	query := fmt.Sprintf(`DELETE FROM terms_accepted WHERE uuid IN (%s)`, strings.Join(placeholders, ", "))

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
