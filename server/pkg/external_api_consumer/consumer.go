package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/tsntt/footballapi/internal/dto"
	"github.com/tsntt/footballapi/internal/model"
)

type FootballAPIClient struct {
	baseURL string
	token   string
	client  *http.Client
}

func NewFootballAPIClient(baseURL, token string) *FootballAPIClient {
	return &FootballAPIClient{
		baseURL: baseURL,
		token:   token,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *FootballAPIClient) GetChampionships(ctx context.Context) ([]model.Championship, error) {
	url := fmt.Sprintf("%s/competitions", c.baseURL)

	resp, err := c.makeRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response dto.ChampionshipsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.Competitions, nil
}

func (c *FootballAPIClient) GetMatches(ctx context.Context, championshipID int, team, stage string) ([]model.Match, error) {
	baseURL := fmt.Sprintf("%s/competitions/%d/matches", c.baseURL, championshipID)

	// Build query parameters
	params := url.Values{}
	if stage != "" {
		params.Add("stage", stage)
	}

	fullURL := baseURL
	if len(params) > 0 {
		fullURL += "?" + params.Encode()
	}

	resp, err := c.makeRequest(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response dto.MatchesResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if team != "" {
		var filteredResponse dto.MatchesResponse
		for _, match := range response.Matches {
			if match.HomeTeam.ShortName == team {
				filteredResponse.Matches = append(filteredResponse.Matches, match)
			}
		}

		return filteredResponse.Matches, nil
	}

	return response.Matches, nil
}

func (c *FootballAPIClient) GetMatch(ctx context.Context, matchID int) (*model.Match, error) {
	url := fmt.Sprintf("%s/matches/%d", c.baseURL, matchID)

	resp, err := c.makeRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var match model.Match
	if err := json.Unmarshal(body, &match); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &match, nil
}

func (c *FootballAPIClient) makeRequest(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Auth-Token", c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return resp, nil
}
