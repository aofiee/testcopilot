package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"hmac-service/internal/domain"
	"hmac-service/internal/usecase"
)

func setupTestApp() (*fiber.App, *Handler) {
	hmacService := usecase.NewHMACService()
	handler := NewHandler(hmacService)
	
	app := fiber.New()
	SetupRoutes(app, handler)
	
	return app, handler
}

func TestHandler_Encrypt(t *testing.T) {
	app, _ := setupTestApp()

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectError    bool
	}{
		{
			name: "valid request",
			requestBody: domain.HMACRequest{
				Message: "hello world",
				Key:     "secret",
			},
			expectedStatus: 200,
			expectError:    false,
		},
		{
			name: "missing message",
			requestBody: domain.HMACRequest{
				Key: "secret",
			},
			expectedStatus: 400,
			expectError:    true,
		},
		{
			name: "missing key",
			requestBody: domain.HMACRequest{
				Message: "hello world",
			},
			expectedStatus: 400,
			expectError:    true,
		},
		{
			name:           "invalid json",
			requestBody:    "invalid json",
			expectedStatus: 400,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			var err error
			
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			}

			req := httptest.NewRequest("POST", "/hmac/encrypt", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			responseBody, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}

			if tt.expectError {
				var errorResp domain.ErrorResponse
				if err := json.Unmarshal(responseBody, &errorResp); err != nil {
					t.Errorf("Failed to unmarshal error response: %v", err)
				}
				if errorResp.Error == "" {
					t.Errorf("Expected error in response, but got none")
				}
			} else {
				var successResp domain.HMACEncryptResponse
				if err := json.Unmarshal(responseBody, &successResp); err != nil {
					t.Errorf("Failed to unmarshal success response: %v", err)
				}
				if successResp.HMAC == "" {
					t.Errorf("Expected HMAC in response, but got empty string")
				}
				if successResp.Algorithm != "HMAC-SHA256" {
					t.Errorf("Expected algorithm 'HMAC-SHA256', got '%s'", successResp.Algorithm)
				}
			}
		})
	}
}

func TestHandler_Decrypt(t *testing.T) {
	app, _ := setupTestApp()

	// Generate a valid HMAC first
	hmacService := usecase.NewHMACService()
	message := "hello world"
	key := "secret"
	validHMAC, err := hmacService.Generate(message, key)
	if err != nil {
		t.Fatalf("Failed to generate valid HMAC: %v", err)
	}

	tests := []struct {
		name           string
		requestBody    domain.HMACRequest
		expectedStatus int
		expectedValid  bool
		expectError    bool
	}{
		{
			name: "valid hmac",
			requestBody: domain.HMACRequest{
				Message: message,
				Key:     key,
				HMAC:    validHMAC,
			},
			expectedStatus: 200,
			expectedValid:  true,
			expectError:    false,
		},
		{
			name: "invalid hmac",
			requestBody: domain.HMACRequest{
				Message: message,
				Key:     key,
				HMAC:    "invalid-hmac",
			},
			expectedStatus: 200,
			expectedValid:  false,
			expectError:    false,
		},
		{
			name: "missing message",
			requestBody: domain.HMACRequest{
				Key:  key,
				HMAC: validHMAC,
			},
			expectedStatus: 400,
			expectError:    true,
		},
		{
			name: "missing key",
			requestBody: domain.HMACRequest{
				Message: message,
				HMAC:    validHMAC,
			},
			expectedStatus: 400,
			expectError:    true,
		},
		{
			name: "missing hmac",
			requestBody: domain.HMACRequest{
				Message: message,
				Key:     key,
			},
			expectedStatus: 400,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			req := httptest.NewRequest("POST", "/hmac/decrypt", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			responseBody, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}

			if tt.expectError {
				var errorResp domain.ErrorResponse
				if err := json.Unmarshal(responseBody, &errorResp); err != nil {
					t.Errorf("Failed to unmarshal error response: %v", err)
				}
				if errorResp.Error == "" {
					t.Errorf("Expected error in response, but got none")
				}
			} else {
				var successResp domain.HMACDecryptResponse
				if err := json.Unmarshal(responseBody, &successResp); err != nil {
					t.Errorf("Failed to unmarshal success response: %v", err)
				}
				if successResp.Valid != tt.expectedValid {
					t.Errorf("Expected valid=%v, got valid=%v", tt.expectedValid, successResp.Valid)
				}
			}
		})
	}
}

func TestHandler_HealthCheck(t *testing.T) {
	app, _ := setupTestApp()

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var healthResp map[string]interface{}
	if err := json.Unmarshal(responseBody, &healthResp); err != nil {
		t.Errorf("Failed to unmarshal health response: %v", err)
	}

	if status, ok := healthResp["status"]; !ok || status != "healthy" {
		t.Errorf("Expected status 'healthy', got %v", status)
	}
}