package service

import (
	"crypto/rsa"
	"fmt"
	"sync"
	"time"

	jwtgo "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwk"

	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/jwt"
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

// GetJWKS returns the JSON Web Key Set
func (s *JWKService) GetJWKS() ([]jwt.JWK, error) {
	s.keyMutex.RLock()
	defer s.keyMutex.RUnlock()

	var keys []jwt.JWK

	// Add current key
	if s.currentKey != nil {
		currentJWK := jwt.JWK{
			KeyID:     s.currentKeyID,
			KeyType:   "RSA",
			Algorithm: "RS256",
			Use:       "sig",
			N:         fmt.Sprintf("%x", s.currentKey.PublicKey.N),
			E:         fmt.Sprintf("%x", s.currentKey.PublicKey.E),
		}
		keys = append(keys, currentJWK)
	}

	// Add previous key if exists
	if s.previousKey != nil {
		previousJWK := jwt.JWK{
			KeyID:     s.previousKeyID,
			KeyType:   "RSA",
			Algorithm: "RS256",
			Use:       "sig",
			N:         fmt.Sprintf("%x", s.previousKey.PublicKey.N),
			E:         fmt.Sprintf("%x", s.previousKey.PublicKey.E),
		}
		keys = append(keys, previousJWK)
	}

	return keys, nil
}

// GenerateAccessToken generates a new JWT access token
func (s *JWKService) GenerateAccessToken(userID uuid.UUID) (string, error) {
	s.keyMutex.RLock()
	defer s.keyMutex.RUnlock()

	claims := jwtgo.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(15 * time.Minute).Unix(),
		"iat": time.Now().Unix(),
		"kid": s.currentKeyID,
	}

	token := jwtgo.NewWithClaims(jwtgo.SigningMethodRS256, claims)
	token.Header["kid"] = s.currentKeyID

	return token.SignedString(s.currentKey)
}

// GenerateRefreshToken generates a new refresh token
func (s *JWKService) GenerateRefreshToken() (string, error) {
	return uuid.NewString(), nil
}

// ValidateToken validates a JWT token
func (s *JWKService) ValidateToken(tokenString string) (*jwt.Claims, error) {
	s.keyMutex.RLock()
	defer s.keyMutex.RUnlock()

	token, err := jwtgo.Parse(tokenString, func(token *jwtgo.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtgo.SigningMethodRSA); !ok {
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

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwtgo.MapClaims); ok && token.Valid {
		return &jwt.Claims{
			UserID: claims["sub"].(string),
			RegisteredClaims: jwtgo.RegisteredClaims{
				ExpiresAt: jwtgo.NewNumericDate(time.Unix(int64(claims["exp"].(float64)), 0)),
				IssuedAt:  jwtgo.NewNumericDate(time.Unix(int64(claims["iat"].(float64)), 0)),
			},
		}, nil
	}

	return nil, fmt.Errorf("invalid token")
}
