package service

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dgentikian/rent-receipt-generator/internal/models"
	"github.com/jung-kurt/gofpdf"
)

type PDFService struct {
	uploadsDir string
}

func NewPDFService(uploadsDir string) *PDFService {
	return &PDFService{uploadsDir: uploadsDir}
}

// GenerateReceipt creates a PDF receipt and returns the file path
func (s *PDFService) GenerateReceipt(receipt *models.ReceiptWithDetails) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Add UTF-8 font support (using built-in translator)
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	pdf.AddPage()

	// Set font for title
	pdf.SetFont("Arial", "B", 20)
	pdf.CellFormat(0, 15, tr("QUITTANCE DE LOYER"), "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Receipt info
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(0, 8, tr(fmt.Sprintf("Quittance N° : %s", receipt.ReceiptNumber)), "", 1, "", false, 0, "")
	pdf.CellFormat(0, 8, tr(fmt.Sprintf("Date : %s", receipt.CreatedAt.Format("02/01/2006"))), "", 1, "", false, 0, "")
	pdf.CellFormat(0, 8, tr(fmt.Sprintf("Période : %s %d", getMonthName(receipt.PeriodMonth), receipt.PeriodYear)), "", 1, "", false, 0, "")
	pdf.Ln(10)

	// Landlord information
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, tr("Bailleur"), "", 1, "", false, 0, "")
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(0, 6, tr(fmt.Sprintf("%s %s", receipt.Landlord.FirstName, receipt.Landlord.LastName)), "", 1, "", false, 0, "")
	if receipt.Landlord.Address != nil && *receipt.Landlord.Address != "" {
		pdf.CellFormat(0, 6, tr(*receipt.Landlord.Address), "", 1, "", false, 0, "")
	}
	if receipt.Landlord.City != nil && receipt.Landlord.PostalCode != nil && *receipt.Landlord.City != "" && *receipt.Landlord.PostalCode != "" {
		pdf.CellFormat(0, 6, tr(fmt.Sprintf("%s %s", *receipt.Landlord.PostalCode, *receipt.Landlord.City)), "", 1, "", false, 0, "")
	}
	if receipt.Landlord.Phone != nil && *receipt.Landlord.Phone != "" {
		pdf.CellFormat(0, 6, tr(fmt.Sprintf("Tél : %s", *receipt.Landlord.Phone)), "", 1, "", false, 0, "")
	}
	pdf.Ln(8)

	// Tenant information
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, tr("Locataire"), "", 1, "", false, 0, "")
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(0, 6, tr(fmt.Sprintf("%s %s", receipt.Tenant.FirstName, receipt.Tenant.LastName)), "", 1, "", false, 0, "")
	if receipt.Tenant.Email != nil && *receipt.Tenant.Email != "" {
		pdf.CellFormat(0, 6, tr(fmt.Sprintf("Email : %s", *receipt.Tenant.Email)), "", 1, "", false, 0, "")
	}
	if receipt.Tenant.Phone != nil && *receipt.Tenant.Phone != "" {
		pdf.CellFormat(0, 6, tr(fmt.Sprintf("Tél : %s", *receipt.Tenant.Phone)), "", 1, "", false, 0, "")
	}
	pdf.Ln(8)

	// Property information
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, tr("Bien loué"), "", 1, "", false, 0, "")
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(0, 6, tr(receipt.Property.Address), "", 1, "", false, 0, "")
	if receipt.Property.City != nil && receipt.Property.PostalCode != nil && *receipt.Property.City != "" && *receipt.Property.PostalCode != "" {
		pdf.CellFormat(0, 6, tr(fmt.Sprintf("%s %s", *receipt.Property.PostalCode, *receipt.Property.City)), "", 1, "", false, 0, "")
	}
	pdf.Ln(10)

	// Payment details
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, tr("Détail du paiement"), "", 1, "", false, 0, "")
	pdf.SetFont("Arial", "", 11)

	// Draw table
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(120, 8, tr("Loyer"), "1", 0, "", true, 0, "")
	pdf.CellFormat(60, 8, tr(fmt.Sprintf("%.2f €", receipt.RentAmount)), "1", 1, "R", true, 0, "")

	if receipt.ChargesAmount > 0 {
		pdf.CellFormat(120, 8, tr("Charges"), "1", 0, "", false, 0, "")
		pdf.CellFormat(60, 8, tr(fmt.Sprintf("%.2f €", receipt.ChargesAmount)), "1", 1, "R", false, 0, "")
	}

	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(120, 8, tr("Total"), "1", 0, "", true, 0, "")
	pdf.CellFormat(60, 8, tr(fmt.Sprintf("%.2f €", receipt.TotalAmount)), "1", 1, "R", true, 0, "")
	pdf.Ln(8)

	// Payment method and date
	if receipt.PaymentMethod != "" {
		pdf.SetFont("Arial", "", 11)
		pdf.CellFormat(0, 6, tr(fmt.Sprintf("Moyen de paiement : %s", receipt.PaymentMethod)), "", 1, "", false, 0, "")
	}
	if receipt.PaymentDate.Valid {
		pdf.CellFormat(0, 6, tr(fmt.Sprintf("Date de paiement : %s", receipt.PaymentDate.Time.Format("02/01/2006"))), "", 1, "", false, 0, "")
	}
	pdf.Ln(10)

	// Legal text
	pdf.SetFont("Arial", "I", 10)
	legalText := fmt.Sprintf("Je soussigné(e) %s %s, propriétaire du logement ci-dessus désigné, reconnais avoir reçu de %s %s la somme indiquée ci-dessus, au titre du loyer et des charges pour la période mentionnée.",
		receipt.Landlord.FirstName, receipt.Landlord.LastName, receipt.Tenant.FirstName, receipt.Tenant.LastName)
	pdf.MultiCell(0, 5, tr(legalText), "", "", false)
	pdf.Ln(15)

	// Signature
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(0, 6, tr(fmt.Sprintf("Fait le %s", time.Now().Format("02/01/2006"))), "", 1, "R", false, 0, "")
	pdf.Ln(5)

	// Add signature image if available
	if receipt.Landlord.SignatureURL != nil && *receipt.Landlord.SignatureURL != "" && fileExists(*receipt.Landlord.SignatureURL) {
		x := pdf.GetX() + 120
		y := pdf.GetY()
		pdf.Image(*receipt.Landlord.SignatureURL, x, y, 60, 0, false, "", 0, "")
		pdf.Ln(20)
	} else {
		pdf.Ln(15)
	}

	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(0, 6, tr("Signature du bailleur"), "", 1, "R", false, 0, "")

	// Generate filename and save
	filename := fmt.Sprintf("receipt_%s_%d.pdf", receipt.ReceiptNumber, time.Now().Unix())
	filepath := filepath.Join(s.uploadsDir, filename)

	// Ensure uploads directory exists
	if err := os.MkdirAll(s.uploadsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create uploads directory: %w", err)
	}

	if err := pdf.OutputFileAndClose(filepath); err != nil {
		return "", fmt.Errorf("failed to save PDF: %w", err)
	}

	return filepath, nil
}

func getMonthName(month int) string {
	months := []string{
		"Janvier", "Février", "Mars", "Avril", "Mai", "Juin",
		"Juillet", "Août", "Septembre", "Octobre", "Novembre", "Décembre",
	}
	if month < 1 || month > 12 {
		return ""
	}
	return months[month-1]
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
