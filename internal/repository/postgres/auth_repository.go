package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ViitoJooj/Jesterx/internal/domain"
)

func NewAuthRepository(db *sql.DB) *connection {
	return &connection{db: db}
}

func (r *connection) UserRegister(user domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `
		INSERT INTO users (
			id, website_id, first_name, last_name, email, verified_email, password, role,
			cpf_cnpj, account_type, company_name, trade_name, display_name, birth_date,
			gender, bio, instagram, website_url, whatsapp, phone,
			zip_code, address_street, address_number, address_complement, address_district,
			address_city, address_state, address_country,
			updated_at, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8,
			$9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19, $20,
			$21, $22, $23, $24, $25,
			$26, $27, $28,
			$29, $30
		)`
	_, err := r.db.ExecContext(ctx, query,
		user.Id, user.WebsiteId, user.First_name, user.Last_name, user.Email, user.Verified_email, user.Password, user.Role,
		user.CpfCnpj, user.AccountType, user.CompanyName, user.TradeName, user.DisplayName, user.BirthDate,
		user.Gender, user.Bio, user.Instagram, user.WebsiteUrl, user.Whatsapp, user.Phone,
		user.ZipCode, user.AddressStreet, user.AddressNumber, user.AddressComplement, user.AddressDistrict,
		user.AddressCity, user.AddressState, user.AddressCountry,
		user.Updated_at, user.Created_at,
	)
	return err
}

