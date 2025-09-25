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
	broadcastRepo    broadcaster.IBroadcastRepository
	broadcastService broadcaster.IBroadcastService
	validator        *validator.Validate
}

func NewAdminController(
	externalAPI model.IChampionshipAPI,
	fanRepo model.IFanRepository,
	broadcastRepo broadcaster.IBroadcastRepository,
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
	// Para o admin, buscaremos partidas de todos os campeonatos principais
	// Em uma implementação real, você pode querer cachear esses dados
	championships, err := c.externalAPI.GetChampionships(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get championships: %w", err)
	}

	var allMatches []model.Match
	for _, championship := range championships {
		matches, err := c.externalAPI.GetMatches(ctx, championship.ID, "", "")
		if err != nil {
			continue // Continua mesmo se falhar para um campeonato
		}
		allMatches = append(allMatches, matches...)
	}

	// Filtrar apenas partidas com status específicos
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

	// Verificar se já foi enviado broadcast para esta partida
	existing, err := c.broadcastRepo.GetByMatchID(ctx, matchID)
	if err == nil && existing != nil {
		return &dto.APIResponse{
			Message: "Broadcast already sent for this match",
		}, nil
	}

	// Buscar detalhes da partida
	match, err := c.externalAPI.GetMatch(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get match details: %w", err)
	}

	// Buscar fãs dos times envolvidos
	homeFans, err := c.fanRepo.GetByTeam(ctx, match.HomeTeam.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get home team fans: %w", err)
	}

	awayFans, err := c.fanRepo.GetByTeam(ctx, match.AwayTeam.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get away team fans: %w", err)
	}

	// Combinar fãs
	allFans := append(homeFans, awayFans...)
	if len(allFans) == 0 {
		return &dto.APIResponse{
			Message: "No fans found for this match",
		}, nil
	}

	// Criar mensagem
	message := fmt.Sprintf("🏆 %s vs %s - Status: %s",
		match.HomeTeam.Name,
		match.AwayTeam.Name,
		match.Status,
	)

	// Enviar broadcast usando goroutines e channels
	go func() {
		if err := c.broadcastService.SendNotification(context.Background(), matchID, message, allFans); err != nil {
			fmt.Printf("Broadcast error: %v\n", err)
		}
	}()

	// Salvar registro do broadcast
	broadcast := &broadcaster.BroadcastMessage{
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
			"match_id":   matchID,
			"fans_count": len(allFans),
			"message":    message,
		},
	}, nil
}
