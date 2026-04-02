package service

import (
	"context"
	user "likexuser/model"
	"testing"

	likexService "github.com/qosdil/like-x/backend/common/service"
)

type mockRepository struct {
	createOutput user.CreateOutput
	createErr    error
	lastInput    user.CreateInput
}

func (m *mockRepository) Create(ctx context.Context, input user.CreateInput) (user.CreateOutput, error) {
	m.lastInput = input
	return m.createOutput, m.createErr
}

func TestSignUp_Valid(t *testing.T) {
	m := &mockRepository{createOutput: user.CreateOutput{ID: 1, PublicID: "public-1"}}
	svc := NewService(m)

	out, err := svc.SignUp(context.Background(), user.CreateInput{FullName: "John Doe", Password: "secret123"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if out.PublicID != "public-1" {
		t.Fatalf("expected public id 'public-1', got '%s'", out.PublicID)
	}
}

func TestSignUp_InvalidFullName(t *testing.T) {
	m := &mockRepository{}
	svc := NewService(m)

	_, err := svc.SignUp(context.Background(), user.CreateInput{FullName: "Jan", Password: "secret123"})
	if err != likexService.ErrBadRequest {
		t.Fatalf("expected ErrBadRequest, got %v", err)
	}
}

func TestSignUp_InvalidPassword(t *testing.T) {
	m := &mockRepository{}
	svc := NewService(m)

	_, err := svc.SignUp(context.Background(), user.CreateInput{FullName: "John Doe", Password: "123"})
	if err != likexService.ErrBadRequest {
		t.Fatalf("expected ErrBadRequest, got %v", err)
	}
}

func TestSignUp_RepositoryError(t *testing.T) {
	m := &mockRepository{createErr: likexService.ErrInternal}
	svc := NewService(m)

	_, err := svc.SignUp(context.Background(), user.CreateInput{FullName: "John Doe", Password: "secret123"})
	if err != likexService.ErrInternal {
		t.Fatalf("expected ErrInternal, got %v", err)
	}
}
