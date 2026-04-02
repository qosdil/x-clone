package http

import (
	"bytes"
	"context"
	"encoding/json"
	user "likexuser/model"
	"likexuser/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
)

type fakeRepo struct {
	createOutput user.CreateOutput
	createErr    error
}

func (f fakeRepo) Create(ctx context.Context, input user.CreateInput) (user.CreateOutput, error) {
	return f.createOutput, f.createErr
}

func TestHandleSignUp_Success(t *testing.T) {
	app := fiber.New()
	fake := fakeRepo{createOutput: user.CreateOutput{ID: 1, PublicID: "pub-123"}}
	svc := service.NewService(fake)
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
	svc := service.NewService(fake)
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
