package controller_test

import (
	"context"
	"errors"
	"testing"

	"github.com/tsntt/footballapi/internal/controller"
	"github.com/tsntt/footballapi/internal/model"
	"github.com/tsntt/footballapi/pkg/broadcaster"
)

func TestAdminController_GetMatches(t *testing.T) {
	mockAPI := &mockChampionshipAPI{
		getChampionships: func(ctx context.Context) ([]model.Championship, error) {
			return []model.Championship{{ID: 1}}, nil
		},
		getMatches: func(ctx context.Context, championshipID int, dateFrom, dateTo string) ([]model.Match, error) {
			return []model.Match{
				{Status: "SCHEDULED"},
				{Status: "LIVE"},
				{Status: "INVALID"},
			}, nil
		},
	}

	adminController := controller.NewAdminController(mockAPI, nil, nil, nil)
	matches, err := adminController.GetMatches(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(matches) != 2 {
		t.Errorf("expected 2 matches, got %d", len(matches))
	}
}

func TestAdminController_BroadcastMatch(t *testing.T) {
	mockAPI := &mockChampionshipAPI{
		getMatch: func(ctx context.Context, matchID int) (*model.Match, error) {
			return &model.Match{
				HomeTeam: model.Team{ID: 1, Name: "Home"},
				AwayTeam: model.Team{ID: 2, Name: "Away"},
				Status:   "SCHEDULED",
			}, nil
		},
	}
	mockFanRepo := &mockFanRepository{
		getByTeamID: func(ctx context.Context, teamID int) ([]model.Fan, error) {
			if teamID == 1 {
				return []model.Fan{{ID: 1, UserID: 101, TeamID: 1}}, nil
			}
			if teamID == 2 {
				return []model.Fan{{ID: 2, UserID: 102, TeamID: 2}}, nil
			}
			return nil, nil
		},
	}
	mockBroadcastRepo := &mockBroadcastRepository{
		getByMatchID: func(ctx context.Context, matchID int) (*model.BroadcastMessage, error) {
			return nil, errors.New("not found")
		},
		create: func(ctx context.Context, broadcast *model.BroadcastMessage) error {
			return nil
		},
	}
	mockBroadcastSvc := &mockBroadcastService{
		sendNotificationWithID: func(ctx context.Context, notificationID, message string, targets []broadcaster.NotificationTarget) error {
			return nil
		},
	}

	adminController := controller.NewAdminController(mockAPI, mockFanRepo, mockBroadcastRepo, mockBroadcastSvc)
	response, err := adminController.BroadcastMatch(context.Background(), 123)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if response.Message != "Broadcast sent to 2 fans" {
		t.Errorf("unexpected response message: %s", response.Message)
	}
}

func TestAdminController_BroadcastMatch_AlreadySent(t *testing.T) {
	mockBroadcastRepo := &mockBroadcastRepository{
		getByMatchID: func(ctx context.Context, matchID int) (*model.BroadcastMessage, error) {
			return &model.BroadcastMessage{}, nil
		},
	}

	adminController := controller.NewAdminController(nil, nil, mockBroadcastRepo, nil)
	response, err := adminController.BroadcastMatch(context.Background(), 123)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if response.Message != "Broadcast already sent for this match" {
		t.Errorf("unexpected response message: %s", response.Message)
	}
}

func TestAdminController_BroadcastMatch_NoFans(t *testing.T) {
	mockAPI := &mockChampionshipAPI{
		getMatch: func(ctx context.Context, matchID int) (*model.Match, error) {
			return &model.Match{
				HomeTeam: model.Team{ID: 1, Name: "Home"},
				AwayTeam: model.Team{ID: 2, Name: "Away"},
			}, nil
		},
	}
	mockFanRepo := &mockFanRepository{
		getByTeamID: func(ctx context.Context, teamID int) ([]model.Fan, error) {
			return nil, nil
		},
	}
	mockBroadcastRepo := &mockBroadcastRepository{
		getByMatchID: func(ctx context.Context, matchID int) (*model.BroadcastMessage, error) {
			return nil, errors.New("not found")
		},
	}

	adminController := controller.NewAdminController(mockAPI, mockFanRepo, mockBroadcastRepo, nil)
	response, err := adminController.BroadcastMatch(context.Background(), 123)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if response.Message != "No fans found for this match" {
		t.Errorf("unexpected response message: %s", response.Message)
	}
}
