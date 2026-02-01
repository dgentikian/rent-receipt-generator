package models

import "time"

type Landlord struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never send password hash to client
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Address      *string   `json:"address,omitempty"`
	City         *string   `json:"city,omitempty"`
	PostalCode   *string   `json:"postal_code,omitempty"`
	Phone        *string   `json:"phone,omitempty"`
	SignatureURL *string   `json:"signature_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type LandlordCreateRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
	FirstName  string `json:"first_name" binding:"required"`
	LastName   string `json:"last_name" binding:"required"`
	Address    string `json:"address"`
	City       string `json:"city"`
	PostalCode string `json:"postal_code"`
	Phone      string `json:"phone"`
}

type LandlordUpdateRequest struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Address    string `json:"address"`
	City       string `json:"city"`
	PostalCode string `json:"postal_code"`
	Phone      string `json:"phone"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token    string    `json:"token"`
	Landlord *Landlord `json:"landlord"`
}
