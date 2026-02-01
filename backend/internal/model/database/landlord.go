package database

import (
	"database/sql"
	"fmt"

	"github.com/dgentikian/rent-receipt-generator/internal/models"
)

type LandlordRepo struct {
	db *sql.DB
}

func NewLandlordRepository(db *sql.DB) *LandlordRepo {
	return &LandlordRepo{db: db}
}

func (r *LandlordRepo) Create(landlord *models.Landlord) error {
	query := `
		INSERT INTO public.landlords (email, password_hash, first_name, last_name, address, city, postal_code, phone)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		landlord.Email,
		landlord.PasswordHash,
		landlord.FirstName,
		landlord.LastName,
		landlord.Address,
		landlord.City,
		landlord.PostalCode,
		landlord.Phone,
	).Scan(&landlord.ID, &landlord.CreatedAt, &landlord.UpdatedAt)
}

func (r *LandlordRepo) GetByID(id int) (*models.Landlord, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, address, city, postal_code, phone, signature_url, created_at, updated_at
		FROM public.landlords
		WHERE id = $1
	`
	landlord := &models.Landlord{}
	err := r.db.QueryRow(query, id).Scan(
		&landlord.ID,
		&landlord.Email,
		&landlord.PasswordHash,
		&landlord.FirstName,
		&landlord.LastName,
		&landlord.Address,
		&landlord.City,
		&landlord.PostalCode,
		&landlord.Phone,
		&landlord.SignatureURL,
		&landlord.CreatedAt,
		&landlord.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("landlord not found")
	}
	return landlord, err
}

func (r *LandlordRepo) GetByEmail(email string) (*models.Landlord, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, address, city, postal_code, phone, signature_url, created_at, updated_at
		FROM public.landlords
		WHERE email = $1
	`
	landlord := &models.Landlord{}
	err := r.db.QueryRow(query, email).Scan(
		&landlord.ID,
		&landlord.Email,
		&landlord.PasswordHash,
		&landlord.FirstName,
		&landlord.LastName,
		&landlord.Address,
		&landlord.City,
		&landlord.PostalCode,
		&landlord.Phone,
		&landlord.SignatureURL,
		&landlord.CreatedAt,
		&landlord.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("landlord not found")
	}
	if err != nil {
		return nil, err
	}
	return landlord, nil
}

func (r *LandlordRepo) Update(landlord *models.Landlord) error {
	query := `
		UPDATE public.landlords
		SET first_name = $1, last_name = $2, address = $3, city = $4, postal_code = $5, phone = $6
		WHERE id = $7
		RETURNING updated_at
	`
	return r.db.QueryRow(
		query,
		landlord.FirstName,
		landlord.LastName,
		landlord.Address,
		landlord.City,
		landlord.PostalCode,
		landlord.Phone,
		landlord.ID,
	).Scan(&landlord.UpdatedAt)
}

func (r *LandlordRepo) UpdateSignature(id int, signatureURL string) error {
	query := `UPDATE public.landlords SET signature_url = $1 WHERE id = $2`
	_, err := r.db.Exec(query, signatureURL, id)
	return err
}
