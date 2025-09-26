package controller

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/tsntt/footballapi/internal/dto"
	"github.com/tsntt/footballapi/internal/model"
	"github.com/tsntt/footballapi/pkg/broadcaster"
)

type AdminController struct {
	externalAPI      model.IChampionshipAPI
	fanRepo          model.IFanRepository
	broadcastRepo    model.IBroadcastRepository
	broadcastService broadcaster.IBroadcastService
	validator        *validator.Validate
}

func NewAdminController(
	externalAPI model.IChampionshipAPI,
	fanRepo model.IFanRepository,
	broadcastRepo model.IBroadcastRepository,
	broadcastService broadcaster.IBroadcastService,
) *AdminController {
	return &AdminController{
		externalAPI:      externalAPI,
		fanRepo:          fanRepo,
		broadcastRepo:    broadcastRepo,
		broadcastService: broadcastService,
		validator:        validator.New(),
	}
}

func (c *AdminController) GetMatches(ctx context.Context) ([]model.Match, error) {
	// TODO: implement caching for production
	championships, err := c.externalAPI.GetChampionships(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get championships: %w", err)
	}

	var allMatches []model.Match
	for _, championship := range championships {
		matches, err := c.externalAPI.GetMatches(ctx, championship.ID, "", "")
		if err != nil {
			continue
		}
		allMatches = append(allMatches, matches...)
	}

	validStatuses := map[string]bool{
		"SCHEDULED": true,
		"LIVE":      true,
		"FINISHED":  true,
		"POSTPONED": true,
		"SUSPENDED": true,
		"CANCELLED": true,
	}

	var filteredMatches []model.Match
	for _, match := range allMatches {
		if validStatuses[match.Status] {
			filteredMatches = append(filteredMatches, match)
		}
	}

	return filteredMatches, nil
}

func (c *AdminController) BroadcastMatch(ctx context.Context, matchIDStr string) (*dto.APIResponse, error) {
	matchID, err := strconv.Atoi(matchIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid match ID: %w", err)
	}

	// Check if broadcast already sent for this match, avoid duplicates
	existing, err := c.broadcastRepo.GetByMatchID(ctx, matchID)
	if err == nil && existing != nil {
		return &dto.APIResponse{
			Message: "Broadcast already sent for this match",
		}, nil
	}

	match, err := c.externalAPI.GetMatch(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get match details: %w", err)
	}

	homeFans, err := c.fanRepo.GetByTeam(ctx, match.HomeTeam.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get home team fans: %w", err)
	}

	awayFans, err := c.fanRepo.GetByTeam(ctx, match.AwayTeam.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get away team fans: %w", err)
	}

	allFans := append(homeFans, awayFans...)
	if len(allFans) == 0 {
		return &dto.APIResponse{
			Message: "No fans found for this match",
		}, nil
	}

	// []Fan to []NotificationTarget
	var targets []broadcaster.NotificationTarget
	for i, fan := range allFans {
		targets = append(targets, broadcaster.NotificationTarget{
			ID:      fmt.Sprintf("fan_%d", i),
			Type:    "user",
			Address: strconv.Itoa(fan.UserID),
			Metadata: map[string]string{
				"team":     fan.Team,
				"fan_id":   strconv.Itoa(fan.ID),
				"match_id": matchIDStr,
			},
		})
	}

	message := fmt.Sprintf("🏆 %s vs %s - Status: %s",
		match.HomeTeam.Name,
		match.AwayTeam.Name,
		match.Status,
	)

	notificationID := fmt.Sprintf("match_%d", matchID)

	go func() {
		if err := c.broadcastService.SendNotificationWithID(context.Background(), notificationID, message, targets); err != nil {
			fmt.Printf("Broadcast error: %v\n", err)
		}
	}()

	broadcast := &model.BroadcastMessage{
		MatchID: matchID,
		Message: message,
		Status:  "sent",
	}

	if err := c.broadcastRepo.Create(ctx, broadcast); err != nil {
		return nil, fmt.Errorf("failed to save broadcast record: %w", err)
	}

	return &dto.APIResponse{
		Message: fmt.Sprintf("Broadcast sent to %d fans", len(allFans)),
		Data: map[string]interface{}{
			"match_id":        matchID,
			"notification_id": notificationID,
			"targets_count":   len(targets),
			"message":         message,
		},
	}, nil
}
