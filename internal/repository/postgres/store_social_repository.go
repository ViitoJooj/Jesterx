package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/ViitoJooj/Jesterx/internal/domain"
	"github.com/google/uuid"
)

type storeSocialRepo struct {
	db *sql.DB
}

func NewStoreSocialRepository(db *sql.DB) *storeSocialRepo {
	return &storeSocialRepo{db: db}
}

// GetStoreFullInfo returns website info joined with creator, latest version editor type, and managers.
func (r *storeSocialRepo) GetStoreFullInfo(websiteID string) (*domain.StoreFullInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(ctx, `
		SELECT
			w.id, w.name, w.short_description, w.description, w.image, w.website_type,
			w.mature_content, w.rating_avg, w.rating_count,
			COALESCE(wv.source_type, '') AS editor_type,
			u.id, u.first_name || ' ' || u.last_name,
			u.company_name, u.trade_name, u.avatar_url, u.account_type
		FROM websites w
		JOIN users u ON u.id = w.creator_id
		LEFT JOIN LATERAL (
			SELECT source_type FROM website_versions
			WHERE website_id = w.id
			ORDER BY created_at DESC LIMIT 1
		) wv ON true
		WHERE w.id = $1 AND w.banned = false`, websiteID)

	var info domain.StoreFullInfo
	var companyName, tradeName, avatarURL sql.NullString

	err := row.Scan(
		&info.ID, &info.Name, &info.ShortDescription, &info.Description, &info.Image, &info.Type,
		&info.MatureContent, &info.RatingAvg, &info.RatingCount,
		&info.EditorType,
		&info.Creator.ID, &info.Creator.FullName,
		&companyName, &tradeName, &avatarURL, &info.Creator.AccountType,
	)
	if err != nil {
		return nil, err
	}
	if companyName.Valid {
		info.Creator.CompanyName = &companyName.String
	}
	if tradeName.Valid {
		info.Creator.TradeName = &tradeName.String
	}
	if avatarURL.Valid {
		info.Creator.AvatarURL = &avatarURL.String
	}

	// Fetch managers (role = 'manager') with user info
	mrows, err := r.db.QueryContext(ctx, `
		SELECT sm.id, sm.website_id, sm.user_id,
		       u.first_name || ' ' || u.last_name, u.avatar_url,
		       sm.role, sm.invited_by, sm.created_at
		FROM store_members sm
		JOIN users u ON u.id = sm.user_id
		WHERE sm.website_id = $1 AND sm.role = 'manager'
		ORDER BY sm.created_at ASC`, websiteID)
	if err != nil {
		return nil, err
	}
	defer mrows.Close()

	for mrows.Next() {
		var m domain.StoreMember
		var mAvatarURL sql.NullString
		var invitedBy sql.NullString
		if err := mrows.Scan(&m.ID, &m.WebsiteID, &m.UserID, &m.UserName, &mAvatarURL,
			&m.Role, &invitedBy, &m.CreatedAt); err != nil {
			return nil, err
		}
		if mAvatarURL.Valid {
			m.AvatarURL = &mAvatarURL.String
		}
		if invitedBy.Valid {
			m.InvitedBy = &invitedBy.String
		}
		info.Managers = append(info.Managers, m)
	}

	return &info, nil
}

// ─── Comments ────────────────────────────────────────────────────────────────

