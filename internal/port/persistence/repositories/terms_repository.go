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

var _ contracts.TermsContract = (*TermsRepository)(nil)

type TermsRepository struct {
	db *sql.DB
}

func NewTermsRepository(db *sql.DB) *TermsRepository {
	return &TermsRepository{
		db: db,
	}
}

func (r *TermsRepository) CreateTerms(terms *domain.Terms) (*domain.Terms, error) {
	if terms == nil {
		return nil, errors.New("invalid terms")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO terms (name, description)
	VALUES ($1, $2)
	RETURNING uuid, created_at, updated_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		terms.Name,
		terms.Description,
	).Scan(
		&terms.UUID,
		&terms.CreatedAt,
		&terms.UpdatedAt,
	)

	if err != nil {
		return nil, errors.New("could not create terms")
	}

	return terms, nil
}

func (r *TermsRepository) FindTermsByUUID(uuid string) (*domain.Terms, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, name, description, created_at, updated_at
	FROM terms
	WHERE uuid = $1`

	row := r.db.QueryRowContext(ctx, query, uuid)
	return helpers.ScanTerms(row)
}

func (r *TermsRepository) FindTermsByName(name string) (*domain.Terms, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, name, description, created_at, updated_at
	FROM terms
	WHERE name = $1`

	row := r.db.QueryRowContext(ctx, query, name)
	return helpers.ScanTerms(row)
}

func (r *TermsRepository) GetTerms() ([]*domain.Terms, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, name, description, created_at, updated_at
	FROM terms`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return helpers.ScanTermsSlice(rows)
}

func (r *TermsRepository) UpdateTermsByUUID(uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE terms SET updated_at = NOW() WHERE uuid = $1`

	result, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("terms not found")
	}

	return nil
}

func (r *TermsRepository) DeleteTermsByUUID(uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM terms WHERE uuid = $1`

	result, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("terms not found")
	}

	return nil
}

func (r *TermsRepository) DeleteTermsByUUIDS(uuids []string) error {
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

	query := fmt.Sprintf(`DELETE FROM terms WHERE uuid IN (%s)`, strings.Join(placeholders, ", "))

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
