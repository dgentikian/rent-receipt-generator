package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/dgentikian/rent-receipt-generator/internal/models"
	"github.com/dgentikian/rent-receipt-generator/internal/repository"
	"github.com/dgentikian/rent-receipt-generator/internal/repository/postgres"
)

type TenantHandler struct {
	tenantRepo   repository.TenantRepository
	propertyRepo repository.PropertyRepository
}

func NewTenantHandler(tenantRepo repository.TenantRepository, propertyRepo repository.PropertyRepository) *TenantHandler {
	return &TenantHandler{
		tenantRepo:   tenantRepo,
		propertyRepo: propertyRepo,
	}
}

// Create creates a new tenant
func (h *TenantHandler) Create(c *gin.Context) {
	landlordID := c.GetInt("landlord_id")

	var req models.TenantCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify property ownership
	property, err := h.propertyRepo.GetByID(req.PropertyID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Property not found"})
		return
	}
	if property.LandlordID != landlordID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	tenant := &models.Tenant{
		PropertyID: req.PropertyID,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Email:      req.Email,
		Phone:      req.Phone,
		MoveInDate: postgres.ParseNullTime(req.MoveInDate),
		IsActive:   true,
	}

	if err := h.tenantRepo.Create(tenant); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tenant"})
		return
	}

	c.JSON(http.StatusCreated, tenant)
}

// List returns all tenants for a property
func (h *TenantHandler) List(c *gin.Context) {
	landlordID := c.GetInt("landlord_id")
	propertyID, err := strconv.Atoi(c.Query("property_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Property ID is required"})
		return
	}

	// Verify property ownership
	property, err := h.propertyRepo.GetByID(propertyID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Property not found"})
		return
	}
	if property.LandlordID != landlordID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	tenants, err := h.tenantRepo.GetByPropertyID(propertyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tenants"})
		return
	}

	c.JSON(http.StatusOK, tenants)
}

// GetByID returns a specific tenant
func (h *TenantHandler) GetByID(c *gin.Context) {
	landlordID := c.GetInt("landlord_id")
	tenantID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	tenant, err := h.tenantRepo.GetByID(tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tenant not found"})
		return
	}

	// Verify property ownership
	property, err := h.propertyRepo.GetByID(tenant.PropertyID)
	if err != nil || property.LandlordID != landlordID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, tenant)
}

// Update updates a tenant
func (h *TenantHandler) Update(c *gin.Context) {
	landlordID := c.GetInt("landlord_id")
	tenantID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	tenant, err := h.tenantRepo.GetByID(tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tenant not found"})
		return
	}

	// Verify property ownership
	property, err := h.propertyRepo.GetByID(tenant.PropertyID)
	if err != nil || property.LandlordID != landlordID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.TenantUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenant.FirstName = req.FirstName
	tenant.LastName = req.LastName
	tenant.Email = req.Email
	tenant.Phone = req.Phone
	tenant.MoveInDate = postgres.ParseNullTime(req.MoveInDate)
	tenant.MoveOutDate = postgres.ParseNullTime(req.MoveOutDate)
	if req.IsActive != nil {
		tenant.IsActive = *req.IsActive
	}

	if err := h.tenantRepo.Update(tenant); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tenant"})
		return
	}

	c.JSON(http.StatusOK, tenant)
}

// Delete deletes a tenant
func (h *TenantHandler) Delete(c *gin.Context) {
	landlordID := c.GetInt("landlord_id")
	tenantID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	tenant, err := h.tenantRepo.GetByID(tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tenant not found"})
		return
	}

	// Verify property ownership
	property, err := h.propertyRepo.GetByID(tenant.PropertyID)
	if err != nil || property.LandlordID != landlordID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.tenantRepo.Delete(tenantID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tenant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tenant deleted successfully"})
}
