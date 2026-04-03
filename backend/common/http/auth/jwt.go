package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	jwtSecretKeyMinLen = 32
)

// GenerateJWT creates a signed JWT token for the user.
func GenerateJWT(publicID string) (string, error) {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if len(secretKey) < jwtSecretKeyMinLen {
		return "", fmt.Errorf("JWT_SECRET_KEY not set properly")
	}

	now := time.Now()
	tokenExp := now.Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"sub": publicID,
		"exp": tokenExp.Unix(),
		"iat": now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}