func (r *storeSocialRepo) SaveComment(comment domain.StoreComment) (*domain.StoreComment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id, _ := uuid.NewV7()
	comment.ID = id.String()
	now := time.Now()
	comment.CreatedAt = now
	comment.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO store_comments (id, website_id, user_id, content, stars, parent_comment_id, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$7)`,
		comment.ID, comment.WebsiteID, comment.UserID, comment.Content,
		nullableInt(comment.Stars), nullableStr(comment.ParentCommentID), now)
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRowContext(ctx, `
		SELECT sc.id, sc.website_id, sc.user_id,
		       u.first_name || ' ' || u.last_name, u.avatar_url,
		       sc.content, sc.stars, sc.parent_comment_id, sc.created_at, sc.updated_at
		FROM store_comments sc
		JOIN users u ON u.id = sc.user_id
		WHERE sc.id = $1`, comment.ID)

	return scanComment(row)
}

func (r *storeSocialRepo) ListComments(websiteID string) ([]domain.StoreComment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, `
		SELECT sc.id, sc.website_id, sc.user_id,
		       u.first_name || ' ' || u.last_name, u.avatar_url,
		       sc.content, sc.stars, sc.parent_comment_id, sc.created_at, sc.updated_at
		FROM store_comments sc
		JOIN users u ON u.id = sc.user_id
		WHERE sc.website_id = $1
		ORDER BY sc.parent_comment_id NULLS FIRST, sc.created_at ASC`, websiteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	byID := make(map[string]*domain.StoreComment)
	var topLevel []domain.StoreComment

	for rows.Next() {
		var c domain.StoreComment
		var avatarURL, parentID sql.NullString
		var stars sql.NullInt64
		if err := rows.Scan(&c.ID, &c.WebsiteID, &c.UserID, &c.UserName, &avatarURL,
			&c.Content, &stars, &parentID, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		if avatarURL.Valid {
			c.AvatarURL = &avatarURL.String
		}
		if parentID.Valid {
			c.ParentCommentID = &parentID.String
		}
		if stars.Valid {
			v := int(stars.Int64)
			c.Stars = &v
		}
		byID[c.ID] = &c
		if !parentID.Valid {
			topLevel = append(topLevel, c)
		}
	}

	result := make([]domain.StoreComment, 0, len(topLevel))
	for i := range topLevel {
		for _, child := range byID {
			if child.ParentCommentID != nil && *child.ParentCommentID == topLevel[i].ID {
				topLevel[i].Replies = append(topLevel[i].Replies, *child)
			}
		}
		result = append(result, topLevel[i])
	}
	return result, nil
}

func (r *storeSocialRepo) FindCommentByID(id string) (*domain.StoreComment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(ctx, `
		SELECT sc.id, sc.website_id, sc.user_id,
		       u.first_name || ' ' || u.last_name, u.avatar_url,
		       sc.content, sc.stars, sc.parent_comment_id, sc.created_at, sc.updated_at
		FROM store_comments sc
		JOIN users u ON u.id = sc.user_id
		WHERE sc.id = $1`, id)

	return scanComment(row)
}

func (r *storeSocialRepo) DeleteComment(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.db.ExecContext(ctx, `DELETE FROM store_comments WHERE id = $1`, id)
	return err
}

func scanComment(row *sql.Row) (*domain.StoreComment, error) {
	var c domain.StoreComment
	var avatarURL, parentID sql.NullString
	var stars sql.NullInt64
	err := row.Scan(&c.ID, &c.WebsiteID, &c.UserID, &c.UserName, &avatarURL,
		&c.Content, &stars, &parentID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if avatarURL.Valid {
		c.AvatarURL = &avatarURL.String
	}
	if parentID.Valid {
		c.ParentCommentID = &parentID.String
	}
	if stars.Valid {
		v := int(stars.Int64)
		c.Stars = &v
	}
	return &c, nil
}

func nullableStr(s *string) interface{} {
	if s == nil {
		return nil
	}
	return *s
}

func nullableInt(i *int) interface{} {
	if i == nil {
		return nil
	}
	return *i
}

// ─── Ratings ─────────────────────────────────────────────────────────────────

func (r *storeSocialRepo) UpsertRating(websiteID, userID string, stars int) (*domain.StoreRating, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id, _ := uuid.NewV7()
	now := time.Now()

	var rating domain.StoreRating
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO store_ratings (id, website_id, user_id, stars, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$5)
		ON CONFLICT (website_id, user_id)
		DO UPDATE SET stars = EXCLUDED.stars, updated_at = EXCLUDED.updated_at
		RETURNING id, website_id, user_id, stars, created_at, updated_at`,
		id.String(), websiteID, userID, stars, now).
		Scan(&rating.ID, &rating.WebsiteID, &rating.UserID, &rating.Stars, &rating.CreatedAt, &rating.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &rating, nil
}

