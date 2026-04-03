package service

import (
	"context"
	"errors"
	"fmt"
	user "likexuser/model"
	"likexuser/repository"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	likexService "github.com/qosdil/like-x/backend/common/service"
	"golang.org/x/crypto/bcrypt"
)

const (
	jwtSecretKeyMinLen = 32
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// Service defines business logic for user-related operations.
type Service struct {
	repo repository.Repository
}

// NewService constructs a new Service with the provided user repository.
func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

// Authenticate validates user credentials and returns a JWT token.
func (s *Service) Authenticate(ctx context.Context, input user.AuthInput) (user.AuthOutput, error) {
	// Get password hash by public_id
	passwordHash, err := s.repo.FirstPasswordHashByPublicID(ctx, input.PublicID)
	if err != nil {
		if err != likexService.ErrNotFound {
			log.Printf("failed to get password hash by public_id: %v", err)
			return user.AuthOutput{}, likexService.ErrInternal
		}

		log.Printf("debug: public_id %s not found in db", input.PublicID)
		return user.AuthOutput{}, ErrInvalidCredentials
	}

	// Validate password against hash
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(input.Password)); err != nil {
		log.Printf("debug: password and hash not a match for public_id %s", input.PublicID)
		return user.AuthOutput{}, ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := generateJWT(string(input.PublicID))
	if err != nil {
		log.Printf("failed to generate JWT: %v", err)
		return user.AuthOutput{}, likexService.ErrInternal
	}

	log.Printf("authentication successful for user %v", input.PublicID)
	return user.AuthOutput{Token: token}, nil
}

// SignUp validates user input and creates a new user via the repository.
func (s *Service) SignUp(ctx context.Context, input user.CreateInput) (user.CreateOutput, error) {
	// Validate full_name
	if len(input.FullName) < user.FullNameMinLength || len(input.FullName) > user.FullNameMaxlength {
		return user.CreateOutput{}, likexService.ErrBadRequest
	}

	// Validate password
	if len(input.Password) < user.PasswordMinLength || len(input.Password) > user.PasswordMaxLength {
		return user.CreateOutput{}, likexService.ErrBadRequest
	}

	// Create the user in the repository and handle any errors.
	signUp, err := s.repo.Create(ctx, input)
	if err != nil {
		log.Printf("failed to sign up a user: %v", err)
		return user.CreateOutput{}, likexService.ErrInternal
	}

	log.Printf("sign-up successful for %v, %v", input.FullName, input.Password)
	return signUp, nil
}

// generateJWT creates a signed JWT token for the user.
func generateJWT(publicID string) (string, error) {
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
