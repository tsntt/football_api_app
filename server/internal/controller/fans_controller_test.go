package controller_test

import (
	"context"
	"errors"
	"testing"

	"github.com/tsntt/footballapi/internal/controller"
	"github.com/tsntt/footballapi/internal/dto"
	"github.com/tsntt/footballapi/internal/model"
)

func TestFanController_Subscribe(t *testing.T) {
	mockFanRepo := &mockFanRepository{
		create: func(ctx context.Context, fan *model.Fan) error {
			return nil
		},
	}

	fanController := controller.NewFanController(mockFanRepo)
	req := &dto.FanRequest{
		UserID:   1,
		TeamID:   1,
		TeamName: "Test Team",
	}

	resp, err := fanController.Subscribe(context.Background(), req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.Message != "Subscribed to Test Team" {
		t.Errorf("unexpected response message: %s", resp.Message)
	}
}

func TestFanController_Subscribe_ValidationError(t *testing.T) {
	fanController := controller.NewFanController(nil)
	req := &dto.FanRequest{}

	_, err := fanController.Subscribe(context.Background(), req)

	if err == nil {
		t.Fatal("expected a validation error, got nil")
	}
}

func TestFanController_Subscribe_CreateError(t *testing.T) {
	mockFanRepo := &mockFanRepository{
		create: func(ctx context.Context, fan *model.Fan) error {
			return errors.New("create error")
		},
	}

	fanController := controller.NewFanController(mockFanRepo)
	req := &dto.FanRequest{
		UserID:   1,
		TeamID:   1,
		TeamName: "Test Team",
	}

	_, err := fanController.Subscribe(context.Background(), req)

	if err == nil {
		t.Fatal("expected a create error, got nil")
	}
}
