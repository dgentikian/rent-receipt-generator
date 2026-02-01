package database

import (
	"database/sql"
	"fmt"

	"github.com/dgentikian/rent-receipt-generator/internal/models"
)

type PropertyRepo struct {
	db *sql.DB
}

func NewPropertyRepository(db *sql.DB) *PropertyRepo {
	return &PropertyRepo{db: db}
}

func (r *PropertyRepo) Create(property *models.Property) error {
	query := `
		INSERT INTO properties (landlord_id, address, city, postal_code, rent_amount, charges_amount, syndic_name, syndic_address)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		property.LandlordID,
		property.Address,
		property.City,
		property.PostalCode,
		property.RentAmount,
		property.ChargesAmount,
		property.SyndicName,
		property.SyndicAddress,
	).Scan(&property.ID, &property.CreatedAt, &property.UpdatedAt)
}

func (r *PropertyRepo) GetByID(id int) (*models.Property, error) {
	query := `
		SELECT id, landlord_id, address, city, postal_code, rent_amount, charges_amount, syndic_name, syndic_address, created_at, updated_at
		FROM properties
		WHERE id = $1
	`
	property := &models.Property{}
	err := r.db.QueryRow(query, id).Scan(
		&property.ID,
		&property.LandlordID,
		&property.Address,
		&property.City,
		&property.PostalCode,
		&property.RentAmount,
		&property.ChargesAmount,
		&property.SyndicName,
		&property.SyndicAddress,
		&property.CreatedAt,
		&property.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("property not found")
	}
	return property, err
}

func (r *PropertyRepo) GetByLandlordID(landlordID int) ([]*models.Property, error) {
	query := `
		SELECT id, landlord_id, address, city, postal_code, rent_amount, charges_amount, syndic_name, syndic_address, created_at, updated_at
		FROM properties
		WHERE landlord_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, landlordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	properties := []*models.Property{}
	for rows.Next() {
		property := &models.Property{}
		err := rows.Scan(
			&property.ID,
			&property.LandlordID,
			&property.Address,
			&property.City,
			&property.PostalCode,
			&property.RentAmount,
			&property.ChargesAmount,
			&property.SyndicName,
			&property.SyndicAddress,
			&property.CreatedAt,
			&property.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		properties = append(properties, property)
	}
	return properties, nil
}

func (r *PropertyRepo) Update(property *models.Property) error {
	query := `
		UPDATE properties
		SET address = $1, city = $2, postal_code = $3, rent_amount = $4, charges_amount = $5, syndic_name = $6, syndic_address = $7
		WHERE id = $8
		RETURNING updated_at
	`
	return r.db.QueryRow(
		query,
		property.Address,
		property.City,
		property.PostalCode,
		property.RentAmount,
		property.ChargesAmount,
		property.SyndicName,
		property.SyndicAddress,
		property.ID,
	).Scan(&property.UpdatedAt)
}

func (r *PropertyRepo) Delete(id int) error {
	query := `DELETE FROM properties WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("property not found")
	}
	return nil
}
