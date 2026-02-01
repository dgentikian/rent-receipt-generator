package models

import (
	"database/sql"
	"time"
)

type Tenant struct {
	ID          int          `json:"id"`
	PropertyID  int          `json:"property_id"`
	FirstName   string       `json:"first_name"`
	LastName    string       `json:"last_name"`
	Email       *string      `json:"email,omitempty"`
	Phone       *string      `json:"phone,omitempty"`
	MoveInDate  sql.NullTime `json:"move_in_date"`
	MoveOutDate sql.NullTime `json:"move_out_date"`
	IsActive    bool         `json:"is_active"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type TenantCreateRequest struct {
	PropertyID int    `json:"property_id" binding:"required"`
	FirstName  string `json:"first_name" binding:"required"`
	LastName   string `json:"last_name" binding:"required"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	MoveInDate string `json:"move_in_date"` // Format: "2006-01-02"
}

type TenantUpdateRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	MoveInDate  string `json:"move_in_date"`
	MoveOutDate string `json:"move_out_date"`
	IsActive    *bool  `json:"is_active"`
}
