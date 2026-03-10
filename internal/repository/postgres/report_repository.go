package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/ViitoJooj/Jesterx/internal/domain"
	"github.com/google/uuid"
)

type reportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *reportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) SaveReport(report domain.Report) (*domain.Report, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id, _ := uuid.NewV7()
	report.ID = id.String()
	now := time.Now()
	report.CreatedAt = now
	report.UpdatedAt = now
	report.Status = domain.ReportStatusOpen

	row := r.db.QueryRowContext(ctx, `
		INSERT INTO reports (id, website_id, reporter_name, reporter_email, reason, description, status, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id, ticket_number, website_id, reporter_name, reporter_email, reason, description, status,
		          admin_response, resolved_at, created_at, updated_at`,
		report.ID, report.WebsiteID, report.ReporterName, report.ReporterEmail,
		string(report.Reason), report.Description, string(report.Status),
		report.CreatedAt, report.UpdatedAt,
	)

	return scanReport(row)
}

func (r *reportRepository) FindReportByID(id string) (*domain.Report, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(ctx, `
		SELECT id, ticket_number, website_id, reporter_name, reporter_email, reason, description, status,
		       admin_response, resolved_at, created_at, updated_at
		FROM reports WHERE id = $1`, id)

	return scanReport(row)
}

func (r *reportRepository) ListReports(status string, page, perPage int) ([]domain.Report, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	var totalRows *sql.Rows
	var dataRows *sql.Rows
	var err error

	if status != "" {
		totalRows, err = r.db.QueryContext(ctx, `SELECT COUNT(*) FROM reports WHERE status = $1`, status)
	} else {
		totalRows, err = r.db.QueryContext(ctx, `SELECT COUNT(*) FROM reports`)
	}
	if err != nil {
		return nil, 0, err
	}
	defer totalRows.Close()

	var total int
	if totalRows.Next() {
		_ = totalRows.Scan(&total)
	}

	if status != "" {
		dataRows, err = r.db.QueryContext(ctx, `
			SELECT id, ticket_number, website_id, reporter_name, reporter_email, reason, description, status,
			       admin_response, resolved_at, created_at, updated_at
			FROM reports WHERE status = $1
			ORDER BY created_at DESC LIMIT $2 OFFSET $3`, status, perPage, offset)
	} else {
		dataRows, err = r.db.QueryContext(ctx, `
			SELECT id, ticket_number, website_id, reporter_name, reporter_email, reason, description, status,
			       admin_response, resolved_at, created_at, updated_at
			FROM reports
			ORDER BY created_at DESC LIMIT $1 OFFSET $2`, perPage, offset)
	}
	if err != nil {
		return nil, 0, err
	}
	defer dataRows.Close()

	var reports []domain.Report
	for dataRows.Next() {
		rep, err := scanReportRow(dataRows)
		if err != nil {
			return nil, 0, err
		}
		reports = append(reports, *rep)
	}

	return reports, total, nil
}

func (r *reportRepository) ListReportsByWebsiteID(websiteID string) ([]domain.Report, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, ticket_number, website_id, reporter_name, reporter_email, reason, description, status,
		       admin_response, resolved_at, created_at, updated_at
		FROM reports WHERE website_id = $1
		ORDER BY created_at DESC`, websiteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []domain.Report
	for rows.Next() {
		rep, err := scanReportRow(rows)
		if err != nil {
			return nil, err
		}
		reports = append(reports, *rep)
	}
	return reports, nil
}

func (r *reportRepository) UpdateReport(id string, status domain.ReportStatus, adminResponse *string) (*domain.Report, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now()
	var resolvedAt *time.Time
	if status == domain.ReportStatusResolved || status == domain.ReportStatusDismissed {
		resolvedAt = &now
	}

	row := r.db.QueryRowContext(ctx, `
		UPDATE reports
		SET status = $2, admin_response = $3, resolved_at = $4, updated_at = $5
		WHERE id = $1
		RETURNING id, ticket_number, website_id, reporter_name, reporter_email, reason, description, status,
		          admin_response, resolved_at, created_at, updated_at`,
		id, string(status), adminResponse, resolvedAt, now,
	)

	return scanReport(row)
}

func scanReport(row *sql.Row) (*domain.Report, error) {
	var rep domain.Report
	var adminResponse sql.NullString
	var resolvedAt sql.NullTime

	err := row.Scan(
		&rep.ID, &rep.TicketNumber, &rep.WebsiteID,
		&rep.ReporterName, &rep.ReporterEmail,
		&rep.Reason, &rep.Description, &rep.Status,
		&adminResponse, &resolvedAt,
		&rep.CreatedAt, &rep.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if adminResponse.Valid {
		rep.AdminResponse = &adminResponse.String
	}
	if resolvedAt.Valid {
		rep.ResolvedAt = &resolvedAt.Time
	}

	return &rep, nil
}

func scanReportRow(rows *sql.Rows) (*domain.Report, error) {
	var rep domain.Report
	var adminResponse sql.NullString
	var resolvedAt sql.NullTime

	err := rows.Scan(
		&rep.ID, &rep.TicketNumber, &rep.WebsiteID,
		&rep.ReporterName, &rep.ReporterEmail,
		&rep.Reason, &rep.Description, &rep.Status,
		&adminResponse, &resolvedAt,
		&rep.CreatedAt, &rep.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if adminResponse.Valid {
		rep.AdminResponse = &adminResponse.String
	}
	if resolvedAt.Valid {
		rep.ResolvedAt = &resolvedAt.Time
	}

	return &rep, nil
}
