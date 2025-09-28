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
		TeamID: req.TeamID,
	}

	if err := c.fanRepo.Create(ctx, fan); err != nil {
		return nil, fmt.Errorf("failed to subscribe to team: %w", err)
	}

	return &dto.APIResponse{
		Message: fmt.Sprintf("Subscribed to %s", req.TeamName),
	}, nil
}

func (c *FanController) Unsubscribe(ctx context.Context, userID int, req *dto.UnsubscribeRequest) (*dto.APIResponse, error) {
	if err := c.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	if err := c.fanRepo.DeleteByUserIDAndTeam(ctx, userID, req.TeamID); err != nil {
		return nil, fmt.Errorf("failed to unsubscribe from team: %w", err)
	}

	return &dto.APIResponse{
		Message: "Unsubscribed!",
		Data: map[string]interface{}{
			"team_id": req.TeamID,
		},
	}, nil
}

func (c *FanController) GetSubscriptions(ctx context.Context, userID int) ([]model.Fan, error) {
	fans, err := c.fanRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user subscriptions: %w", err)
	}

	return fans, nil
}
