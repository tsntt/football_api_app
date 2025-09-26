package controller

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/tsntt/footballapi/internal/dto"
	"github.com/tsntt/footballapi/internal/model"
	"github.com/tsntt/footballapi/pkg/utils"
)

type UserController struct {
	userRepo   model.IUserRepository
	jwtService *utils.JWTService
	validator  *validator.Validate
}

func NewUserController(userRepo model.IUserRepository, jwtService *utils.JWTService) *UserController {
	return &UserController{
		userRepo:   userRepo,
		jwtService: jwtService,
		validator:  validator.New(),
	}
}

func (c *UserController) Register(ctx context.Context, req *dto.UserRequest) (*dto.APIResponse, error) {
	// Validate data
	if err := c.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Make user
	user := &model.User{
		Name:     req.Name,
		Password: hashedPassword,
		Role:     "default",
	}

	if err := c.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &dto.APIResponse{
		Message: "User successfully created!",
	}, nil
}

func (c *UserController) Login(ctx context.Context, req *dto.UserRequest) (*dto.LoginResponse, error) {
	if err := c.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	user, err := c.userRepo.GetByName(ctx, req.Name)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	token, err := c.jwtService.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &dto.LoginResponse{Token: token}, nil
}

func (c *UserController) Logout(ctx context.Context) (*dto.APIResponse, error) {
	// TODO: remove token from cookie
	return &dto.APIResponse{
		Message: "Logged out successfully",
	}, nil
}
