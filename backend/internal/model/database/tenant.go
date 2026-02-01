package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dgentikian/rent-receipt-generator/internal/models"
)

type TenantRepo struct {
	db *sql.DB
}

func NewTenantRepository(db *sql.DB) *TenantRepo {
	return &TenantRepo{db: db}
}

func (r *TenantRepo) Create(tenant *models.Tenant) error {
	query := `
		INSERT INTO tenants (property_id, first_name, last_name, email, phone, move_in_date, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		tenant.PropertyID,
		tenant.FirstName,
		tenant.LastName,
		tenant.Email,
		tenant.Phone,
		tenant.MoveInDate,
		tenant.IsActive,
	).Scan(&tenant.ID, &tenant.CreatedAt, &tenant.UpdatedAt)
}

func (r *TenantRepo) GetByID(id int) (*models.Tenant, error) {
	query := `
		SELECT id, property_id, first_name, last_name, email, phone, move_in_date, move_out_date, is_active, created_at, updated_at
		FROM tenants
		WHERE id = $1
	`
	tenant := &models.Tenant{}
	err := r.db.QueryRow(query, id).Scan(
		&tenant.ID,
		&tenant.PropertyID,
		&tenant.FirstName,
		&tenant.LastName,
		&tenant.Email,
		&tenant.Phone,
		&tenant.MoveInDate,
		&tenant.MoveOutDate,
		&tenant.IsActive,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("tenant not found")
	}
	return tenant, err
}

func (r *TenantRepo) GetByPropertyID(propertyID int) ([]*models.Tenant, error) {
	query := `
		SELECT id, property_id, first_name, last_name, email, phone, move_in_date, move_out_date, is_active, created_at, updated_at
		FROM tenants
		WHERE property_id = $1
		ORDER BY is_active DESC, created_at DESC
	`
	rows, err := r.db.Query(query, propertyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tenants := []*models.Tenant{}
	for rows.Next() {
		tenant := &models.Tenant{}
		err := rows.Scan(
			&tenant.ID,
			&tenant.PropertyID,
			&tenant.FirstName,
			&tenant.LastName,
			&tenant.Email,
			&tenant.Phone,
			&tenant.MoveInDate,
			&tenant.MoveOutDate,
			&tenant.IsActive,
			&tenant.CreatedAt,
			&tenant.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tenants = append(tenants, tenant)
	}
	return tenants, nil
}

func (r *TenantRepo) GetByLandlordID(landlordID int) ([]*models.Tenant, error) {
	query := `
		SELECT t.id, t.property_id, t.first_name, t.last_name, t.email, t.phone, 
		       t.move_in_date, t.move_out_date, t.is_active, t.created_at, t.updated_at
		FROM tenants t
		INNER JOIN properties p ON t.property_id = p.id
		WHERE p.landlord_id = $1
		ORDER BY t.is_active DESC, t.created_at DESC
	`
	rows, err := r.db.Query(query, landlordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tenants := []*models.Tenant{}
	for rows.Next() {
		tenant := &models.Tenant{}
		err := rows.Scan(
			&tenant.ID,
			&tenant.PropertyID,
			&tenant.FirstName,
			&tenant.LastName,
			&tenant.Email,
			&tenant.Phone,
			&tenant.MoveInDate,
			&tenant.MoveOutDate,
			&tenant.IsActive,
			&tenant.CreatedAt,
			&tenant.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tenants = append(tenants, tenant)
	}
	return tenants, nil
}

func (r *TenantRepo) Update(tenant *models.Tenant) error {
	query := `
		UPDATE tenants
		SET first_name = $1, last_name = $2, email = $3, phone = $4, move_in_date = $5, move_out_date = $6, is_active = $7
		WHERE id = $8
		RETURNING updated_at
	`
	return r.db.QueryRow(
		query,
		tenant.FirstName,
		tenant.LastName,
		tenant.Email,
		tenant.Phone,
		tenant.MoveInDate,
		tenant.MoveOutDate,
		tenant.IsActive,
		tenant.ID,
	).Scan(&tenant.UpdatedAt)
}

func (r *TenantRepo) Delete(id int) error {
	query := `DELETE FROM tenants WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("tenant not found")
	}
	return nil
}

// Helper to parse date string to sql.NullTime
func ParseNullTime(dateStr string) sql.NullTime {
	if dateStr == "" {
		return sql.NullTime{Valid: false}
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: t, Valid: true}
}
