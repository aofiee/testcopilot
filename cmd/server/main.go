package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"hmac-service/internal/adapter/config"
	"hmac-service/internal/adapter/http"
	"hmac-service/internal/usecase"
)

func main() {
	// Load configuration
	cfg := config.NewConfig()

	// Initialize services (dependency injection)
	hmacService := usecase.NewHMACService()

	// Initialize HTTP handler
	handler := http.NewHandler(hmacService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "HMAC Service",
		ServerHeader: "hmac-service",
	})

	// Setup routes
	http.SetupRoutes(app, handler)

	// Start server
	log.Printf("Starting server on port %s", cfg.Port)
	log.Printf("Health check available at: http://localhost:%s/health", cfg.Port)
	log.Printf("HMAC encrypt endpoint: http://localhost:%s/hmac/encrypt", cfg.Port)
	log.Printf("HMAC decrypt endpoint: http://localhost:%s/hmac/decrypt", cfg.Port)

	if err := app.Listen(cfg.GetPort()); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}