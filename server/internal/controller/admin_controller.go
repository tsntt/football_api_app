package controller

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"github.com/tsntt/footballapi/internal/dto"
	"github.com/tsntt/footballapi/internal/model"
	"github.com/tsntt/footballapi/pkg/broadcast"
)

type AdminController struct {
	externalAPI      model.IChampionshipAPI
	fanRepo          model.IFanRepository
	broadcastRepo    model.IBroadcastRepository
	broadcastService *broadcast.BroadcastService
	validator        *validator.Validate
}

func NewAdminController(
	externalAPI model.IChampionshipAPI,
	fanRepo model.IFanRepository,
	broadcastRepo model.IBroadcastRepository,
	broadcastService *broadcast.BroadcastService,
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
	// INFO: should be cached for production
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

	// TODO: filter also by fan in db
	var filteredMatches []model.Match
	for _, match := range allMatches {
		if validStatuses[match.Status] {
			filteredMatches = append(filteredMatches, match)
		}
	}

	return filteredMatches, nil
}

func (c *AdminController) BroadcastMatch(ctx context.Context, matchID int) (*dto.APIResponse, error) {
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

	homeFans, err := c.fanRepo.GetByTeamID(ctx, match.HomeTeam.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get home team fans: %w", err)
	}

	awayFans, err := c.fanRepo.GetByTeamID(ctx, match.AwayTeam.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get away team fans: %w", err)
	}

	allFans := append(homeFans, awayFans...)
	if len(allFans) == 0 {
		return &dto.APIResponse{
			Message: "No fans found for this match",
		}, nil
	}

	for _, fan := range allFans {
		c.broadcastService.AddSubscription(broadcast.Subscription{
			UserID:           fan.ID,
			ChannelID:        fan.TeamID,
			NotificationType: broadcast.NotificationType(fan.NotificationType),
			Address:          fan.Address,
		})
	}

	message := fmt.Sprintf("🏆 %s vs %s - Status: %s",
		match.HomeTeam.Name,
		match.AwayTeam.Name,
		match.Status,
	)

	notificationID := fmt.Sprintf("match_%d", matchID)

	msg := broadcast.Message{
		Title:   "Football APP",
		Content: message,
	}

	go c.broadcastService.BroadCastToChannel(context.Background(), 5, match.HomeTeam.ID, msg)

	broadcast := &model.BroadcastMessage{
		MatchID:            matchID,
		MessageContentHash: broadcast.GenerateContentHash(msg),
		Status:             "sent",
	}

	if err := c.broadcastRepo.Create(ctx, broadcast); err != nil {
		return nil, fmt.Errorf("failed to save broadcast record: %w", err)
	}

	return &dto.APIResponse{
		Message: fmt.Sprintf("Broadcast started! notifying %d %s fans and %d %s fans.", len(homeFans), match.HomeTeam.Name, len(awayFans), match.AwayTeam.Name),
		Data: map[string]interface{}{
			"match_id":        matchID,
			"notification_id": notificationID,
			"targets_count":   len(allFans),
			"message":         message,
		},
	}, nil
}

func (c *AdminController) RegisterWS(conn *websocket.Conn) {
	c.broadcastService.RegisterAdmConn(conn)
}

func (c *AdminController) UnregisterWS(conn *websocket.Conn) {
	c.broadcastService.UnregisterAdmConn(conn)
}
