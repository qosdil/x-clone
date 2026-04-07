package http

import (
	"bytes"
	"context"
	"encoding/json"
	user "likexuser/model"
	"likexuser/service"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/qosdil/like-x/backend/common/http/auth"
	likexService "github.com/qosdil/like-x/backend/common/service"
	"golang.org/x/crypto/bcrypt"
)

type fakeRepo struct {
	createOutput user.CreateOutput
	createErr    error
	firstHash    string
	firstErr     error
	firstID      user.ID
	firstIDErr   error
}

func (f fakeRepo) Create(ctx context.Context, input user.CreateInput) (user.CreateOutput, error) {
	return f.createOutput, f.createErr
}

func (f fakeRepo) FirstPasswordHashByPublicID(ctx context.Context, publicID user.PublicID) (string, error) {
	return f.firstHash, f.firstErr
}

func (f fakeRepo) FirstIDByPublicID(ctx context.Context, publicID user.PublicID) (user.ID, error) {
	return f.firstID, f.firstIDErr
}

type fakeAuthenticator struct {
	token string
	err   error
}

func (a fakeAuthenticator) CompareHashAndPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (a fakeAuthenticator) GenerateToken(_ string) (string, error) {
	return a.token, a.err
}

func TestHandleSignUp_Success(t *testing.T) {
	app := fiber.New()
	fake := fakeRepo{createOutput: user.CreateOutput{ID: 1, PublicID: "pub-123"}}
	svc := service.NewService(fakeAuthenticator{token: "token"}, fakeAuthenticator{token: "token"}, fake)
	h := NewHandler(svc)
	app.Post("/v1/users/sign-up", h.HandleSignUp)

	body := map[string]string{"full_name": "John Doe", "password": "secret123"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/users/sign-up", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var respBody map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		t.Fatalf("failed decode body: %v", err)
	}
	if respBody["id"] != "pub-123" {
		t.Fatalf("expected id pub-123, got %q", respBody["id"])
	}
}

func TestHandleSignUp_BadRequest(t *testing.T) {
	app := fiber.New()
	fake := fakeRepo{createOutput: user.CreateOutput{}, createErr: nil}
	svc := service.NewService(fakeAuthenticator{token: "token"}, fakeAuthenticator{token: "token"}, fake)
	h := NewHandler(svc)
	app.Post("/v1/users/sign-up", h.HandleSignUp)

	body := map[string]string{"full_name": "J", "password": "1"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/users/sign-up", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestHandleAuthenticate_Success(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "supersecretkeythatisatleast32characterslong")
	defer os.Unsetenv("JWT_SECRET_KEY")

	passHash, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	app := fiber.New()
	fake := fakeRepo{firstHash: string(passHash)}
	svc := service.NewService(fakeAuthenticator{token: "token"}, fakeAuthenticator{token: "token"}, fake)
	h := NewHandler(svc)
	app.Post("/v1/users/authenticate", h.HandleAuthenticate)

	body := map[string]string{"id": "pub-1", "password": "secret123"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/users/authenticate", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var respBody map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		t.Fatalf("failed decode body: %v", err)
	}
	if respBody["token"] == "" {
		t.Fatal("expected token in response")
	}
}

func TestHandleAuthenticate_Unauthorized(t *testing.T) {
	app := fiber.New()
	fake := fakeRepo{firstErr: likexService.ErrNotFound}
	svc := service.NewService(fakeAuthenticator{token: "token"}, fakeAuthenticator{token: "token"}, fake)
	h := NewHandler(svc)
	app.Post("/v1/users/authenticate", h.HandleAuthenticate)

	body := map[string]string{"id": "pub-1", "password": "wrong"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/users/authenticate", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test error: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", resp.StatusCode)
	}
}

func TestHandleInternalAuthenticate_Success(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "supersecretkeythatisatleast32characterslong")
	defer os.Unsetenv("JWT_SECRET_KEY")

	token, err := auth.GenerateJWT("pub-1")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	app := fiber.New()
	fake := fakeRepo{firstID: 42}
	svc := service.NewService(fakeAuthenticator{token: "token"}, fakeAuthenticator{token: "token"}, fake)
	h := NewHandler(svc)
	app.Post("/v1/users/internal/authenticate", h.HandleInternalAuthenticate)

	body := map[string]string{"token": token}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/users/internal/authenticate", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var respBody struct {
		ID uint `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		t.Fatalf("failed decode body: %v", err)
	}
	if respBody.ID != 42 {
		t.Fatalf("expected id 42, got %d", respBody.ID)
	}
}

func TestHandleInternalAuthenticate_Unauthorized(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "supersecretkeythatisatleast32characterslong")
	defer os.Unsetenv("JWT_SECRET_KEY")

	app := fiber.New()
	fake := fakeRepo{firstIDErr: likexService.ErrNotFound}
	svc := service.NewService(fakeAuthenticator{token: "token"}, fakeAuthenticator{token: "token"}, fake)
	h := NewHandler(svc)
	app.Post("/v1/users/internal/authenticate", h.HandleInternalAuthenticate)

	body := map[string]string{"token": "bad.token.value"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/users/internal/authenticate", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test error: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", resp.StatusCode)
	}
}
