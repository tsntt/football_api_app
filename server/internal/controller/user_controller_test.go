package controller_test

import (
	"context"
	"errors"
	"testing"

	"github.com/tsntt/footballapi/internal/controller"
	"github.com/tsntt/footballapi/internal/dto"
	"github.com/tsntt/footballapi/internal/model"
	"github.com/tsntt/footballapi/pkg/utils"
)

func BenchmarkUserController_Register(b *testing.B) {
	mockUserRepo := &mockUserRepository{
		create: func(ctx context.Context, user *model.User) error {
			return nil
		},
	}

	userController := controller.NewUserController(mockUserRepo, nil)
	req := &dto.UserRequest{
		Name:     "testuser",
		Password: "password",
	}

	for i := 0; i < b.N; i++ {
		_, _ = userController.Register(context.Background(), req)
	}
}

func BenchmarkUserController_Login(b *testing.B) {
	hashedPassword, _ := utils.HashPassword("password")
	mockUserRepo := &mockUserRepository{
		getByName: func(ctx context.Context, name string) (*model.User, error) {
			return &model.User{ID: 1, Name: "testuser", Password: hashedPassword, Role: "default"}, nil
		},
	}

	jms := utils.NewJWTService("secret", 24)
	userController := controller.NewUserController(mockUserRepo, jms)
	req := &dto.UserRequest{
		Name:     "testuser",
		Password: "password",
	}

	for i := 0; i < b.N; i++ {
		_, _ = userController.Login(context.Background(), req)
	}
}

func TestUserController_Register(t *testing.T) {
	mockUserRepo := &mockUserRepository{
		create: func(ctx context.Context, user *model.User) error {
			return nil
		},
	}

	userController := controller.NewUserController(mockUserRepo, nil)
	req := &dto.UserRequest{
		Name:     "testuser",
		Password: "password",
	}

	resp, err := userController.Register(context.Background(), req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.Message != "User successfully created!" {
		t.Errorf("unexpected response message: %s", resp.Message)
	}
}

func TestUserController_Register_ValidationError(t *testing.T) {
	userController := controller.NewUserController(nil, nil)
	req := &dto.UserRequest{}

	_, err := userController.Register(context.Background(), req)

	if err == nil {
		t.Fatal("expected a validation error, got nil")
	}
}

func TestUserController_Register_CreateError(t *testing.T) {
	mockUserRepo := &mockUserRepository{
		create: func(ctx context.Context, user *model.User) error {
			return errors.New("create error")
		},
	}

	userController := controller.NewUserController(mockUserRepo, nil)
	req := &dto.UserRequest{
		Name:     "testuser",
		Password: "password",
	}

	_, err := userController.Register(context.Background(), req)

	if err == nil {
		t.Fatal("expected a create error, got nil")
	}
}

func TestUserController_Login(t *testing.T) {
	hashedPassword, _ := utils.HashPassword("password")
	mockUserRepo := &mockUserRepository{
		getByName: func(ctx context.Context, name string) (*model.User, error) {
			return &model.User{ID: 1, Name: "testuser", Password: hashedPassword, Role: "default"}, nil
		},
	}

	jms := utils.NewJWTService("secret", 24)
	userController := controller.NewUserController(mockUserRepo, jms)
	req := &dto.UserRequest{
		Name:     "testuser",
		Password: "password",
	}

	resp, err := userController.Login(context.Background(), req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.Token == "" {
		t.Error("expected a token, got an empty string")
	}
}

func TestUserController_Login_InvalidCredentials(t *testing.T) {
	mockUserRepo := &mockUserRepository{
		getByName: func(ctx context.Context, name string) (*model.User, error) {
			return nil, errors.New("not found")
		},
	}

	userController := controller.NewUserController(mockUserRepo, nil)
	req := &dto.UserRequest{
		Name:     "testuser",
		Password: "password",
	}

	_, err := userController.Login(context.Background(), req)

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestUserController_Logout(t *testing.T) {
	userController := controller.NewUserController(nil, nil)
	resp, err := userController.Logout(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.Message != "Logged out successfully" {
		t.Errorf("unexpected response message: %s", resp.Message)
	}
}