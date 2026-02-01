package handlers

import (
	"net/http"
	"strconv"

	"github.com/dgentikian/rent-receipt-generator/internal/models"
	"github.com/dgentikian/rent-receipt-generator/internal/model"
	"github.com/dgentikian/rent-receipt-generator/pkg/utils"
	"github.com/gin-gonic/gin"
)

type PropertyHandler struct {
	repo model.PropertyRepository
}

func NewPropertyHandler(repo model.PropertyRepository) *PropertyHandler {
	return &PropertyHandler{repo: repo}
}

// Create creates a new property
func (h *PropertyHandler) Create(c *gin.Context) {
	landlordID := c.GetInt("landlord_id")

	var req models.PropertyCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	property := &models.Property{
		LandlordID:    landlordID,
		Address:       req.Address,
		City:          utils.StringToPtr(req.City),
		PostalCode:    utils.StringToPtr(req.PostalCode),
		RentAmount:    req.RentAmount,
		ChargesAmount: req.ChargesAmount,
		SyndicName:    utils.StringToPtr(req.SyndicName),
		SyndicAddress: utils.StringToPtr(req.SyndicAddress),
	}

	if err := h.repo.Create(property); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create property"})
		return
	}

	c.JSON(http.StatusCreated, property)
}

// List returns all properties for the current landlord
func (h *PropertyHandler) List(c *gin.Context) {
	landlordID := c.GetInt("landlord_id")

	properties, err := h.repo.GetByLandlordID(landlordID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch properties"})
		return
	}

	c.JSON(http.StatusOK, properties)
}

// GetByID returns a specific property
func (h *PropertyHandler) GetByID(c *gin.Context) {
	landlordID := c.GetInt("landlord_id")
	propertyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid property ID"})
		return
	}

	property, err := h.repo.GetByID(propertyID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Property not found"})
		return
	}

	// Verify ownership
	if property.LandlordID != landlordID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, property)
}

// Update updates a property
func (h *PropertyHandler) Update(c *gin.Context) {
	landlordID := c.GetInt("landlord_id")
	propertyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid property ID"})
		return
	}

	property, err := h.repo.GetByID(propertyID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Property not found"})
		return
	}

	// Verify ownership
	if property.LandlordID != landlordID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.PropertyUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	property.Address = req.Address
	property.City = utils.StringToPtr(req.City)
	property.PostalCode = utils.StringToPtr(req.PostalCode)
	property.RentAmount = req.RentAmount
	property.ChargesAmount = req.ChargesAmount
	property.SyndicName = utils.StringToPtr(req.SyndicName)
	property.SyndicAddress = utils.StringToPtr(req.SyndicAddress)

	if err := h.repo.Update(property); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update property"})
		return
	}

	c.JSON(http.StatusOK, property)
}

// Delete deletes a property
func (h *PropertyHandler) Delete(c *gin.Context) {
	landlordID := c.GetInt("landlord_id")
	propertyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid property ID"})
		return
	}

	property, err := h.repo.GetByID(propertyID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Property not found"})
		return
	}

	// Verify ownership
	if property.LandlordID != landlordID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.repo.Delete(propertyID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete property"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Property deleted successfully"})
}
