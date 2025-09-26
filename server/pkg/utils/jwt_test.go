package utils_test

import (
	"testing"
	"time"

	"github.com/tsntt/footballapi/internal/model"
	"github.com/tsntt/footballapi/pkg/utils"
)

func TestJWTService_GenerateAndValidateToken(t *testing.T) {
	jms := utils.NewJWTService("secret", 1)
	user := &model.User{
		ID:   1,
		Name: "testuser",
		Role: "default",
	}

	token, err := jms.GenerateToken(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	claims, err := jms.ValidateToken(token)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if claims.UserID != user.ID {
		t.Errorf("expected user ID %d, got %d", user.ID, claims.UserID)
	}

	if claims.Name != user.Name {
		t.Errorf("expected user name %s, got %s", user.Name, claims.Name)
	}

	if claims.Role != user.Role {
		t.Errorf("expected user role %s, got %s", user.Role, claims.Role)
	}
}

func TestJWTService_ValidateToken_InvalidToken(t *testing.T) {
	jms := utils.NewJWTService("secret", 1)

	_, err := jms.ValidateToken("invalid-token")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestJWTService_ValidateToken_ExpiredToken(t *testing.T) {
	jms := utils.NewJWTService("secret", -1)
	user := &model.User{
		ID:   1,
		Name: "testuser",
		Role: "default",
	}

	token, err := jms.GenerateToken(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Wait for the token to expire
	time.Sleep(2 * time.Second)

	_, err = jms.ValidateToken(token)
	if err == nil {
		t.Fatal("expected an error for expired token, got nil")
	}
}
