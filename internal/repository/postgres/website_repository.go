package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ViitoJooj/Jesterx/internal/domain"
)

func NewWebSiteRepository(db *sql.DB) *connection {
	return &connection{db: db}
}

func (r *connection) SaveWebSite(website domain.WebSite) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `INSERT INTO websites (id, website_type, image, name, short_description, description, creator_id, banned, updated_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.ExecContext(ctx, query, website.Id, website.Type, website.Image, website.Name, website.Short_description, website.Description, website.Creator_id, website.Banned, website.Updated_at, website.Created_at)
	return err
}

func (r *connection) FindWebSiteByID(id string) (*domain.WebSite, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `SELECT id, website_type, image, name, short_description, description, creator_id, banned, updated_at, created_at, mature_content, rating_avg, rating_count FROM websites WHERE id = $1`
	var website domain.WebSite
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&website.Id, &website.Type, &website.Image, &website.Name, &website.Short_description, &website.Description, &website.Creator_id, &website.Banned, &website.Updated_at, &website.Created_at, &website.MatureContent, &website.RatingAvg, &website.RatingCount)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &website, nil
}

func (r *connection) ListWebSitesByUserID(id string) ([]domain.WebSite, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, website_type, image, name, short_description, description, creator_id, banned, updated_at, created_at
		FROM websites
		WHERE creator_id = $1
		ORDER BY created_at DESC
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	websites := make([]domain.WebSite, 0)
	for rows.Next() {
		var website domain.WebSite
		if err := rows.Scan(
			&website.Id,
			&website.Type,
			&website.Image,
			&website.Name,
			&website.Short_description,
			&website.Description,
			&website.Creator_id,
			&website.Banned,
			&website.Updated_at,
			&website.Created_at,
		); err != nil {
			return nil, err
		}
		websites = append(websites, website)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return websites, nil
}

func (r *connection) FindWebSiteByName(name string) (*domain.WebSite, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, website_type, image, name, short_description, description, creator_id, banned, updated_at, created_at, mature_content, rating_avg, rating_count FROM websites WHERE name = $1`

	var website domain.WebSite
	err := r.db.QueryRowContext(ctx, query, name).
		Scan(&website.Id, &website.Type, &website.Image, &website.Name, &website.Short_description, &website.Description, &website.Creator_id, &website.Banned, &website.Updated_at, &website.Created_at, &website.MatureContent, &website.RatingAvg, &website.RatingCount)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &website, nil
}

func (r *connection) UpdateWebSiteByID(id string, website domain.WebSite) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `UPDATE websites SET website_type = $1, image = $2, name = $3, short_description = $4, description = $5, creator_id = $6, banned = $7, updated_at = $8 WHERE id = $9`
	_, err := r.db.ExecContext(ctx, query, website.Type, website.Image, website.Name, website.Short_description, website.Description, website.Creator_id, website.Banned, website.Updated_at, id)
	return err
}

func (r *connection) DeleteWebSiteByID(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM website_versions WHERE website_id = $1`, id); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM website_routes WHERE website_id = $1`, id); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM websites WHERE id = $1`, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *connection) CountWebSitesByUserID(userID string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM websites WHERE creator_id=$1`, userID).Scan(&count)
	return count, err
}

func (r *connection) ReplaceRoutesByWebsiteID(websiteID string, routes []domain.WebSiteRoute) error {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM website_routes WHERE website_id = $1`, websiteID); err != nil {
		tx.Rollback()
		return err
	}

	for _, route := range routes {
		_, err := tx.ExecContext(
			ctx,
			`INSERT INTO website_routes (id, website_id, path, title, requires_auth, position, updated_at, created_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			route.Id,
			route.WebsiteId,
			route.Path,
			route.Title,
			route.RequiresAuth,
			route.Position,
			route.Updated_at,
			route.Created_at,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *connection) ListRoutesByWebsiteID(websiteID string) ([]domain.WebSiteRoute, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, website_id, path, title, requires_auth, position, updated_at, created_at
		FROM website_routes
		WHERE website_id = $1
		ORDER BY position ASC, created_at ASC
	`, websiteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	routes := make([]domain.WebSiteRoute, 0)
	for rows.Next() {
		var route domain.WebSiteRoute
		if err := rows.Scan(
			&route.Id,
			&route.WebsiteId,
			&route.Path,
			&route.Title,
			&route.RequiresAuth,
			&route.Position,
			&route.Updated_at,
			&route.Created_at,
		); err != nil {
			return nil, err
		}
		routes = append(routes, route)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return routes, nil
}

func (r *connection) FindRouteByWebsiteIDAndPath(websiteID string, path string) (*domain.WebSiteRoute, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var route domain.WebSiteRoute
	err := r.db.QueryRowContext(ctx, `
		SELECT id, website_id, path, title, requires_auth, position, updated_at, created_at
		FROM website_routes
		WHERE website_id = $1 AND path = $2
		LIMIT 1
	`, websiteID, path).Scan(
		&route.Id,
		&route.WebsiteId,
		&route.Path,
		&route.Title,
		&route.RequiresAuth,
		&route.Position,
		&route.Updated_at,
		&route.Created_at,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &route, nil
}

func (r *connection) SaveVersion(version domain.WebSiteVersion) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO website_versions
			(id, website_id, version, source_type, source, compiled_html, scan_status, scan_score, scan_findings, published, published_at, created_by, updated_at, created_at)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`,
		version.Id,
		version.WebsiteId,
		version.Version,
		version.SourceType,
		version.Source,
		version.CompiledHTML,
		version.ScanStatus,
		version.ScanScore,
		version.ScanFindings,
		version.Published,
		version.PublishedAt,
		version.CreatedBy,
		version.Updated_at,
		version.Created_at,
	)
	return err
}

