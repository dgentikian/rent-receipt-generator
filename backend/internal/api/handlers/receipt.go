package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/dgentikian/rent-receipt-generator/internal/models"
	"github.com/dgentikian/rent-receipt-generator/internal/service"
)

type ReceiptHandler struct {
	service *service.ReceiptService
}

func NewReceiptHandler(service *service.ReceiptService) *ReceiptHandler {
	return &ReceiptHandler{service: service}
}

// Create generates a new receipt
func (h *ReceiptHandler) Create(c *gin.Context) {
	landlordID := c.GetInt("landlord_id")

	var req models.ReceiptCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	receipt, err := h.service.Create(landlordID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, receipt)
}

// List returns all receipts for the landlord
func (h *ReceiptHandler) List(c *gin.Context) {
	landlordID := c.GetInt("landlord_id")

	var query models.ReceiptListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default limit
	if query.Limit == 0 {
		query.Limit = 50
	}

	receipts, err := h.service.List(landlordID, &query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch receipts"})
		return
	}

	c.JSON(http.StatusOK, receipts)
}

// GetByID returns a specific receipt
func (h *ReceiptHandler) GetByID(c *gin.Context) {
	landlordID := c.GetInt("landlord_id")
	receiptID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid receipt ID"})
		return
	}

	receipt, err := h.service.GetByID(receiptID, landlordID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, receipt)
}

// GetWithDetails returns a receipt with all related details
func (h *ReceiptHandler) GetWithDetails(c *gin.Context) {
	landlordID := c.GetInt("landlord_id")
	receiptID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid receipt ID"})
		return
	}

	receipt, err := h.service.GetWithDetails(receiptID, landlordID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, receipt)
}

// DownloadPDF serves the PDF file
func (h *ReceiptHandler) DownloadPDF(c *gin.Context) {
	landlordID := c.GetInt("landlord_id")
	receiptID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid receipt ID"})
		return
	}

	receipt, err := h.service.GetByID(receiptID, landlordID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if receipt.PDFURL == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "PDF not found"})
		return
	}

	c.File(receipt.PDFURL)
}
