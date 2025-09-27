package domain

// HMACService defines the interface for HMAC operations
type HMACService interface {
	Generate(message, key string) (string, error)
	Verify(message, key, hmac string) (bool, error)
}

// HMACRequest represents the request payload for HMAC operations
type HMACRequest struct {
	Message string `json:"message" validate:"required"`
	Key     string `json:"key" validate:"required"`
	HMAC    string `json:"hmac,omitempty"`
}

// HMACEncryptResponse represents the response for HMAC encryption
type HMACEncryptResponse struct {
	HMAC      string `json:"hmac"`
	Algorithm string `json:"algorithm"`
}

// HMACDecryptResponse represents the response for HMAC verification
type HMACDecryptResponse struct {
	Valid bool `json:"valid"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}