package service

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dgentikian/rent-receipt-generator/internal/models"
	"github.com/dgentikian/rent-receipt-generator/internal/repository"
)

type ReceiptService struct {
	receiptRepo  repository.ReceiptRepository
	propertyRepo repository.PropertyRepository
	tenantRepo   repository.TenantRepository
	pdfService   *PDFService
}

func NewReceiptService(
	receiptRepo repository.ReceiptRepository,
	propertyRepo repository.PropertyRepository,
	tenantRepo repository.TenantRepository,
	pdfService *PDFService,
) *ReceiptService {
	return &ReceiptService{
		receiptRepo:  receiptRepo,
		propertyRepo: propertyRepo,
		tenantRepo:   tenantRepo,
		pdfService:   pdfService,
	}
}

func (s *ReceiptService) Create(landlordID int, req *models.ReceiptCreateRequest) (*models.Receipt, error) {
	// Validate property belongs to landlord
	property, err := s.propertyRepo.GetByID(req.PropertyID)
	if err != nil {
		return nil, fmt.Errorf("property not found")
	}
	if property.LandlordID != landlordID {
		return nil, fmt.Errorf("unauthorized: property does not belong to landlord")
	}

	// Validate tenant belongs to property
	tenant, err := s.tenantRepo.GetByID(req.TenantID)
	if err != nil {
		return nil, fmt.Errorf("tenant not found")
	}
	if tenant.PropertyID != req.PropertyID {
		return nil, fmt.Errorf("tenant does not belong to this property")
	}

	// Check for duplicate receipt
	isDuplicate, err := s.receiptRepo.CheckDuplicate(req.PropertyID, req.TenantID, req.PeriodYear, req.PeriodMonth)
	if err != nil {
		return nil, fmt.Errorf("failed to check duplicate: %w", err)
	}
	if isDuplicate {
		return nil, fmt.Errorf("receipt already exists for this period")
	}

	// Generate receipt number
	receiptNumber := fmt.Sprintf("QT-%d-%04d%02d-%d", landlordID, req.PeriodYear, req.PeriodMonth, time.Now().Unix()%100000)

	// Calculate total
	totalAmount := req.RentAmount + req.ChargesAmount

	// Parse payment date
	var paymentDate time.Time
	if req.PaymentDate != "" {
		paymentDate, err = time.Parse("2006-01-02", req.PaymentDate)
		if err != nil {
			return nil, fmt.Errorf("invalid payment date format")
		}
	}

	receipt := &models.Receipt{
		LandlordID:    landlordID,
		PropertyID:    req.PropertyID,
		TenantID:      req.TenantID,
		ReceiptNumber: receiptNumber,
		PeriodMonth:   req.PeriodMonth,
		PeriodYear:    req.PeriodYear,
		RentAmount:    req.RentAmount,
		ChargesAmount: req.ChargesAmount,
		TotalAmount:   totalAmount,
		PaymentMethod: req.PaymentMethod,
		Notes:         req.Notes,
	}

	if req.PaymentDate != "" {
		receipt.PaymentDate = sql.NullTime{Time: paymentDate, Valid: true}
	}

	// Create receipt in database
	if err := s.receiptRepo.Create(receipt); err != nil {
		return nil, fmt.Errorf("failed to create receipt: %w", err)
	}

	// Generate PDF
	receiptWithDetails, err := s.receiptRepo.GetWithDetails(receipt.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get receipt details: %w", err)
	}

	pdfPath, err := s.pdfService.GenerateReceipt(receiptWithDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	// Update receipt with PDF URL
	if err := s.receiptRepo.UpdatePDFURL(receipt.ID, pdfPath); err != nil {
		return nil, fmt.Errorf("failed to update PDF URL: %w", err)
	}

	receipt.PDFURL = pdfPath

	return receipt, nil
}

func (s *ReceiptService) GetByID(id int, landlordID int) (*models.Receipt, error) {
	receipt, err := s.receiptRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if receipt.LandlordID != landlordID {
		return nil, fmt.Errorf("unauthorized")
	}

	return receipt, nil
}

func (s *ReceiptService) GetWithDetails(id int, landlordID int) (*models.ReceiptWithDetails, error) {
	receipt, err := s.receiptRepo.GetWithDetails(id)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if receipt.LandlordID != landlordID {
		return nil, fmt.Errorf("unauthorized")
	}

	return receipt, nil
}

func (s *ReceiptService) List(landlordID int, query *models.ReceiptListQuery) ([]*models.Receipt, error) {
	return s.receiptRepo.GetByLandlordID(landlordID, query)
}
