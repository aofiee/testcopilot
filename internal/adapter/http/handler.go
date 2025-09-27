package http

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"hmac-service/internal/domain"
)

// Handler holds the HTTP handlers
type Handler struct {
	hmacService domain.HMACService
}

// NewHandler creates a new HTTP handler
func NewHandler(hmacService domain.HMACService) *Handler {
	return &Handler{
		hmacService: hmacService,
	}
}

// Encrypt handles the POST /hmac/encrypt endpoint
func (h *Handler) Encrypt(c *fiber.Ctx) error {
	var req domain.HMACRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid JSON format",
		})
	}

	// Validate required fields
	if req.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{
			Error:   "validation_error",
			Message: "message is required",
		})
	}
	if req.Key == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{
			Error:   "validation_error",
			Message: "key is required",
		})
	}

	// Generate HMAC
	hmacValue, err := h.hmacService.Generate(req.Message, req.Key)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse{
			Error:   "hmac_generation_failed",
			Message: fmt.Sprintf("Failed to generate HMAC: %v", err),
		})
	}

	// Return response
	response := domain.HMACEncryptResponse{
		HMAC:      hmacValue,
		Algorithm: "HMAC-SHA256",
	}

	return c.JSON(response)
}

// Decrypt handles the POST /hmac/decrypt endpoint
func (h *Handler) Decrypt(c *fiber.Ctx) error {
	var req domain.HMACRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid JSON format",
		})
	}

	// Validate required fields
	if req.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{
			Error:   "validation_error",
			Message: "message is required",
		})
	}
	if req.Key == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{
			Error:   "validation_error",
			Message: "key is required",
		})
	}
	if req.HMAC == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{
			Error:   "validation_error",
			Message: "hmac is required",
		})
	}

	// Verify HMAC
	valid, err := h.hmacService.Verify(req.Message, req.Key, req.HMAC)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse{
			Error:   "hmac_verification_failed",
			Message: fmt.Sprintf("Failed to verify HMAC: %v", err),
		})
	}

	// Return response
	response := domain.HMACDecryptResponse{
		Valid: valid,
	}

	return c.JSON(response)
}

// HealthCheck handles the GET /health endpoint
func (h *Handler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "healthy",
		"service": "hmac-service",
	})
}