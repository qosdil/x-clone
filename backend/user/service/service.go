package service

import (
	"context"
	"errors"
	user "likexuser/model"
	"likexuser/repository"
	"log"

	"github.com/qosdil/like-x/backend/common/http/auth"
	likexService "github.com/qosdil/like-x/backend/common/service"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// Service defines business logic for user-related operations.
type Service struct {
	auth     authenticator
	httpauth httpauthenticator
	repo     repository.Repository
}

// NewService constructs a new Service with the provided user repository.
func NewService(auth authenticator, httpauth httpauthenticator, repo repository.Repository) *Service {
	return &Service{auth: auth, httpauth: httpauth, repo: repo}
}

// Authenticate validates user credentials and returns an auth token.
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
	if err := s.auth.CompareHashAndPassword(passwordHash, input.Password); err != nil {
		log.Printf("debug: password and hash not a match for public_id %s", input.PublicID)
		return user.AuthOutput{}, ErrInvalidCredentials
	}

	token, err := s.httpauth.GenerateToken(string(input.PublicID))
	if err != nil {
		log.Printf("failed to generate token: %v", err)
		return user.AuthOutput{}, likexService.ErrInternal
	}

	log.Printf("authentication successful for user %v", input.PublicID)
	return user.AuthOutput{Token: token}, nil
}

func (s *Service) AuthenticateInternal(ctx context.Context, authToken string) (user.AuthInternalOutput, error) {
	publicIDStr, err := auth.ParsePublicIDFromToken(authToken)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidToken) {
			return user.AuthInternalOutput{}, auth.ErrInvalidToken
		}
		log.Printf("failed to parse auth token: %v", err)
		return user.AuthInternalOutput{}, likexService.ErrInternal
	}
	publicID := user.PublicID(publicIDStr)

	id, err := s.repo.FirstIDByPublicID(ctx, publicID)
	if err != nil {
		if err == likexService.ErrNotFound {
			log.Printf("debug: public_id %s not found in db", publicID)
			return user.AuthInternalOutput{}, auth.ErrInvalidToken
		}
		log.Printf("failed to get user id by public_id: %v", err)
		return user.AuthInternalOutput{}, likexService.ErrInternal
	}

	return user.AuthInternalOutput{ID: id}, nil
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