func (r *connection) FindUserByID(id string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT 
		u.id, u.website_id, u.first_name, u.last_name, u.email, u.verified_email,
		u.password, u.role, u.updated_at, u.created_at, u.cpf_cnpj, u.avatar_url,
		u.account_type, u.company_name, u.trade_name, u.display_name, u.birth_date, u.gender, u.bio,
		u.instagram, u.website_url, u.whatsapp, u.phone,
		u.zip_code, u.address_street, u.address_number, u.address_complement, u.address_district,
		u.address_city, u.address_state, u.address_country,
		p.name AS plan_name
	FROM users u
	LEFT JOIN LATERAL (
		SELECT plan_id
		FROM payments pay
		WHERE pay.user_id = u.id
		  AND pay.status = 'completed'
		ORDER BY pay.purchased_in DESC
		LIMIT 1
	) pay ON TRUE
	LEFT JOIN plans p ON p.id = pay.plan_id
	WHERE u.id = $1
	`

	var user domain.User
	var planName, cpfCnpj, avatarUrl sql.NullString
	var accountType, companyName, tradeName, displayName sql.NullString
	var gender, bio, instagram, websiteUrl, whatsapp, phone sql.NullString
	var zipCode, addressStreet, addressNumber, addressComplement, addressDistrict sql.NullString
	var addressCity, addressState, addressCountry sql.NullString
	var birthDate sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).
		Scan(
			&user.Id, &user.WebsiteId, &user.First_name, &user.Last_name, &user.Email,
			&user.Verified_email, &user.Password, &user.Role, &user.Updated_at, &user.Created_at,
			&cpfCnpj, &avatarUrl,
			&accountType, &companyName, &tradeName, &displayName, &birthDate, &gender, &bio,
			&instagram, &websiteUrl, &whatsapp, &phone,
			&zipCode, &addressStreet, &addressNumber, &addressComplement, &addressDistrict,
			&addressCity, &addressState, &addressCountry,
			&planName,
		)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if planName.Valid {
		user.Plan = &planName.String
	}
	if cpfCnpj.Valid {
		user.CpfCnpj = &cpfCnpj.String
	}
	if avatarUrl.Valid {
		user.AvatarUrl = &avatarUrl.String
	}
	user.AccountType = accountType.String
	if companyName.Valid {
		user.CompanyName = &companyName.String
	}
	if tradeName.Valid {
		user.TradeName = &tradeName.String
	}
	if displayName.Valid {
		user.DisplayName = &displayName.String
	}
	if birthDate.Valid {
		user.BirthDate = &birthDate.Time
	}
	if gender.Valid {
		user.Gender = &gender.String
	}
	if bio.Valid {
		user.Bio = &bio.String
	}
	if instagram.Valid {
		user.Instagram = &instagram.String
	}
	if websiteUrl.Valid {
		user.WebsiteUrl = &websiteUrl.String
	}
	if whatsapp.Valid {
		user.Whatsapp = &whatsapp.String
	}
	if phone.Valid {
		user.Phone = &phone.String
	}
	if zipCode.Valid {
		user.ZipCode = &zipCode.String
	}
	if addressStreet.Valid {
		user.AddressStreet = &addressStreet.String
	}
	if addressNumber.Valid {
		user.AddressNumber = &addressNumber.String
	}
	if addressComplement.Valid {
		user.AddressComplement = &addressComplement.String
	}
	if addressDistrict.Valid {
		user.AddressDistrict = &addressDistrict.String
	}
	if addressCity.Valid {
		user.AddressCity = &addressCity.String
	}
	if addressState.Valid {
		user.AddressState = &addressState.String
	}
	if addressCountry.Valid {
		user.AddressCountry = &addressCountry.String
	}

	return &user, nil
}

func (r *connection) FindUserByEmail(email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT 
		u.id, u.website_id, u.first_name, u.last_name, u.email, u.verified_email,
		u.password, u.role, u.updated_at, u.created_at, u.cpf_cnpj, u.avatar_url,
		u.account_type, u.company_name, u.trade_name, u.display_name, u.birth_date, u.gender, u.bio,
		u.instagram, u.website_url, u.whatsapp, u.phone,
		u.zip_code, u.address_street, u.address_number, u.address_complement, u.address_district,
		u.address_city, u.address_state, u.address_country,
		p.name AS plan_name
	FROM users u
	LEFT JOIN LATERAL (
		SELECT plan_id FROM payments pay
		WHERE pay.user_id = u.id AND pay.status = 'completed'
		ORDER BY pay.purchased_in DESC LIMIT 1
	) pay ON TRUE
	LEFT JOIN plans p ON p.id = pay.plan_id
	WHERE u.email = $1
	`

	var user domain.User
	var planName, cpfCnpj, avatarUrl sql.NullString
	var accountType, companyName, tradeName, displayName sql.NullString
	var gender, bio, instagram, websiteUrl, whatsapp, phone sql.NullString
	var zipCode, addressStreet, addressNumber, addressComplement, addressDistrict sql.NullString
	var addressCity, addressState, addressCountry sql.NullString
	var birthDate sql.NullTime

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.Id, &user.WebsiteId, &user.First_name, &user.Last_name, &user.Email,
		&user.Verified_email, &user.Password, &user.Role, &user.Updated_at, &user.Created_at,
		&cpfCnpj, &avatarUrl,
		&accountType, &companyName, &tradeName, &displayName, &birthDate, &gender, &bio,
		&instagram, &websiteUrl, &whatsapp, &phone,
		&zipCode, &addressStreet, &addressNumber, &addressComplement, &addressDistrict,
		&addressCity, &addressState, &addressCountry,
		&planName,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if planName.Valid {
		user.Plan = &planName.String
	}
	if cpfCnpj.Valid {
		user.CpfCnpj = &cpfCnpj.String
	}
	if avatarUrl.Valid {
		user.AvatarUrl = &avatarUrl.String
	}
	user.AccountType = accountType.String
	if companyName.Valid {
		user.CompanyName = &companyName.String
	}
	if tradeName.Valid {
		user.TradeName = &tradeName.String
	}
	if displayName.Valid {
		user.DisplayName = &displayName.String
	}
	if birthDate.Valid {
		user.BirthDate = &birthDate.Time
	}
	if gender.Valid {
		user.Gender = &gender.String
	}
	if bio.Valid {
		user.Bio = &bio.String
	}
	if instagram.Valid {
		user.Instagram = &instagram.String
	}
	if websiteUrl.Valid {
		user.WebsiteUrl = &websiteUrl.String
	}
	if whatsapp.Valid {
		user.Whatsapp = &whatsapp.String
	}
	if phone.Valid {
		user.Phone = &phone.String
	}
	if zipCode.Valid {
		user.ZipCode = &zipCode.String
	}
	if addressStreet.Valid {
		user.AddressStreet = &addressStreet.String
	}
	if addressNumber.Valid {
		user.AddressNumber = &addressNumber.String
	}
	if addressComplement.Valid {
		user.AddressComplement = &addressComplement.String
	}
	if addressDistrict.Valid {
		user.AddressDistrict = &addressDistrict.String
	}
	if addressCity.Valid {
		user.AddressCity = &addressCity.String
	}
	if addressState.Valid {
		user.AddressState = &addressState.String
	}
	if addressCountry.Valid {
		user.AddressCountry = &addressCountry.String
	}
	return &user, nil
}

