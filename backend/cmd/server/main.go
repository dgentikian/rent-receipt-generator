package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dgentikian/rent-receipt-generator/internal/api"
	"github.com/dgentikian/rent-receipt-generator/internal/api/handlers"
	"github.com/dgentikian/rent-receipt-generator/internal/config"
	"github.com/dgentikian/rent-receipt-generator/internal/database"
	"github.com/dgentikian/rent-receipt-generator/internal/repository/postgres"
	"github.com/dgentikian/rent-receipt-generator/internal/service"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
)

func main() {
	log.Printf("Starting Rent Receipt Generator API v%s (built: %s)\n", Version, BuildTime)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Environment: %s", cfg.Server.Env)

	// Connect to database
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Database connection established")

	// Ensure uploads directory exists
	if err := os.MkdirAll(cfg.App.UploadsDir, 0755); err != nil {
		log.Fatalf("Failed to create uploads directory: %v", err)
	}

	// Initialize repositories
	landlordRepo := postgres.NewLandlordRepository(db.DB)
	propertyRepo := postgres.NewPropertyRepository(db.DB)
	tenantRepo := postgres.NewTenantRepository(db.DB)
	receiptRepo := postgres.NewReceiptRepository(db.DB)

	// Initialize services
	landlordService, err := service.NewLandlordService(landlordRepo, cfg.JWT.Secret, cfg.JWT.Expiry)
	if err != nil {
		log.Fatalf("Failed to create landlord service: %v", err)
	}

	pdfService := service.NewPDFService(cfg.App.UploadsDir)
	receiptService := service.NewReceiptService(receiptRepo, propertyRepo, tenantRepo, pdfService)

	// Initialize handlers
	landlordHandler := handlers.NewLandlordHandler(landlordService)
	propertyHandler := handlers.NewPropertyHandler(propertyRepo)
	tenantHandler := handlers.NewTenantHandler(tenantRepo, propertyRepo)
	receiptHandler := handlers.NewReceiptHandler(receiptService)

	// Setup router
	router := api.NewRouter(
		landlordHandler,
		propertyHandler,
		tenantHandler,
		receiptHandler,
		cfg,
	)

	engine := router.Setup()

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Server listening on %s", addr)

	if err := engine.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
