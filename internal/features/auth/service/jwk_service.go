package service

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwk"
)

type JWKService struct {
	currentKey     *rsa.PrivateKey
	currentKeyID   string
	previousKey    *rsa.PrivateKey
	previousKeyID  string
	keyMutex       sync.RWMutex
	rotationPeriod time.Duration
}

type JWKResponse struct {
	Keys []jwk.Key `json:"keys"`
}

// NewJWKService creates a new JWK service with automatic key rotation
func NewJWKService(rotationPeriod time.Duration) (*JWKService, error) {
	service := &JWKService{
		rotationPeriod: rotationPeriod,
	}

	// Generate initial key
	if err := service.rotateKeys(); err != nil {
		return nil, err
	}

	// Start key rotation goroutine
	go service.startKeyRotation()

	return service, nil
}

// rotateKeys generates new RSA key pair and rotates the current and previous keys
func (s *JWKService) rotateKeys() error {
	s.keyMutex.Lock()
	defer s.keyMutex.Unlock()

	// Generate new key
	privateKey, err := rsa.GenerateKey(nil, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate RSA key: %w", err)
	}

	// Rotate keys
	s.previousKey = s.currentKey
	s.previousKeyID = s.currentKeyID
	s.currentKey = privateKey
	s.currentKeyID = uuid.New().String()

	return nil
}

// startKeyRotation starts the key rotation process
func (s *JWKService) startKeyRotation() {
	ticker := time.NewTicker(s.rotationPeriod)
	for range ticker.C {
		if err := s.rotateKeys(); err != nil {
			// Log error but continue running
			fmt.Printf("Error rotating keys: %v\n", err)
		}
	}
}

// GetJWKS returns the public JWK set
func (s *JWKService) GetJWKS() ([]byte, error) {
	s.keyMutex.RLock()
	defer s.keyMutex.RUnlock()

	// Create key set
	keySet := make([]interface{}, 0, 2)

	// Add current key
	if s.currentKey != nil {
		key, err := jwk.New(s.currentKey.Public())
		if err != nil {
			return nil, fmt.Errorf("failed to create JWK from current key: %w", err)
		}
		if err := key.Set(jwk.KeyIDKey, s.currentKeyID); err != nil {
			return nil, err
		}
		keySet = append(keySet, key)
	}

	// Add previous key if exists
	if s.previousKey != nil {
		key, err := jwk.New(s.previousKey.Public())
		if err != nil {
			return nil, fmt.Errorf("failed to create JWK from previous key: %w", err)
		}
		if err := key.Set(jwk.KeyIDKey, s.previousKeyID); err != nil {
			return nil, err
		}
		keySet = append(keySet, key)
	}

	// Marshal to JSON
	return json.Marshal(struct {
		Keys []interface{} `json:"keys"`
	}{
		Keys: keySet,
	})
}

// GenerateAccessToken generates a new JWT access token
func (s *JWKService) GenerateAccessToken(userID uuid.UUID) (string, error) {
	s.keyMutex.RLock()
	defer s.keyMutex.RUnlock()

	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(15 * time.Minute).Unix(),
		"iat": time.Now().Unix(),
		"kid": s.currentKeyID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = s.currentKeyID

	return token.SignedString(s.currentKey)
}

// GenerateRefreshToken generates a new refresh token
func (s *JWKService) GenerateRefreshToken() (string, error) {
	return uuid.NewString(), nil
}

// ValidateToken validates a JWT token
func (s *JWKService) ValidateToken(tokenString string) (*jwt.Token, error) {
	s.keyMutex.RLock()
	defer s.keyMutex.RUnlock()

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid header not found")
		}

		switch kid {
		case s.currentKeyID:
			return &s.currentKey.PublicKey, nil
		case s.previousKeyID:
			if s.previousKey != nil {
				return &s.previousKey.PublicKey, nil
			}
		}

		return nil, fmt.Errorf("key not found")
	})
}
