package auth

import (
	"errors"
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
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	claims := &Claims{}
	if sub, ok := mapClaims["user_id"].(string); ok {
		claims.Sub = sub
	}
	if exp, ok := mapClaims["exp"].(float64); ok {
		claims.Exp = int64(exp)
	}
	if iat, ok := mapClaims["iat"].(float64); ok {
		claims.Iat = int64(iat)
	}
	return claims, nil
}
