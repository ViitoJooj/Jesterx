package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ViitoJooj/Jesterx/internal/domain"
	"github.com/google/uuid"
)

func NewOrderRepository(db *sql.DB) *connection {
	return &connection{db: db}
}

func (r *connection) Create(order *domain.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id, _ := uuid.NewV7()
	order.ID = id.String()
	now := time.Now()
	order.CreatedAt = now
	order.UpdatedAt = now

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO orders (
			id, website_id, buyer_user_id, buyer_name, buyer_email, buyer_phone, buyer_document,
			shipping_name, shipping_phone, shipping_zip_code, shipping_address_street, shipping_address_number,
			shipping_address_complement, shipping_address_district, shipping_address_city, shipping_address_state,
			shipping_address_country, shipping_method, shipping_cost, discount_total, tax_total, currency,
			status, subtotal, platform_fee, total, notes, created_at, updated_at
		)
		VALUES (
			$1,$2,$3,$4,$5,$6,$7,
			$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,
			$23,$24,$25,$26,$27,$28,$29
		)`,
		order.ID, order.WebsiteID, nullableStringPtr(order.BuyerUserID), order.BuyerName, order.BuyerEmail,
		nullableString(order.BuyerPhone), nullableString(order.BuyerDocument),
		nullableString(order.ShippingName), nullableString(order.ShippingPhone), nullableString(order.ShippingZipCode),
		nullableString(order.ShippingAddressStreet), nullableString(order.ShippingAddressNumber),
		nullableString(order.ShippingAddressComplement), nullableString(order.ShippingAddressDistrict),
		nullableString(order.ShippingAddressCity), nullableString(order.ShippingAddressState),
		nullableString(order.ShippingAddressCountry), nullableString(order.ShippingMethod),
		order.ShippingCost, order.DiscountTotal, order.TaxTotal, order.Currency,
		string(order.Status), order.Subtotal, order.PlatformFee, order.Total,
		nullableString(order.Notes), order.CreatedAt, order.UpdatedAt,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	for i := range order.Items {
		itemID, _ := uuid.NewV7()
		order.Items[i].ID = itemID.String()
		order.Items[i].OrderID = order.ID

		_, err = tx.ExecContext(ctx, `
			INSERT INTO order_items (id, order_id, product_id, product_name, unit_price, qty, total)
			VALUES ($1,$2,$3,$4,$5,$6,$7)`,
			order.Items[i].ID, order.ID,
			order.Items[i].ProductID, order.Items[i].ProductName,
			order.Items[i].UnitPrice, order.Items[i].Qty, order.Items[i].Total,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *connection) GetByID(orderID string) (*domain.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	order, err := r.scanOrderRow(r.db.QueryRowContext(ctx, `
		SELECT id, website_id, buyer_user_id, buyer_name, buyer_email, buyer_phone, buyer_document,
		       shipping_name, shipping_phone, shipping_zip_code, shipping_address_street, shipping_address_number,
		       shipping_address_complement, shipping_address_district, shipping_address_city, shipping_address_state,
		       shipping_address_country, shipping_method, shipping_cost, discount_total, tax_total, currency,
		       status, subtotal, platform_fee, total, notes, created_at, updated_at
		FROM orders WHERE id = $1`, orderID))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	items, err := r.fetchOrderItems(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	order.Items = items
	return order, nil
}

func (r *connection) ListBySite(websiteID string) ([]domain.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, website_id, buyer_user_id, buyer_name, buyer_email, buyer_phone, buyer_document,
		       shipping_name, shipping_phone, shipping_zip_code, shipping_address_street, shipping_address_number,
		       shipping_address_complement, shipping_address_district, shipping_address_city, shipping_address_state,
		       shipping_address_country, shipping_method, shipping_cost, discount_total, tax_total, currency,
		       status, subtotal, platform_fee, total, notes, created_at, updated_at
		FROM orders WHERE website_id = $1
		ORDER BY created_at DESC`, websiteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.collectOrdersWithItems(ctx, rows)
}

