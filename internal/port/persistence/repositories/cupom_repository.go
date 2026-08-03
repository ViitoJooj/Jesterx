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

var _ contracts.CupomContract = (*CupomRepository)(nil)

type CupomRepository struct {
	db *sql.DB
}

func NewCupomRepository(db *sql.DB) *CupomRepository {
	return &CupomRepository{
		db: db,
	}
}

func (r *CupomRepository) CreateCupom(cupom *domain.Cupons) (*domain.Cupons, error) {
	if cupom == nil {
		return nil, errors.New("invalid cupom")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO cupons (tag_uuid, label, description, value, value_type)
	VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.ExecContext(
		ctx,
		query,
		cupom.TagUUID,
		cupom.Label,
		cupom.Description,
		cupom.Value,
		cupom.ValueType,
	)

	if err != nil {
		return nil, errors.New("could not create cupom")
	}

	return cupom, nil
}

func (r *CupomRepository) FindCupomByUUID(uuid string) (*domain.Cupons, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, tag_uuid, label, description, value, value_type
	FROM cupons
	WHERE uuid = $1`

	row := r.db.QueryRowContext(ctx, query, uuid)
	return helpers.ScanCupom(row)
}

func (r *CupomRepository) FindCupomByLabel(label string) (*domain.Cupons, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, tag_uuid, label, description, value, value_type
	FROM cupons
	WHERE label = $1`

	row := r.db.QueryRowContext(ctx, query, label)
	return helpers.ScanCupom(row)
}

func (r *CupomRepository) GetCuponsFromTag(tagUUID string) ([]*domain.Cupons, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, tag_uuid, label, description, value, value_type
	FROM cupons
	WHERE tag_uuid = $1`

	rows, err := r.db.QueryContext(ctx, query, tagUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return helpers.ScanCupoms(rows)
}

func (r *CupomRepository) GetCupons() ([]*domain.Cupons, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, tag_uuid, label, description, value, value_type
	FROM cupons`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return helpers.ScanCupoms(rows)
}

func (r *CupomRepository) UpdateCupomByUUID(uuid string) error {
	return nil
}

func (r *CupomRepository) DeleteCupomByUUID(uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM cupons WHERE uuid = $1`

	result, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("cupom not found")
	}

	return nil
}

func (r *CupomRepository) DeleteCupomsByUUIDS(uuids []string) error {
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

	query := fmt.Sprintf(`DELETE FROM cupons WHERE uuid IN (%s)`, strings.Join(placeholders, ", "))

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