func (r *connection) DeleteVersionsByWebsiteID(websiteID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.db.ExecContext(ctx, `DELETE FROM website_versions WHERE website_id = $1`, websiteID)
	return err
}

func (r *connection) FindLatestVersionByWebsiteID(websiteID string) (*domain.WebSiteVersion, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, website_id, version, source_type, source, compiled_html, scan_status, scan_score, scan_findings, published, published_at, created_by, updated_at, created_at
		FROM website_versions
		WHERE website_id = $1
		ORDER BY version DESC
		LIMIT 1
	`

	var version domain.WebSiteVersion
	err := r.db.QueryRowContext(ctx, query, websiteID).Scan(
		&version.Id,
		&version.WebsiteId,
		&version.Version,
		&version.SourceType,
		&version.Source,
		&version.CompiledHTML,
		&version.ScanStatus,
		&version.ScanScore,
		&version.ScanFindings,
		&version.Published,
		&version.PublishedAt,
		&version.CreatedBy,
		&version.Updated_at,
		&version.Created_at,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &version, nil
}

func (r *connection) FindVersionByWebsiteID(websiteID string, versionNumber int) (*domain.WebSiteVersion, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, website_id, version, source_type, source, compiled_html, scan_status, scan_score, scan_findings, published, published_at, created_by, updated_at, created_at
		FROM website_versions
		WHERE website_id = $1 AND version = $2
		LIMIT 1
	`

	var version domain.WebSiteVersion
	err := r.db.QueryRowContext(ctx, query, websiteID, versionNumber).Scan(
		&version.Id,
		&version.WebsiteId,
		&version.Version,
		&version.SourceType,
		&version.Source,
		&version.CompiledHTML,
		&version.ScanStatus,
		&version.ScanScore,
		&version.ScanFindings,
		&version.Published,
		&version.PublishedAt,
		&version.CreatedBy,
		&version.Updated_at,
		&version.Created_at,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &version, nil
}

func (r *connection) ListVersionsByWebsiteID(websiteID string) ([]domain.WebSiteVersion, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, website_id, version, source_type, source, compiled_html, scan_status, scan_score, scan_findings, published, published_at, created_by, updated_at, created_at
		FROM website_versions
		WHERE website_id = $1
		ORDER BY version DESC
	`, websiteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	versions := make([]domain.WebSiteVersion, 0)
	for rows.Next() {
		var version domain.WebSiteVersion
		if err := rows.Scan(
			&version.Id,
			&version.WebsiteId,
			&version.Version,
			&version.SourceType,
			&version.Source,
			&version.CompiledHTML,
			&version.ScanStatus,
			&version.ScanScore,
			&version.ScanFindings,
			&version.Published,
			&version.PublishedAt,
			&version.CreatedBy,
			&version.Updated_at,
			&version.Created_at,
		); err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return versions, nil
}

func (r *connection) FindPublishedVersionByWebsiteID(websiteID string) (*domain.WebSiteVersion, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var version domain.WebSiteVersion
	err := r.db.QueryRowContext(ctx, `
		SELECT id, website_id, version, source_type, source, compiled_html, scan_status, scan_score, scan_findings, published, published_at, created_by, updated_at, created_at
		FROM website_versions
		WHERE website_id = $1 AND published = true
		ORDER BY version DESC
		LIMIT 1
	`, websiteID).Scan(
		&version.Id,
		&version.WebsiteId,
		&version.Version,
		&version.SourceType,
		&version.Source,
		&version.CompiledHTML,
		&version.ScanStatus,
		&version.ScanScore,
		&version.ScanFindings,
		&version.Published,
		&version.PublishedAt,
		&version.CreatedBy,
		&version.Updated_at,
		&version.Created_at,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &version, nil
}

func (r *connection) UpdateVersionPublishState(websiteID string, version int, published bool, publishedAt *time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if published {
		if _, err := tx.ExecContext(ctx, `UPDATE website_versions SET published = false, published_at = NULL, updated_at = NOW() WHERE website_id = $1`, websiteID); err != nil {
			tx.Rollback()
			return err
		}
	}

	if _, err := tx.ExecContext(
		ctx,
		`UPDATE website_versions
		 SET published = $1, published_at = $2, updated_at = NOW()
		 WHERE website_id = $3 AND version = $4`,
		published, publishedAt, websiteID, version,
	); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