func (r *storeSocialRepo) GetUserRating(websiteID, userID string) (*domain.StoreRating, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var rating domain.StoreRating
	err := r.db.QueryRowContext(ctx, `
		SELECT id, website_id, user_id, stars, created_at, updated_at
		FROM store_ratings WHERE website_id = $1 AND user_id = $2`,
		websiteID, userID).
		Scan(&rating.ID, &rating.WebsiteID, &rating.UserID, &rating.Stars, &rating.CreatedAt, &rating.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &rating, nil
}

func (r *storeSocialRepo) RecalcRating(websiteID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, `
		UPDATE websites
		SET rating_avg   = COALESCE((SELECT AVG(stars) FROM store_ratings WHERE website_id = $1), 0),
		    rating_count = (SELECT COUNT(*) FROM store_ratings WHERE website_id = $1)
		WHERE id = $1`, websiteID)
	return err
}

// ─── Visits ──────────────────────────────────────────────────────────────────

func (r *storeSocialRepo) RecordVisit(websiteID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO store_visits (website_id, visit_date, visit_count)
		VALUES ($1, CURRENT_DATE, 1)
		ON CONFLICT (website_id, visit_date)
		DO UPDATE SET visit_count = store_visits.visit_count + 1`, websiteID)
	return err
}

func (r *storeSocialRepo) GetVisitStats(websiteID string, days int) ([]domain.VisitDay, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if days <= 0 || days > 90 {
		days = 30
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT TO_CHAR(d.day, 'YYYY-MM-DD'),
		       COALESCE(sv.visit_count, 0)
		FROM generate_series(
		    CURRENT_DATE - ($2 - 1)::int,
		    CURRENT_DATE,
		    '1 day'::interval
		) AS d(day)
		LEFT JOIN store_visits sv
		       ON sv.website_id = $1 AND sv.visit_date = d.day::date
		ORDER BY d.day ASC`, websiteID, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []domain.VisitDay
	for rows.Next() {
		var v domain.VisitDay
		if err := rows.Scan(&v.Date, &v.Count); err != nil {
			return nil, err
		}
		stats = append(stats, v)
	}
	return stats, nil
}

// ─── Admin ───────────────────────────────────────────────────────────────────

func (r *storeSocialRepo) SetMatureContent(websiteID string, mature bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.db.ExecContext(ctx,
		`UPDATE websites SET mature_content = $2, updated_at = NOW() WHERE id = $1`, websiteID, mature)
	return err
}

// ─── Owner update ────────────────────────────────────────────────────────────

func (r *storeSocialRepo) UpdateStoreProfile(websiteID, name, shortDesc, description string, image []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.db.ExecContext(ctx, `
		UPDATE websites
		SET name = COALESCE(NULLIF($2,''), name),
		    short_description = COALESCE(NULLIF($3,''), short_description),
		    description = COALESCE(NULLIF($4,''), description),
		    image = CASE WHEN $5::bytea IS NOT NULL THEN $5 ELSE image END,
		    updated_at = NOW()
		WHERE id = $1`, websiteID, name, shortDesc, description, imageOrNil(image))
	return err
}

func imageOrNil(b []byte) interface{} {
	if len(b) == 0 {
		return nil
	}
	return b
}

// ─── Team members ────────────────────────────────────────────────────────────

