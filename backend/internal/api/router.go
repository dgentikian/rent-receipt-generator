package api

import (
	"github.com/gin-gonic/gin"
	"github.com/dgentikian/rent-receipt-generator/internal/api/handlers"
	"github.com/dgentikian/rent-receipt-generator/internal/api/middleware"
	"github.com/dgentikian/rent-receipt-generator/internal/config"
)

type Router struct {
	landlordHandler *handlers.LandlordHandler
	propertyHandler *handlers.PropertyHandler
	tenantHandler   *handlers.TenantHandler
	receiptHandler  *handlers.ReceiptHandler
	config          *config.Config
}

func NewRouter(
	landlordHandler *handlers.LandlordHandler,
	propertyHandler *handlers.PropertyHandler,
	tenantHandler *handlers.TenantHandler,
	receiptHandler *handlers.ReceiptHandler,
	config *config.Config,
) *Router {
	return &Router{
		landlordHandler: landlordHandler,
		propertyHandler: propertyHandler,
		tenantHandler:   tenantHandler,
		receiptHandler:  receiptHandler,
		config:          config,
	}
}

func (r *Router) Setup() *gin.Engine {
	// Set Gin mode
	if r.config.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware(r.config.App.CORSOrigins))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", r.landlordHandler.Register)
			auth.POST("/login", r.landlordHandler.Login)
		}

		// Protected routes (authentication required)
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware(r.config.JWT.Secret))
		{
			// Landlord routes
			landlords := protected.Group("/landlord")
			{
				landlords.GET("/profile", r.landlordHandler.GetProfile)
				landlords.PUT("/profile", r.landlordHandler.UpdateProfile)
				landlords.POST("/signature", r.landlordHandler.UploadSignature)
			}

			// Property routes
			properties := protected.Group("/properties")
			{
				properties.POST("", r.propertyHandler.Create)
				properties.GET("", r.propertyHandler.List)
				properties.GET("/:id", r.propertyHandler.GetByID)
				properties.PUT("/:id", r.propertyHandler.Update)
				properties.DELETE("/:id", r.propertyHandler.Delete)
			}

			// Tenant routes
			tenants := protected.Group("/tenants")
			{
				tenants.POST("", r.tenantHandler.Create)
				tenants.GET("", r.tenantHandler.List) // ?property_id=X
				tenants.GET("/:id", r.tenantHandler.GetByID)
				tenants.PUT("/:id", r.tenantHandler.Update)
				tenants.DELETE("/:id", r.tenantHandler.Delete)
			}

			// Receipt routes
			receipts := protected.Group("/receipts")
			{
				receipts.POST("", r.receiptHandler.Create)
				receipts.GET("", r.receiptHandler.List) // Supports filters
				receipts.GET("/:id", r.receiptHandler.GetByID)
				receipts.GET("/:id/details", r.receiptHandler.GetWithDetails)
				receipts.GET("/:id/pdf", r.receiptHandler.DownloadPDF)
			}
		}
	}

	// Serve uploaded files
	router.Static("/uploads", r.config.App.UploadsDir)

	return router
}
