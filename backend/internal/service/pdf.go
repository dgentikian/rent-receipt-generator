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
	pdf.AddPage()

	// Set font for title
	pdf.SetFont("Arial", "B", 20)
	pdf.CellFormat(0, 15, "QUITTANCE DE LOYER", "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Receipt info
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(0, 8, fmt.Sprintf("Quittance N° : %s", receipt.ReceiptNumber), "", 1, "", false, 0, "")
	pdf.CellFormat(0, 8, fmt.Sprintf("Date : %s", receipt.CreatedAt.Format("02/01/2006")), "", 1, "", false, 0, "")
	pdf.CellFormat(0, 8, fmt.Sprintf("Période : %s %d", getMonthName(receipt.PeriodMonth), receipt.PeriodYear), "", 1, "", false, 0, "")
	pdf.Ln(10)

	// Landlord information
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, "Bailleur", "", 1, "", false, 0, "")
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(0, 6, fmt.Sprintf("%s %s", receipt.Landlord.FirstName, receipt.Landlord.LastName), "", 1, "", false, 0, "")
	if receipt.Landlord.Address != "" {
		pdf.CellFormat(0, 6, receipt.Landlord.Address, "", 1, "", false, 0, "")
	}
	if receipt.Landlord.City != "" && receipt.Landlord.PostalCode != "" {
		pdf.CellFormat(0, 6, fmt.Sprintf("%s %s", receipt.Landlord.PostalCode, receipt.Landlord.City), "", 1, "", false, 0, "")
	}
	if receipt.Landlord.Phone != "" {
		pdf.CellFormat(0, 6, fmt.Sprintf("Tél : %s", receipt.Landlord.Phone), "", 1, "", false, 0, "")
	}
	pdf.Ln(8)

	// Tenant information
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, "Locataire", "", 1, "", false, 0, "")
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(0, 6, fmt.Sprintf("%s %s", receipt.Tenant.FirstName, receipt.Tenant.LastName), "", 1, "", false, 0, "")
	if receipt.Tenant.Email != "" {
		pdf.CellFormat(0, 6, fmt.Sprintf("Email : %s", receipt.Tenant.Email), "", 1, "", false, 0, "")
	}
	if receipt.Tenant.Phone != "" {
		pdf.CellFormat(0, 6, fmt.Sprintf("Tél : %s", receipt.Tenant.Phone), "", 1, "", false, 0, "")
	}
	pdf.Ln(8)

	// Property information
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, "Bien loué", "", 1, "", false, 0, "")
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(0, 6, receipt.Property.Address, "", 1, "", false, 0, "")
	if receipt.Property.City != "" && receipt.Property.PostalCode != "" {
		pdf.CellFormat(0, 6, fmt.Sprintf("%s %s", receipt.Property.PostalCode, receipt.Property.City), "", 1, "", false, 0, "")
	}
	pdf.Ln(10)

	// Payment details
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, "Détail du paiement", "", 1, "", false, 0, "")
	pdf.SetFont("Arial", "", 11)

	// Draw table
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(120, 8, "Loyer", "1", 0, "", true, 0, "")
	pdf.CellFormat(60, 8, fmt.Sprintf("%.2f €", receipt.RentAmount), "1", 1, "R", true, 0, "")

	if receipt.ChargesAmount > 0 {
		pdf.CellFormat(120, 8, "Charges", "1", 0, "", false, 0, "")
		pdf.CellFormat(60, 8, fmt.Sprintf("%.2f €", receipt.ChargesAmount), "1", 1, "R", false, 0, "")
	}

	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(120, 8, "Total", "1", 0, "", true, 0, "")
	pdf.CellFormat(60, 8, fmt.Sprintf("%.2f €", receipt.TotalAmount), "1", 1, "R", true, 0, "")
	pdf.Ln(8)

	// Payment method and date
	if receipt.PaymentMethod != "" {
		pdf.SetFont("Arial", "", 11)
		pdf.CellFormat(0, 6, fmt.Sprintf("Moyen de paiement : %s", receipt.PaymentMethod), "", 1, "", false, 0, "")
	}
	if receipt.PaymentDate.Valid {
		pdf.CellFormat(0, 6, fmt.Sprintf("Date de paiement : %s", receipt.PaymentDate.Time.Format("02/01/2006")), "", 1, "", false, 0, "")
	}
	pdf.Ln(10)

	// Legal text
	pdf.SetFont("Arial", "I", 10)
	pdf.MultiCell(0, 5, "Je soussigné(e) "+receipt.Landlord.FirstName+" "+receipt.Landlord.LastName+", propriétaire du logement ci-dessus désigné, reconnais avoir reçu de "+receipt.Tenant.FirstName+" "+receipt.Tenant.LastName+" la somme indiquée ci-dessus, au titre du loyer et des charges pour la période mentionnée.", "", "", false)
	pdf.Ln(15)

	// Signature
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(0, 6, fmt.Sprintf("Fait le %s", time.Now().Format("02/01/2006")), "", 1, "R", false, 0, "")
	pdf.Ln(5)

	// Add signature image if available
	if receipt.Landlord.SignatureURL != "" && fileExists(receipt.Landlord.SignatureURL) {
		x := pdf.GetX() + 120
		y := pdf.GetY()
		pdf.Image(receipt.Landlord.SignatureURL, x, y, 60, 0, false, "", 0, "")
		pdf.Ln(20)
	} else {
		pdf.Ln(15)
	}

	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(0, 6, "Signature du bailleur", "", 1, "R", false, 0, "")

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
