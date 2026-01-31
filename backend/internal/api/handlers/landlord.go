package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/dgentikian/rent-receipt-generator/internal/models"
	"github.com/dgentikian/rent-receipt-generator/internal/service"
)

type LandlordHandler struct {
	service *service.LandlordService
}

func NewLandlordHandler(service *service.LandlordService) *LandlordHandler {
	return &LandlordHandler{service: service}
}

// Register creates a new landlord account
func (h *LandlordHandler) Register(c *gin.Context) {
	var req models.LandlordCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	landlord, err := h.service.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, landlord)
}

// Login authenticates a landlord and returns JWT token
func (h *LandlordHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetProfile returns the current landlord's profile
func (h *LandlordHandler) GetProfile(c *gin.Context) {
	landlordID := c.GetInt("landlord_id")

	landlord, err := h.service.GetByID(landlordID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Landlord not found"})
		return
	}

	c.JSON(http.StatusOK, landlord)
}

// UpdateProfile updates the current landlord's profile
func (h *LandlordHandler) UpdateProfile(c *gin.Context) {
	landlordID := c.GetInt("landlord_id")

	var req models.LandlordUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	landlord, err := h.service.Update(landlordID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, landlord)
}

// UploadSignature handles signature image upload
func (h *LandlordHandler) UploadSignature(c *gin.Context) {
	landlordID := c.GetInt("landlord_id")

	file, err := c.FormFile("signature")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Signature file is required"})
		return
	}

	// Save file
	filename := "signature_" + strconv.Itoa(landlordID) + "_" + file.Filename
	filepath := "./uploads/" + filename

	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save signature"})
		return
	}

	// Update database
	if err := h.service.UpdateSignature(landlordID, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update signature"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"signature_url": filepath})
}
