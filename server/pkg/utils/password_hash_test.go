package utils_test

import (
	"testing"

	"github.com/tsntt/footballapi/pkg/utils"
)

func TestHashPasswordAndCheckPasswordHash(t *testing.T) {
	password := "password123"

	hash, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if !utils.CheckPasswordHash(password, hash) {
		t.Error("password hash does not match")
	}
}

func TestCheckPasswordHash_InvalidPassword(t *testing.T) {
	password := "password123"
	hash, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if utils.CheckPasswordHash("wrong-password", hash) {
		t.Error("expected password check to fail for invalid password")
	}
}
