package models

import "time"

type Property struct {
	ID            int       `json:"id"`
	LandlordID    int       `json:"landlord_id"`
	Address       string    `json:"address"`
	City          string    `json:"city"`
	PostalCode    string    `json:"postal_code"`
	RentAmount    float64   `json:"rent_amount"`
	ChargesAmount float64   `json:"charges_amount"`
	SyndicName    string    `json:"syndic_name"`
	SyndicAddress string    `json:"syndic_address"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type PropertyCreateRequest struct {
	Address       string  `json:"address" binding:"required"`
	City          string  `json:"city"`
	PostalCode    string  `json:"postal_code"`
	RentAmount    float64 `json:"rent_amount" binding:"required,gt=0"`
	ChargesAmount float64 `json:"charges_amount"`
	SyndicName    string  `json:"syndic_name"`
	SyndicAddress string  `json:"syndic_address"`
}

type PropertyUpdateRequest struct {
	Address       string  `json:"address"`
	City          string  `json:"city"`
	PostalCode    string  `json:"postal_code"`
	RentAmount    float64 `json:"rent_amount" binding:"gt=0"`
	ChargesAmount float64 `json:"charges_amount"`
	SyndicName    string  `json:"syndic_name"`
	SyndicAddress string  `json:"syndic_address"`
}