func (r *connection) FindUserByEmailAndWebsite(email string, websiteId string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT 
		u.id, u.website_id, u.first_name, u.last_name, u.email, u.verified_email,
		u.password, u.role, u.updated_at, u.created_at, u.cpf_cnpj, u.avatar_url,
		u.account_type, u.company_name, u.trade_name, u.display_name, u.birth_date, u.gender, u.bio,
		u.instagram, u.website_url, u.whatsapp, u.phone,
		u.zip_code, u.address_street, u.address_number, u.address_complement, u.address_district,
		u.address_city, u.address_state, u.address_country,
		p.name AS plan_name
	FROM users u
	LEFT JOIN LATERAL (
		SELECT plan_id FROM payments pay
		WHERE pay.user_id = u.id AND pay.status = 'completed'
		ORDER BY pay.purchased_in DESC LIMIT 1
	) pay ON TRUE
	LEFT JOIN plans p ON p.id = pay.plan_id
	WHERE u.email = $1 AND u.website_id = $2
	`

	var user domain.User
	var planName, cpfCnpj, avatarUrl sql.NullString
	var accountType, companyName, tradeName, displayName sql.NullString
	var gender, bio, instagram, websiteUrl, whatsapp, phone sql.NullString
	var zipCode, addressStreet, addressNumber, addressComplement, addressDistrict sql.NullString
	var addressCity, addressState, addressCountry sql.NullString
	var birthDate sql.NullTime

	err := r.db.QueryRowContext(ctx, query, email, websiteId).Scan(
		&user.Id, &user.WebsiteId, &user.First_name, &user.Last_name, &user.Email,
		&user.Verified_email, &user.Password, &user.Role, &user.Updated_at, &user.Created_at,
		&cpfCnpj, &avatarUrl,
		&accountType, &companyName, &tradeName, &displayName, &birthDate, &gender, &bio,
		&instagram, &websiteUrl, &whatsapp, &phone,
		&zipCode, &addressStreet, &addressNumber, &addressComplement, &addressDistrict,
		&addressCity, &addressState, &addressCountry,
		&planName,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if planName.Valid {
		user.Plan = &planName.String
	}
	if cpfCnpj.Valid {
		user.CpfCnpj = &cpfCnpj.String
	}
	if avatarUrl.Valid {
		user.AvatarUrl = &avatarUrl.String
	}
	user.AccountType = accountType.String
	if companyName.Valid {
		user.CompanyName = &companyName.String
	}
	if tradeName.Valid {
		user.TradeName = &tradeName.String
	}
	if displayName.Valid {
		user.DisplayName = &displayName.String
	}
	if birthDate.Valid {
		user.BirthDate = &birthDate.Time
	}
	if gender.Valid {
		user.Gender = &gender.String
	}
	if bio.Valid {
		user.Bio = &bio.String
	}
	if instagram.Valid {
		user.Instagram = &instagram.String
	}
	if websiteUrl.Valid {
		user.WebsiteUrl = &websiteUrl.String
	}
	if whatsapp.Valid {
		user.Whatsapp = &whatsapp.String
	}
	if phone.Valid {
		user.Phone = &phone.String
	}
	if zipCode.Valid {
		user.ZipCode = &zipCode.String
	}
	if addressStreet.Valid {
		user.AddressStreet = &addressStreet.String
	}
	if addressNumber.Valid {
		user.AddressNumber = &addressNumber.String
	}
	if addressComplement.Valid {
		user.AddressComplement = &addressComplement.String
	}
	if addressDistrict.Valid {
		user.AddressDistrict = &addressDistrict.String
	}
	if addressCity.Valid {
		user.AddressCity = &addressCity.String
	}
	if addressState.Valid {
		user.AddressState = &addressState.String
	}
	if addressCountry.Valid {
		user.AddressCountry = &addressCountry.String
	}
	return &user, nil
}

func (r *connection) UpdateUserProfile(id string, data domain.UpdateProfileData) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.db.ExecContext(ctx, `
		UPDATE users SET
			first_name=$1, last_name=$2, cpf_cnpj=$3, avatar_url=$4,
			company_name=$5, trade_name=$6, display_name=$7, birth_date=$8,
			gender=$9, bio=$10, instagram=$11, website_url=$12, whatsapp=$13, phone=$14,
			zip_code=$15, address_street=$16, address_number=$17,
			address_complement=$18, address_district=$19, address_city=$20, address_state=$21,
			address_country=$22, updated_at=NOW()
		WHERE id=$23`,
		data.FirstName, data.LastName, data.CpfCnpj, data.AvatarUrl,
		data.CompanyName, data.TradeName, data.DisplayName, data.BirthDate,
		data.Gender, data.Bio, data.Instagram, data.WebsiteUrl, data.Whatsapp, data.Phone,
		data.ZipCode, data.AddressStreet, data.AddressNumber,
		data.AddressComplement, data.AddressDistrict, data.AddressCity, data.AddressState,
		data.AddressCountry,
		id,
	)
	return err
}

func (r *connection) DeleteUserByID(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM websites WHERE creator_id = $1`, id); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *connection) DeleteExpiredUnverifiedUsers() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, `
		DELETE FROM users
		WHERE verified_email = false
		AND created_at < NOW() - INTERVAL '10 minutes'
	`)
	return err
}

func (r *connection) UpdateVerifiedEmailToTrue(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `UPDATE users SET verified_email = true WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *connection) UpdateVerifiedEmailToTrueByWebsite(id string, websiteID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `UPDATE users SET verified_email = true WHERE id = $1 AND website_id = $2`
	_, err := r.db.ExecContext(ctx, query, id, websiteID)
	return err
}
