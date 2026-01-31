package models

import (
	"database/sql"
	"time"
)

type Receipt struct {
	ID            int          `json:"id"`
	LandlordID    int          `json:"landlord_id"`
	PropertyID    int          `json:"property_id"`
	TenantID      int          `json:"tenant_id"`
	ReceiptNumber string       `json:"receipt_number"`
	PeriodMonth   int          `json:"period_month"`
	PeriodYear    int          `json:"period_year"`
	RentAmount    float64      `json:"rent_amount"`
	ChargesAmount float64      `json:"charges_amount"`
	TotalAmount   float64      `json:"total_amount"`
	PaymentMethod string       `json:"payment_method"`
	PaymentDate   sql.NullTime `json:"payment_date"`
	Notes         string       `json:"notes"`
	PDFURL        string       `json:"pdf_url"`
	CreatedAt     time.Time    `json:"created_at"`
}

type ReceiptWithDetails struct {
	Receipt
	Landlord *Landlord `json:"landlord"`
	Property *Property `json:"property"`
	Tenant   *Tenant   `json:"tenant"`
}

type ReceiptCreateRequest struct {
	PropertyID    int     `json:"property_id" binding:"required"`
	TenantID      int     `json:"tenant_id" binding:"required"`
	PeriodMonth   int     `json:"period_month" binding:"required,min=1,max=12"`
	PeriodYear    int     `json:"period_year" binding:"required,min=2000"`
	RentAmount    float64 `json:"rent_amount" binding:"required,gt=0"`
	ChargesAmount float64 `json:"charges_amount"`
	PaymentMethod string  `json:"payment_method"`
	PaymentDate   string  `json:"payment_date"` // Format: "2006-01-02"
	Notes         string  `json:"notes"`
}

type ReceiptListQuery struct {
	PropertyID *int `form:"property_id"`
	TenantID   *int `form:"tenant_id"`
	Year       *int `form:"year"`
	Month      *int `form:"month"`
	Limit      int  `form:"limit"`
	Offset     int  `form:"offset"`
}