func (r *connection) ListSince(from, to time.Time) ([]domain.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, website_id, buyer_user_id, buyer_name, buyer_email, buyer_phone, buyer_document,
		       shipping_name, shipping_phone, shipping_zip_code, shipping_address_street, shipping_address_number,
		       shipping_address_complement, shipping_address_district, shipping_address_city, shipping_address_state,
		       shipping_address_country, shipping_method, shipping_cost, discount_total, tax_total, currency,
		       status, subtotal, platform_fee, total, notes, created_at, updated_at
		FROM orders
		WHERE created_at >= $1 AND created_at < $2
		ORDER BY created_at DESC`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.collectOrdersWithItems(ctx, rows)
}

func (r *connection) UpdateStatus(orderID string, status domain.OrderStatus) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := r.db.ExecContext(ctx,
		`UPDATE orders SET status=$1, updated_at=NOW() WHERE id=$2`,
		string(status), orderID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("order not found")
	}
	return nil
}

// helpers

func (r *connection) scanOrderRow(row *sql.Row) (*domain.Order, error) {
	var o domain.Order
	var buyerUserID, buyerPhone, buyerDocument sql.NullString
	var shippingName, shippingPhone, shippingZipCode sql.NullString
	var shippingStreet, shippingNumber, shippingComplement sql.NullString
	var shippingDistrict, shippingCity, shippingState, shippingCountry sql.NullString
	var shippingMethod, notes sql.NullString
	err := row.Scan(
		&o.ID, &o.WebsiteID, &buyerUserID, &o.BuyerName, &o.BuyerEmail,
		&buyerPhone, &buyerDocument,
		&shippingName, &shippingPhone, &shippingZipCode, &shippingStreet, &shippingNumber,
		&shippingComplement, &shippingDistrict, &shippingCity, &shippingState,
		&shippingCountry, &shippingMethod, &o.ShippingCost, &o.DiscountTotal, &o.TaxTotal, &o.Currency,
		&o.Status, &o.Subtotal, &o.PlatformFee, &o.Total,
		&notes, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if buyerUserID.Valid {
		o.BuyerUserID = &buyerUserID.String
	}
	if buyerPhone.Valid {
		o.BuyerPhone = buyerPhone.String
	}
	if buyerDocument.Valid {
		o.BuyerDocument = buyerDocument.String
	}
	if shippingName.Valid {
		o.ShippingName = shippingName.String
	}
	if shippingPhone.Valid {
		o.ShippingPhone = shippingPhone.String
	}
	if shippingZipCode.Valid {
		o.ShippingZipCode = shippingZipCode.String
	}
	if shippingStreet.Valid {
		o.ShippingAddressStreet = shippingStreet.String
	}
	if shippingNumber.Valid {
		o.ShippingAddressNumber = shippingNumber.String
	}
	if shippingComplement.Valid {
		o.ShippingAddressComplement = shippingComplement.String
	}
	if shippingDistrict.Valid {
		o.ShippingAddressDistrict = shippingDistrict.String
	}
	if shippingCity.Valid {
		o.ShippingAddressCity = shippingCity.String
	}
	if shippingState.Valid {
		o.ShippingAddressState = shippingState.String
	}
	if shippingCountry.Valid {
		o.ShippingAddressCountry = shippingCountry.String
	}
	if shippingMethod.Valid {
		o.ShippingMethod = shippingMethod.String
	}
	if notes.Valid {
		o.Notes = notes.String
	}
	return &o, nil
}

func (r *connection) fetchOrderItems(ctx context.Context, orderID string) ([]domain.OrderItem, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, order_id, product_id, product_name, unit_price, qty, total
		 FROM order_items WHERE order_id = $1`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.OrderItem, 0)
	for rows.Next() {
		var it domain.OrderItem
		if err := rows.Scan(&it.ID, &it.OrderID, &it.ProductID, &it.ProductName, &it.UnitPrice, &it.Qty, &it.Total); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	return items, rows.Err()
}

func (r *connection) collectOrdersWithItems(ctx context.Context, rows *sql.Rows) ([]domain.Order, error) {
	orders := make([]domain.Order, 0)
	for rows.Next() {
		var o domain.Order
		var buyerUserID, buyerPhone, buyerDocument sql.NullString
		var shippingName, shippingPhone, shippingZipCode sql.NullString
		var shippingStreet, shippingNumber, shippingComplement sql.NullString
		var shippingDistrict, shippingCity, shippingState, shippingCountry sql.NullString
		var shippingMethod, notes sql.NullString
		if err := rows.Scan(
			&o.ID, &o.WebsiteID, &buyerUserID, &o.BuyerName, &o.BuyerEmail,
			&buyerPhone, &buyerDocument,
			&shippingName, &shippingPhone, &shippingZipCode, &shippingStreet, &shippingNumber,
			&shippingComplement, &shippingDistrict, &shippingCity, &shippingState,
			&shippingCountry, &shippingMethod, &o.ShippingCost, &o.DiscountTotal, &o.TaxTotal, &o.Currency,
			&o.Status, &o.Subtotal, &o.PlatformFee, &o.Total,
			&notes, &o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if buyerUserID.Valid {
			o.BuyerUserID = &buyerUserID.String
		}
		if buyerPhone.Valid {
			o.BuyerPhone = buyerPhone.String
		}
		if buyerDocument.Valid {
			o.BuyerDocument = buyerDocument.String
		}
		if shippingName.Valid {
			o.ShippingName = shippingName.String
		}
		if shippingPhone.Valid {
			o.ShippingPhone = shippingPhone.String
		}
		if shippingZipCode.Valid {
			o.ShippingZipCode = shippingZipCode.String
		}
		if shippingStreet.Valid {
			o.ShippingAddressStreet = shippingStreet.String
		}
		if shippingNumber.Valid {
			o.ShippingAddressNumber = shippingNumber.String
		}
		if shippingComplement.Valid {
			o.ShippingAddressComplement = shippingComplement.String
		}
		if shippingDistrict.Valid {
			o.ShippingAddressDistrict = shippingDistrict.String
		}
		if shippingCity.Valid {
			o.ShippingAddressCity = shippingCity.String
		}
		if shippingState.Valid {
			o.ShippingAddressState = shippingState.String
		}
		if shippingCountry.Valid {
			o.ShippingAddressCountry = shippingCountry.String
		}
		if shippingMethod.Valid {
			o.ShippingMethod = shippingMethod.String
		}
		if notes.Valid {
			o.Notes = notes.String
		}
		orders = append(orders, o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i := range orders {
		items, err := r.fetchOrderItems(ctx, orders[i].ID)
		if err != nil {
			return nil, err
		}
		orders[i].Items = items
	}
	return orders, nil
}

func nullableString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}

func nullableStringPtr(s *string) sql.NullString {
	if s == nil || *s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}
