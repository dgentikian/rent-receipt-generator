package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dgentikian/rent-receipt-generator/internal/api"
	"github.com/dgentikian/rent-receipt-generator/internal/api/handlers"
	"github.com/dgentikian/rent-receipt-generator/internal/config"
	"github.com/dgentikian/rent-receipt-generator/internal/controller"
	"github.com/dgentikian/rent-receipt-generator/internal/model/database"
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

	// Initialize repositories (Model layer)
	landlordRepo := database.NewLandlordRepository(db.DB)
	propertyRepo := database.NewPropertyRepository(db.DB)
	tenantRepo := database.NewTenantRepository(db.DB)
	receiptRepo := database.NewReceiptRepository(db.DB)

	// Initialize controllers (Controller layer)
	landlordController, err := controller.NewLandlordService(landlordRepo, cfg.JWT.Secret, cfg.JWT.Expiry)
	if err != nil {
		log.Fatalf("Failed to create landlord controller: %v", err)
	}

	pdfController := controller.NewPDFService(cfg.App.UploadsDir)
	receiptController := controller.NewReceiptService(receiptRepo, propertyRepo, tenantRepo, pdfController)

	// Initialize handlers (View layer)
	landlordHandler := handlers.NewLandlordHandler(landlordController)
	propertyHandler := handlers.NewPropertyHandler(propertyRepo)
	tenantHandler := handlers.NewTenantHandler(tenantRepo, propertyRepo)
	receiptHandler := handlers.NewReceiptHandler(receiptController)

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
