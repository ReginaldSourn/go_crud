package auth

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

func (t *TokenService) Generate(userID string, roles []string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"roles":   roles,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(t.secret)
}

func (t *TokenService) Validate(tokenString string) (*Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return t.secret, nil
	})
	// Extract claims, verify expiry
}
