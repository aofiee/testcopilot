package usecase

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"

	"hmac-service/internal/domain"
)

type hmacService struct{}

// NewHMACService creates a new instance of HMACService
func NewHMACService() domain.HMACService {
	return &hmacService{}
}

// Generate creates a HMAC-SHA256 hash from the given message and key
func (s *hmacService) Generate(message, key string) (string, error) {
	if key == "" {
		return "", errors.New("key is required")
	}
	if message == "" {
		return "", errors.New("message is required")
	}

	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(message))
	sum := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(sum), nil
}

// Verify validates the given HMAC against the message and key
func (s *hmacService) Verify(message, key, receivedHMAC string) (bool, error) {
	if key == "" {
		return false, errors.New("key is required")
	}
	if message == "" {
		return false, errors.New("message is required")
	}
	if receivedHMAC == "" {
		return false, errors.New("hmac is required")
	}

	expected, err := s.Generate(message, key)
	if err != nil {
		return false, err
	}

	// Use hmac.Equal to prevent timing attacks
	return hmac.Equal([]byte(expected), []byte(receivedHMAC)), nil
}