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
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO users (
			id, website_id, first_name, last_name, email, verified_email, password, role,
			cpf, display_name, birth_date, gender, bio, instagram, website_url, whatsapp, phone,
			updated_at, created_at
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,
			$9,$10,$11,$12,$13,$14,$15,$16,$17,
			$18,$19
		)`,
		user.Id, user.WebsiteId, user.First_name, user.Last_name, user.Email, user.Verified_email, user.Password, user.Role,
		user.Cpf, user.DisplayName, user.BirthDate, user.Gender, user.Bio, user.Instagram, user.WebsiteUrl, user.Whatsapp, user.Phone,
		user.Updated_at, user.Created_at,
	)
	return err
}

func (r *connection) CompanyRegister(company domain.Company) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO companies (
			id, owner_user_id, company_name, trade_name, cnpj, phone,
			zip_code, address_street, address_number, address_complement, address_district,
			address_city, address_state, address_country,
			created_at, updated_at
		) VALUES (
			$1,$2,$3,$4,$5,$6,
			$7,$8,$9,$10,$11,
			$12,$13,$14,
			$15,$16
		)`,
		company.Id, company.OwnerUserId, company.CompanyName, company.TradeName, company.Cnpj, company.Phone,
		company.ZipCode, company.AddressStreet, company.AddressNumber, company.AddressComplement, company.AddressDistrict,
		company.AddressCity, company.AddressState, company.AddressCountry,
		company.CreatedAt, company.UpdatedAt,
	)
	return err
}

const userSelectCols = `
	u.id, u.website_id, u.first_name, u.last_name, u.email, u.verified_email,
	u.password, u.role, u.updated_at, u.created_at,
	u.cpf, u.avatar_url, u.display_name, u.birth_date, u.gender, u.bio,
	u.instagram, u.website_url, u.whatsapp, u.phone,
	u.is_active,
	p.name AS plan_name,
	c.id, c.company_name, c.trade_name, c.cnpj, c.phone,
	c.zip_code, c.address_street, c.address_number, c.address_complement, c.address_district,
	c.address_city, c.address_state, c.address_country,
	c.is_active, c.created_at, c.updated_at
`

const userJoins = `
	LEFT JOIN LATERAL (
		SELECT plan_id FROM payments pay
		WHERE pay.user_id = u.id AND pay.status = 'completed'
		ORDER BY pay.purchased_in DESC LIMIT 1
	) pay ON TRUE
	LEFT JOIN plans p ON p.id = pay.plan_id
	LEFT JOIN companies c ON c.owner_user_id = u.id AND c.is_active = TRUE
`