func (r *storeSocialRepo) AddMember(member domain.StoreMember) (*domain.StoreMember, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id, _ := uuid.NewV7()
	member.ID = id.String()
	member.CreatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO store_members (id, website_id, user_id, role, invited_by, created_at)
		VALUES ($1,$2,$3,$4,$5,$6)
		ON CONFLICT (website_id, user_id) DO UPDATE SET role = EXCLUDED.role`,
		member.ID, member.WebsiteID, member.UserID, member.Role,
		nullableStr(member.InvitedBy), member.CreatedAt)
	if err != nil {
		return nil, err
	}
	return r.FindMember(member.WebsiteID, member.UserID)
}

func (r *storeSocialRepo) RemoveMember(websiteID, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM store_members WHERE website_id = $1 AND user_id = $2`, websiteID, userID)
	return err
}

func (r *storeSocialRepo) ListMembers(websiteID string) ([]domain.StoreMember, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, `
		SELECT sm.id, sm.website_id, sm.user_id,
		       u.first_name || ' ' || u.last_name, u.avatar_url,
		       sm.role, sm.invited_by, sm.created_at
		FROM store_members sm
		JOIN users u ON u.id = sm.user_id
		WHERE sm.website_id = $1
		ORDER BY sm.created_at ASC`, websiteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []domain.StoreMember
	for rows.Next() {
		m, err := scanMember(rows)
		if err != nil {
			return nil, err
		}
		members = append(members, *m)
	}
	return members, nil
}

func (r *storeSocialRepo) FindMember(websiteID, userID string) (*domain.StoreMember, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, `
		SELECT sm.id, sm.website_id, sm.user_id,
		       u.first_name || ' ' || u.last_name, u.avatar_url,
		       sm.role, sm.invited_by, sm.created_at
		FROM store_members sm
		JOIN users u ON u.id = sm.user_id
		WHERE sm.website_id = $1 AND sm.user_id = $2`, websiteID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		return scanMember(rows)
	}
	return nil, nil
}

// GetUserRoleInStore resolves the effective role of a user for a store.
// Returns "owner", a member role constant, "admin" (global), or "" for no access.
func (r *storeSocialRepo) GetUserRoleInStore(userID, websiteID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check global admin first — avoids failure if website lookup fails.
	var globalRole string
	if err := r.db.QueryRowContext(ctx,
		`SELECT COALESCE(role,'') FROM users WHERE id = $1`, userID).Scan(&globalRole); err == nil && globalRole == "admin" {
		return "admin", nil
	}

	// Check if owner
	var ownerID string
	err := r.db.QueryRowContext(ctx,
		`SELECT creator_id FROM websites WHERE id = $1`, websiteID).Scan(&ownerID)
	if err != nil {
		return "", err
	}
	if ownerID == userID {
		return "owner", nil
	}

	// Check store member role
	var memberRole string
	err = r.db.QueryRowContext(ctx,
		`SELECT role FROM store_members WHERE website_id = $1 AND user_id = $2`, websiteID, userID).
		Scan(&memberRole)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return memberRole, nil
}

// GetPlatformRole returns the user's global platform role from the users table.
func (r *storeSocialRepo) GetPlatformRole(userID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var role string
	err := r.db.QueryRowContext(ctx,
		`SELECT COALESCE(role,'user') FROM users WHERE id = $1`, userID).Scan(&role)
	if err == sql.ErrNoRows {
		return "user", nil
	}
	return role, err
}

func scanMember(rows *sql.Rows) (*domain.StoreMember, error) {
	var m domain.StoreMember
	var avatarURL, invitedBy sql.NullString
	err := rows.Scan(&m.ID, &m.WebsiteID, &m.UserID, &m.UserName, &avatarURL,
		&m.Role, &invitedBy, &m.CreatedAt)
	if err != nil {
		return nil, err
	}
	if avatarURL.Valid {
		m.AvatarURL = &avatarURL.String
	}
	if invitedBy.Valid {
		m.InvitedBy = &invitedBy.String
	}
	return &m, nil
}
