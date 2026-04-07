package auth

import (
	"os"
	"testing"
)

func TestParsePublicIDFromToken_Success(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "supersecretkeythatisatleast32characterslong")
	defer os.Unsetenv("JWT_SECRET_KEY")

	token, err := GenerateJWT("pub-1")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	got, err := ParsePublicIDFromToken(token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "pub-1" {
		t.Fatalf("expected public id 'pub-1', got '%s'", got)
	}
}

func TestParsePublicIDFromToken_InvalidToken(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "supersecretkeythatisatleast32characterslong")
	defer os.Unsetenv("JWT_SECRET_KEY")

	_, err := ParsePublicIDFromToken("bad.token.value")
	if err == nil {
		t.Fatal("expected error for invalid token, got nil")
	}
}
