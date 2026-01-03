package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
)

// TokenService issues and validates HS256 JWTs for this service.
type TokenService struct {
	secret []byte
	ttl    time.Duration
}

func NewTokenService(secret []byte, ttl time.Duration) (*TokenService, error) {
	if len(secret) == 0 {
		return nil, errors.New("secret is required")
	}
	if ttl <= 0 {
		return nil, errors.New("ttl must be positive")
	}
	return &TokenService{secret: secret, ttl: ttl}, nil
}

func (s *TokenService) GenerateToken(subject string) (string, error) {
	headerJSON, err := json.Marshal(map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	})
	if err != nil {
		return "", fmt.Errorf("encode header: %w", err)
	}

	now := time.Now().Unix()
	payload := Claims{
		Sub: subject,
		Exp: time.Now().Add(s.ttl).Unix(),
		Iat: now,
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("encode claims: %w", err)
	}

	header := base64.RawURLEncoding.EncodeToString(headerJSON)
	body := base64.RawURLEncoding.EncodeToString(payloadJSON)
	unsigned := header + "." + body

	signature := signHS256(unsigned, s.secret)
	return unsigned + "." + signature, nil
}

func (s *TokenService) ParseToken(token string) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Claims{}, ErrInvalidToken
	}

	unsigned := parts[0] + "." + parts[1]
	expected := signHS256(unsigned, s.secret)
	if !hmac.Equal([]byte(expected), []byte(parts[2])) {
		return Claims{}, ErrInvalidToken
	}

	payloadJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return Claims{}, ErrInvalidToken
	}

	var claims Claims
	if err := json.Unmarshal(payloadJSON, &claims); err != nil {
		return Claims{}, ErrInvalidToken
	}

	if time.Now().Unix() > claims.Exp {
		return Claims{}, ErrExpiredToken
	}

	return claims, nil
}

func signHS256(unsigned string, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(unsigned))
	sum := mac.Sum(nil)
	return base64.RawURLEncoding.EncodeToString(sum)
}
