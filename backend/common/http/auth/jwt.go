package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	jwtSecretKeyMinLen = 32
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

// GenerateJWT creates a signed JWT token for the user.
func GenerateJWT(publicID string) (string, error) {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if len(secretKey) < jwtSecretKeyMinLen {
		return "", fmt.Errorf("JWT_SECRET_KEY not set properly")
	}

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": publicID,
		"exp": now.Add(24 * time.Hour).Unix(),
		"iat": now.Unix(),
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ParsePublicIDFromToken validates the JWT and returns the public ID from the subject claim.
func ParsePublicIDFromToken(tokenString string) (string, error) {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if len(secretKey) < jwtSecretKeyMinLen {
		return "", fmt.Errorf("JWT_SECRET_KEY not set properly")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method)
		}
		return []byte(secretKey), nil
	})
	if err != nil || token == nil || !token.Valid {
		return "", ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrInvalidToken
	}

	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return "", ErrInvalidToken
	}

	return sub, nil
}
