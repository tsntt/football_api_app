package controller

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/tsntt/footballapi/internal/model"
)

type ChampionshipController struct {
	externalAPI model.IChampionshipAPI
	validator   *validator.Validate
}

func NewChampionshipController(externalAPI model.IChampionshipAPI) *ChampionshipController {
	return &ChampionshipController{
		externalAPI: externalAPI,
		validator:   validator.New(),
	}
}

func (c *ChampionshipController) GetChampionships(ctx context.Context) ([]model.Championship, error) {
	championships, err := c.externalAPI.GetChampionships(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get championships: %w", err)
	}
	return championships, nil
}

func (c *ChampionshipController) GetMatches(ctx context.Context, championshipIDStr, team, stage string) ([]model.Match, error) {
	championshipID, err := strconv.Atoi(championshipIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid championship ID: %w", err)
	}

	matches, err := c.externalAPI.GetMatches(ctx, championshipID, team, stage)
	if err != nil {
		return nil, fmt.Errorf("failed to get matches: %w", err)
	}
	return matches, nil
}
