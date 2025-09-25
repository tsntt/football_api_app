package controller

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/tsntt/footballapi/internal/dto"
	"github.com/tsntt/footballapi/internal/model"
)

type FanController struct {
	fanRepo   model.IFanRepository
	validator *validator.Validate
}

func NewFanController(fanRepo model.IFanRepository) *FanController {
	return &FanController{
		fanRepo:   fanRepo,
		validator: validator.New(),
	}
}

func (c *FanController) Subscribe(ctx context.Context, req *dto.FanRequest) (*dto.APIResponse, error) {
	// Validate data
	if err := c.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	fan := &model.Fan{
		UserID: req.UserID,
		Team:   req.Team,
	}

	if err := c.fanRepo.Create(ctx, fan); err != nil {
		return nil, fmt.Errorf("failed to subscribe to team: %w", err)
	}

	return &dto.APIResponse{
		Message: fmt.Sprintf("Subscribed to %s", req.Team),
	}, nil
}
