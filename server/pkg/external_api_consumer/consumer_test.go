package consumer_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	consumer "github.com/tsntt/footballapi/pkg/external_api_consumer"
	"github.com/tsntt/footballapi/internal/dto"
	"github.com/tsntt/footballapi/internal/model"
)

func TestFootballAPIClient_GetChampionships(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := dto.ChampionshipsResponse{
			Competitions: []model.Championship{
				{ID: 1, Name: "Test Championship"},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := consumer.NewFootballAPIClient(server.URL, "test-token")
	championships, err := client.GetChampionships(context.Background())

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

func TestFootballAPIClient_GetMatches(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := dto.MatchesResponse{
			Matches: []model.Match{
				{ID: 1, Status: "SCHEDULED"},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := consumer.NewFootballAPIClient(server.URL, "test-token")
	matches, err := client.GetMatches(context.Background(), 1, "", "")

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

func TestFootballAPIClient_GetMatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		match := model.Match{ID: 1, Status: "SCHEDULED"}
		json.NewEncoder(w).Encode(match)
	}))
	defer server.Close()

	client := consumer.NewFootballAPIClient(server.URL, "test-token")
	match, err := client.GetMatch(context.Background(), 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if match.Status != "SCHEDULED" {
		t.Errorf("unexpected match status: %s", match.Status)
	}
}

func TestFootballAPIClient_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := consumer.NewFootballAPIClient(server.URL, "test-token")
	_, err := client.GetChampionships(context.Background())

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}
