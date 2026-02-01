package postgres

import (
	"database/sql"
	"fmt"

	"github.com/dgentikian/rent-receipt-generator/internal/models"
)

type ReceiptRepo struct {
	db *sql.DB
}

func NewReceiptRepository(db *sql.DB) *ReceiptRepo {
	return &ReceiptRepo{db: db}
}

func (r *ReceiptRepo) Create(receipt *models.Receipt) error {
	query := `
		INSERT INTO receipts (landlord_id, property_id, tenant_id, receipt_number, period_month, period_year, 
			rent_amount, charges_amount, total_amount, payment_method, payment_date, notes, pdf_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at
	`
	return r.db.QueryRow(
		query,
		receipt.LandlordID,
		receipt.PropertyID,
		receipt.TenantID,
		receipt.ReceiptNumber,
		receipt.PeriodMonth,
		receipt.PeriodYear,
		receipt.RentAmount,
		receipt.ChargesAmount,
		receipt.TotalAmount,
		receipt.PaymentMethod,
		receipt.PaymentDate,
		receipt.Notes,
		receipt.PDFURL,
	).Scan(&receipt.ID, &receipt.CreatedAt)
}

func (r *ReceiptRepo) GetByID(id int) (*models.Receipt, error) {
	query := `
		SELECT id, landlord_id, property_id, tenant_id, receipt_number, period_month, period_year,
			rent_amount, charges_amount, total_amount, payment_method, payment_date, notes, pdf_url, created_at
		FROM receipts
		WHERE id = $1
	`
	receipt := &models.Receipt{}
	err := r.db.QueryRow(query, id).Scan(
		&receipt.ID,
		&receipt.LandlordID,
		&receipt.PropertyID,
		&receipt.TenantID,
		&receipt.ReceiptNumber,
		&receipt.PeriodMonth,
		&receipt.PeriodYear,
		&receipt.RentAmount,
		&receipt.ChargesAmount,
		&receipt.TotalAmount,
		&receipt.PaymentMethod,
		&receipt.PaymentDate,
		&receipt.Notes,
		&receipt.PDFURL,
		&receipt.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("receipt not found")
	}
	return receipt, err
}

func (r *ReceiptRepo) GetByLandlordID(landlordID int, query *models.ReceiptListQuery) ([]*models.Receipt, error) {
	sql := `
		SELECT id, landlord_id, property_id, tenant_id, receipt_number, period_month, period_year,
			rent_amount, charges_amount, total_amount, payment_method, payment_date, notes, pdf_url, created_at
		FROM receipts
		WHERE landlord_id = $1
	`

	args := []interface{}{landlordID}
	argCount := 1

	if query.PropertyID != nil {
		argCount++
		sql += fmt.Sprintf(" AND property_id = $%d", argCount)
		args = append(args, *query.PropertyID)
	}
	if query.TenantID != nil {
		argCount++
		sql += fmt.Sprintf(" AND tenant_id = $%d", argCount)
		args = append(args, *query.TenantID)
	}
	if query.Year != nil {
		argCount++
		sql += fmt.Sprintf(" AND period_year = $%d", argCount)
		args = append(args, *query.Year)
	}
	if query.Month != nil {
		argCount++
		sql += fmt.Sprintf(" AND period_month = $%d", argCount)
		args = append(args, *query.Month)
	}

	sql += " ORDER BY period_year DESC, period_month DESC, created_at DESC"

	if query.Limit > 0 {
		argCount++
		sql += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, query.Limit)
	}
	if query.Offset > 0 {
		argCount++
		sql += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, query.Offset)
	}

	rows, err := r.db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	receipts := []*models.Receipt{}
	for rows.Next() {
		receipt := &models.Receipt{}
		err := rows.Scan(
			&receipt.ID,
			&receipt.LandlordID,
			&receipt.PropertyID,
			&receipt.TenantID,
			&receipt.ReceiptNumber,
			&receipt.PeriodMonth,
			&receipt.PeriodYear,
			&receipt.RentAmount,
			&receipt.ChargesAmount,
			&receipt.TotalAmount,
			&receipt.PaymentMethod,
			&receipt.PaymentDate,
			&receipt.Notes,
			&receipt.PDFURL,
			&receipt.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		receipts = append(receipts, receipt)
	}
	return receipts, nil
}

func (r *ReceiptRepo) GetWithDetails(id int) (*models.ReceiptWithDetails, error) {
	query := `
		SELECT 
			r.id, r.landlord_id, r.property_id, r.tenant_id, r.receipt_number, r.period_month, r.period_year,
			r.rent_amount, r.charges_amount, r.total_amount, r.payment_method, r.payment_date, r.notes, r.pdf_url, r.created_at,
			l.email, l.first_name, l.last_name, l.address, l.city, l.postal_code, l.phone, l.signature_url,
			p.address, p.city, p.postal_code, p.rent_amount, p.charges_amount, p.syndic_name, p.syndic_address,
			t.first_name, t.last_name, t.email, t.phone, t.move_in_date
		FROM receipts r
		JOIN landlords l ON r.landlord_id = l.id
		JOIN properties p ON r.property_id = p.id
		JOIN tenants t ON r.tenant_id = t.id
		WHERE r.id = $1
	`

	result := &models.ReceiptWithDetails{
		Landlord: &models.Landlord{},
		Property: &models.Property{},
		Tenant:   &models.Tenant{},
	}

	err := r.db.QueryRow(query, id).Scan(
		&result.ID, &result.LandlordID, &result.PropertyID, &result.TenantID,
		&result.ReceiptNumber, &result.PeriodMonth, &result.PeriodYear,
		&result.RentAmount, &result.ChargesAmount, &result.TotalAmount,
		&result.PaymentMethod, &result.PaymentDate, &result.Notes, &result.PDFURL, &result.CreatedAt,
		&result.Landlord.Email, &result.Landlord.FirstName, &result.Landlord.LastName,
		&result.Landlord.Address, &result.Landlord.City, &result.Landlord.PostalCode,
		&result.Landlord.Phone, &result.Landlord.SignatureURL,
		&result.Property.Address, &result.Property.City, &result.Property.PostalCode,
		&result.Property.RentAmount, &result.Property.ChargesAmount,
		&result.Property.SyndicName, &result.Property.SyndicAddress,
		&result.Tenant.FirstName, &result.Tenant.LastName, &result.Tenant.Email,
		&result.Tenant.Phone, &result.Tenant.MoveInDate,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("receipt not found")
	}
	return result, err
}

func (r *ReceiptRepo) UpdatePDFURL(id int, pdfURL string) error {
	query := `UPDATE receipts SET pdf_url = $1 WHERE id = $2`
	_, err := r.db.Exec(query, pdfURL, id)
	return err
}

func (r *ReceiptRepo) CheckDuplicate(propertyID, tenantID, year, month int) (bool, error) {
	query := `
		SELECT COUNT(*) FROM receipts 
		WHERE property_id = $1 AND tenant_id = $2 AND period_year = $3 AND period_month = $4
	`
	var count int
	err := r.db.QueryRow(query, propertyID, tenantID, year, month).Scan(&count)
	return count > 0, err
}