func scanUser(row *sql.Row) (*domain.User, error) {
	var user domain.User
	var planName, cpf, avatarUrl, displayName sql.NullString
	var gender, bio, instagram, websiteUrl, whatsapp, phone sql.NullString
	var birthDate sql.NullTime

	var companyID, companyName, tradeName, cnpj, companyPhone sql.NullString
	var companyZip, companyStreet, companyNumber, companyComplement, companyDistrict sql.NullString
	var companyCity, companyState, companyCountry sql.NullString
	var companyIsActive sql.NullBool
	var companyCreatedAt, companyUpdatedAt sql.NullTime

	err := row.Scan(
		&user.Id, &user.WebsiteId, &user.First_name, &user.Last_name, &user.Email,
		&user.Verified_email, &user.Password, &user.Role, &user.Updated_at, &user.Created_at,
		&cpf, &avatarUrl, &displayName, &birthDate, &gender, &bio,
		&instagram, &websiteUrl, &whatsapp, &phone,
		&user.IsActive,
		&planName,
		&companyID, &companyName, &tradeName, &cnpj, &companyPhone,
		&companyZip, &companyStreet, &companyNumber, &companyComplement, &companyDistrict,
		&companyCity, &companyState, &companyCountry,
		&companyIsActive, &companyCreatedAt, &companyUpdatedAt,
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
	if cpf.Valid {
		user.Cpf = &cpf.String
	}
	if avatarUrl.Valid {
		user.AvatarUrl = &avatarUrl.String
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

	if companyID.Valid {
		co := &domain.Company{
			Id:          companyID.String,
			CompanyName: companyName.String,
			IsActive:    companyIsActive.Bool,
		}
		if tradeName.Valid {
			co.TradeName = &tradeName.String
		}
		if cnpj.Valid {
			co.Cnpj = &cnpj.String
		}
		if companyPhone.Valid {
			co.Phone = &companyPhone.String
		}
		if companyZip.Valid {
			co.ZipCode = &companyZip.String
		}
		if companyStreet.Valid {
			co.AddressStreet = &companyStreet.String
		}
		if companyNumber.Valid {
			co.AddressNumber = &companyNumber.String
		}
		if companyComplement.Valid {
			co.AddressComplement = &companyComplement.String
		}
		if companyDistrict.Valid {
			co.AddressDistrict = &companyDistrict.String
		}
		if companyCity.Valid {
			co.AddressCity = &companyCity.String
		}
		if companyState.Valid {
			co.AddressState = &companyState.String
		}
		if companyCountry.Valid {
			co.AddressCountry = &companyCountry.String
		}
		if companyCreatedAt.Valid {
			co.CreatedAt = companyCreatedAt.Time
		}
		if companyUpdatedAt.Valid {
			co.UpdatedAt = companyUpdatedAt.Time
		}
		user.Company = co
	}

	return &user, nil
}

func (r *connection) FindUserByID(id string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(ctx, `
		SELECT `+userSelectCols+`
		FROM users u
		`+userJoins+`
		WHERE u.id = $1 AND u.is_active = TRUE
	`, id)
	return scanUser(row)
}

func (r *connection) FindUserByEmail(email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(ctx, `
		SELECT `+userSelectCols+`
		FROM users u
		`+userJoins+`
		WHERE u.email = $1 AND u.is_active = TRUE
	`, email)
	return scanUser(row)
}

func (r *connection) FindUserByEmailAndWebsite(email, websiteID string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(ctx, `
		SELECT `+userSelectCols+`
		FROM users u
		`+userJoins+`
		WHERE u.email = $1 AND u.website_id = $2 AND u.is_active = TRUE
	`, email, websiteID)
	return scanUser(row)
}

func (r *connection) FindCompanyByOwnerUserID(ownerUserID string) (*domain.Company, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var co domain.Company
	var tradeName, cnpj, phone, zipCode, street, number, complement, district, city, state, country sql.NullString

	err := r.db.QueryRowContext(ctx, `
		SELECT id, owner_user_id, company_name, trade_name, cnpj, phone,
		       zip_code, address_street, address_number, address_complement, address_district,
		       address_city, address_state, address_country,
		       is_active, created_at, updated_at
		FROM companies
		WHERE owner_user_id = $1 AND is_active = TRUE
		LIMIT 1
	`, ownerUserID).Scan(
		&co.Id, &co.OwnerUserId, &co.CompanyName,
		&tradeName, &cnpj, &phone,
		&zipCode, &street, &number, &complement, &district,
		&city, &state, &country,
		&co.IsActive, &co.CreatedAt, &co.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if tradeName.Valid {
		co.TradeName = &tradeName.String
	}
	if cnpj.Valid {
		co.Cnpj = &cnpj.String
	}
	if phone.Valid {
		co.Phone = &phone.String
	}
	if zipCode.Valid {
		co.ZipCode = &zipCode.String
	}
	if street.Valid {
		co.AddressStreet = &street.String
	}
	if number.Valid {
		co.AddressNumber = &number.String
	}
	if complement.Valid {
		co.AddressComplement = &complement.String
	}
	if district.Valid {
		co.AddressDistrict = &district.String
	}
	if city.Valid {
		co.AddressCity = &city.String
	}
	if state.Valid {
		co.AddressState = &state.String
	}
	if country.Valid {
		co.AddressCountry = &country.String
	}
	return &co, nil
}

func (r *connection) UpdateUserProfile(id string, data domain.UpdateProfileData) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.db.ExecContext(ctx, `
		UPDATE users SET
			first_name=$1, last_name=$2, cpf=$3, avatar_url=$4,
			display_name=$5, birth_date=$6, gender=$7, bio=$8,
			instagram=$9, website_url=$10, whatsapp=$11, phone=$12,
			updated_at=NOW()
		WHERE id=$13`,
		data.FirstName, data.LastName, data.Cpf, data.AvatarUrl,
		data.DisplayName, data.BirthDate, data.Gender, data.Bio,
		data.Instagram, data.WebsiteUrl, data.Whatsapp, data.Phone,
		id,
	)
	return err
}

func (r *connection) ListUserAddresses(userID string) ([]*domain.UserAddress, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, label, zip_code, street, number, complement, district, city, state, country, is_default, created_at, updated_at
		FROM user_addresses
		WHERE user_id = $1
		ORDER BY is_default DESC, created_at ASC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addrs []*domain.UserAddress
	for rows.Next() {
		a := &domain.UserAddress{}
		var label, zipCode, street, number, complement, district, city, state sql.NullString
		if err := rows.Scan(
			&a.Id, &a.UserId, &label, &zipCode, &street, &number, &complement,
			&district, &city, &state, &a.Country, &a.IsDefault, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if label.Valid {
			a.Label = &label.String
		}
		if zipCode.Valid {
			a.ZipCode = &zipCode.String
		}
		if street.Valid {
			a.Street = &street.String
		}
		if number.Valid {
			a.Number = &number.String
		}
		if complement.Valid {
			a.Complement = &complement.String
		}
		if district.Valid {
			a.District = &district.String
		}
		if city.Valid {
			a.City = &city.String
		}
		if state.Valid {
			a.State = &state.String
		}
		addrs = append(addrs, a)
	}
	return addrs, rows.Err()
}

func (r *connection) GetDefaultUserAddress(userID string) (*domain.UserAddress, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	a := &domain.UserAddress{}
	var label, zipCode, street, number, complement, district, city, state sql.NullString
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, label, zip_code, street, number, complement, district, city, state, country, is_default, created_at, updated_at
		FROM user_addresses
		WHERE user_id = $1 AND is_default = TRUE
		LIMIT 1
	`, userID).Scan(
		&a.Id, &a.UserId, &label, &zipCode, &street, &number, &complement,
		&district, &city, &state, &a.Country, &a.IsDefault, &a.CreatedAt, &a.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if label.Valid {
		a.Label = &label.String
	}
	if zipCode.Valid {
		a.ZipCode = &zipCode.String
	}
	if street.Valid {
		a.Street = &street.String
	}
	if number.Valid {
		a.Number = &number.String
	}
	if complement.Valid {
		a.Complement = &complement.String
	}
	if district.Valid {
		a.District = &district.String
	}
	if city.Valid {
		a.City = &city.String
	}
	if state.Valid {
		a.State = &state.String
	}
	return a, nil
}

func (r *connection) CreateUserAddress(addr domain.UserAddress) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if addr.IsDefault {
		if _, err := tx.ExecContext(ctx, `UPDATE user_addresses SET is_default = FALSE WHERE user_id = $1`, addr.UserId); err != nil {
			return err
		}
	}
	_, err = tx.ExecContext(ctx, `
		INSERT INTO user_addresses (id, user_id, label, zip_code, street, number, complement, district, city, state, country, is_default, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,NOW(),NOW())
	`, addr.Id, addr.UserId, addr.Label, addr.ZipCode, addr.Street, addr.Number, addr.Complement,
		addr.District, addr.City, addr.State, addr.Country, addr.IsDefault,
	)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *connection) UpdateUserAddress(id, userID string, data domain.UpsertAddressData) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	country := "BR"
	if data.Country != nil && *data.Country != "" {
		country = *data.Country
	}
	_, err := r.db.ExecContext(ctx, `
		UPDATE user_addresses SET
			label=$1, zip_code=$2, street=$3, number=$4, complement=$5,
			district=$6, city=$7, state=$8, country=$9, updated_at=NOW()
		WHERE id=$10 AND user_id=$11
	`, data.Label, data.ZipCode, data.Street, data.Number, data.Complement,
		data.District, data.City, data.State, country, id, userID,
	)
	return err
}

func (r *connection) DeleteUserAddress(id, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.db.ExecContext(ctx, `DELETE FROM user_addresses WHERE id=$1 AND user_id=$2`, id, userID)
	return err
}

func (r *connection) SetDefaultAddress(id, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `UPDATE user_addresses SET is_default = FALSE WHERE user_id = $1`, userID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE user_addresses SET is_default = TRUE WHERE id = $1 AND user_id = $2`, id, userID); err != nil {
		return err
	}
	return tx.Commit()
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

func (r *connection) DeactivateUserByID(id string, deleteAfter time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, `
		UPDATE users
		SET is_active = FALSE,
		    deactivated_at = NOW(),
		    delete_after = $2,
		    updated_at = NOW()
		WHERE id = $1 AND is_active = TRUE
	`, id, deleteAfter)
	return err
}

func (r *connection) DeleteExpiredUnverifiedUsers() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, `
		DELETE FROM users
		WHERE verified_email = FALSE
		  AND created_at < NOW() - INTERVAL '10 minutes'
	`)
	return err
}

func (r *connection) DeleteExpiredDeactivatedUsers() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `
		DELETE FROM websites
		WHERE creator_id IN (
			SELECT id FROM users
			WHERE is_active = FALSE
			  AND delete_after IS NOT NULL
			  AND delete_after <= NOW()
		)
	`); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `
		DELETE FROM users
		WHERE is_active = FALSE
		  AND delete_after IS NOT NULL
		  AND delete_after <= NOW()
	`); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *connection) UpdateVerifiedEmailToTrue(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.db.ExecContext(ctx, `UPDATE users SET verified_email = TRUE WHERE id = $1`, id)
	return err
}

func (r *connection) UpdateVerifiedEmailToTrueByWebsite(id, websiteID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.db.ExecContext(ctx, `UPDATE users SET verified_email = TRUE WHERE id = $1 AND website_id = $2`, id, websiteID)
	return err
}
