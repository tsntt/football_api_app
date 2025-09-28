package controller_test

import (
	"context"
	"errors"
	"testing"

	"github.com/tsntt/footballapi/internal/controller"
	"github.com/tsntt/footballapi/internal/model"
)

func BenchmarkChampionshipController_GetChampionships(b *testing.B) {
	mockAPI := &mockChampionshipAPI{
		getChampionships: func(ctx context.Context) ([]model.Championship, error) {
			return []model.Championship{{ID: 1, Name: "Test Championship"}}, nil
		},
	}

	championshipController := controller.NewChampionshipController(mockAPI)

	for i := 0; i < b.N; i++ {
		_, _ = championshipController.GetChampionships(context.Background())
	}
}

func BenchmarkChampionshipController_GetMatches(b *testing.B) {
	mockAPI := &mockChampionshipAPI{
		getMatches: func(ctx context.Context, championshipID int, dateFrom, dateTo string) ([]model.Match, error) {
			return []model.Match{{ID: 1, Status: "SCHEDULED"}}, nil
		},
	}

	championshipController := controller.NewChampionshipController(mockAPI)

	for i := 0; i < b.N; i++ {
		_, _ = championshipController.GetMatches(context.Background(), "1", "", "")
	}
}

func TestChampionshipController_GetChampionships(t *testing.T) {
	mockAPI := &mockChampionshipAPI{
		getChampionships: func(ctx context.Context) ([]model.Championship, error) {
			return []model.Championship{{ID: 1, Name: "Test Championship"}}, nil
		},
	}

	championshipController := controller.NewChampionshipController(mockAPI)
	championships, err := championshipController.GetChampionships(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(championships) != 1 {
		t.Errorf("expected 1 championship, got %d", len(championships))
	}

	if championships[0].Name != "Test Championship" {
		t.Errorf("unexpected championship name: %s", championships[0].Name)
	}
}

func TestChampionshipController_GetMatches(t *testing.T) {
	mockAPI := &mockChampionshipAPI{
		getMatches: func(ctx context.Context, championshipID int, dateFrom, dateTo string) ([]model.Match, error) {
			return []model.Match{{ID: 1, Status: "SCHEDULED"}}, nil
		},
	}

	championshipController := controller.NewChampionshipController(mockAPI)
	matches, err := championshipController.GetMatches(context.Background(), "1", "", "")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(matches) != 1 {
		t.Errorf("expected 1 match, got %d", len(matches))
	}

	if matches[0].Status != "SCHEDULED" {
		t.Errorf("unexpected match status: %s", matches[0].Status)
	}
}

func TestChampionshipController_GetMatches_InvalidID(t *testing.T) {
	championshipController := controller.NewChampionshipController(nil)
	_, err := championshipController.GetMatches(context.Background(), "invalid", "", "")

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestChampionshipController_GetChampionships_Error(t *testing.T) {
	mockAPI := &mockChampionshipAPI{
		getChampionships: func(ctx context.Context) ([]model.Championship, error) {
			return nil, errors.New("some error")
		},
	}

	championshipController := controller.NewChampionshipController(mockAPI)
	_, err := championshipController.GetChampionships(context.Background())

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}