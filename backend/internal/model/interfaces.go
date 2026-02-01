package model

import "github.com/dgentikian/rent-receipt-generator/internal/models"

type LandlordRepository interface {
	Create(landlord *models.Landlord) error
	GetByID(id int) (*models.Landlord, error)
	GetByEmail(email string) (*models.Landlord, error)
	Update(landlord *models.Landlord) error
	UpdateSignature(id int, signatureURL string) error
}

type PropertyRepository interface {
	Create(property *models.Property) error
	GetByID(id int) (*models.Property, error)
	GetByLandlordID(landlordID int) ([]*models.Property, error)
	Update(property *models.Property) error
	Delete(id int) error
}

type TenantRepository interface {
	Create(tenant *models.Tenant) error
	GetByID(id int) (*models.Tenant, error)
	GetByPropertyID(propertyID int) ([]*models.Tenant, error)
	GetByLandlordID(landlordID int) ([]*models.Tenant, error)
	Update(tenant *models.Tenant) error
	Delete(id int) error
}

type ReceiptRepository interface {
	Create(receipt *models.Receipt) error
	GetByID(id int) (*models.Receipt, error)
	GetByLandlordID(landlordID int, query *models.ReceiptListQuery) ([]*models.Receipt, error)
	GetWithDetails(id int) (*models.ReceiptWithDetails, error)
	UpdatePDFURL(id int, pdfURL string) error
	CheckDuplicate(propertyID, tenantID, year, month int) (bool, error)
}
