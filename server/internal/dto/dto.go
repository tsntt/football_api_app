package dto

import "github.com/tsntt/footballapi/internal/model"

type UserRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=50"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type JWTClaims struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	Exp    int64  `json:"exp"`
}

type APIResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ChampionshipsResponse struct {
	Competitions []model.Championship `json:"competitions"`
}

type MatchesResponse struct {
	Matches []model.Match `json:"matches"`
}

type FanRequest struct {
	UserID   int    `json:"user_id" validate:"required"`
	TeamID   int    `json:"team_id" validate:"required"`
	TeamName string `json:"team_name" validate:"required"`
}
