package usecase

import (
	"testing"
)

func TestHMACService_Generate(t *testing.T) {
	service := NewHMACService()

	tests := []struct {
		name      string
		message   string
		key       string
		wantErr   bool
		errMsg    string
	}{
		{
			name:    "valid input",
			message: "hello world",
			key:     "secret",
			wantErr: false,
		},
		{
			name:    "empty key",
			message: "hello world",
			key:     "",
			wantErr: true,
			errMsg:  "key is required",
		},
		{
			name:    "empty message",
			message: "",
			key:     "secret",
			wantErr: true,
			errMsg:  "message is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Generate(tt.message, tt.key)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("Generate() expected error but got none")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("Generate() error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("Generate() unexpected error = %v", err)
				return
			}

			if result == "" {
				t.Errorf("Generate() returned empty result")
			}
		})
	}
}

func TestHMACService_Verify(t *testing.T) {
	service := NewHMACService()
	
	// Generate a valid HMAC first
	message := "hello world"
	key := "secret"
	validHMAC, err := service.Generate(message, key)
	if err != nil {
		t.Fatalf("Failed to generate valid HMAC: %v", err)
	}

	tests := []struct {
		name     string
		message  string
		key      string
		hmac     string
		expected bool
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid hmac",
			message:  message,
			key:      key,
			hmac:     validHMAC,
			expected: true,
			wantErr:  false,
		},
		{
			name:     "invalid hmac",
			message:  message,
			key:      key,
			hmac:     "invalid-hmac",
			expected: false,
			wantErr:  false,
		},
		{
			name:     "wrong key",
			message:  message,
			key:      "wrong-key",
			hmac:     validHMAC,
			expected: false,
			wantErr:  false,
		},
		{
			name:    "empty key",
			message: message,
			key:     "",
			hmac:    validHMAC,
			wantErr: true,
			errMsg:  "key is required",
		},
		{
			name:    "empty message",
			message: "",
			key:     key,
			hmac:    validHMAC,
			wantErr: true,
			errMsg:  "message is required",
		},
		{
			name:    "empty hmac",
			message: message,
			key:     key,
			hmac:    "",
			wantErr: true,
			errMsg:  "hmac is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Verify(tt.message, tt.key, tt.hmac)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("Verify() expected error but got none")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("Verify() error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("Verify() unexpected error = %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Verify() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHMACService_ConsistentGeneration(t *testing.T) {
	service := NewHMACService()
	
	message := "test message"
	key := "test key"
	
	// Generate HMAC multiple times
	hmac1, err1 := service.Generate(message, key)
	hmac2, err2 := service.Generate(message, key)
	
	if err1 != nil || err2 != nil {
		t.Fatalf("Unexpected errors: %v, %v", err1, err2)
	}
	
	if hmac1 != hmac2 {
		t.Errorf("Generate() should produce consistent results, got %v and %v", hmac1, hmac2)
	}
}