package model

import (
	"context"
	"time"
)

type Championship struct {
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	Code          string   `json:"code"`
	Type          string   `json:"type"`
	Emblem        string   `json:"emblem"`
	CurrentSeason *Season  `json:"currentSeason"`
	Seasons       []Season `json:"seasons"`
}

type Season struct {
	ID              int    `json:"id"`
	StartDate       string `json:"startDate"`
	EndDate         string `json:"endDate"`
	CurrentMatchday int    `json:"currentMatchday"`
	Winner          *Team  `json:"winner"`
}

type Team struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
	TLA       string `json:"tla"`
	Crest     string `json:"crest"`
}

type Match struct {
	ID          int          `json:"id"`
	UTCDate     time.Time    `json:"utcDate"`
	Status      string       `json:"status"`
	Matchday    int          `json:"matchday"`
	Stage       string       `json:"stage"`
	Group       string       `json:"group"`
	LastUpdated time.Time    `json:"lastUpdated"`
	HomeTeam    Team         `json:"homeTeam"`
	AwayTeam    Team         `json:"awayTeam"`
	Score       Score        `json:"score"`
	Competition Championship `json:"competition"`
	Season      Season       `json:"season"`
}

type Score struct {
	Winner   string    `json:"winner"`
	Duration string    `json:"duration"`
	FullTime ScoreTime `json:"fullTime"`
	HalfTime ScoreTime `json:"halfTime"`
}

type ScoreTime struct {
	Home int `json:"home"`
	Away int `json:"away"`
}

type IChampionshipAPI interface {
	GetChampionships(ctx context.Context) ([]Championship, error)
	GetMatches(ctx context.Context, championshipID int, team, stage string) ([]Match, error)
	GetMatch(ctx context.Context, matchID int) (*Match, error)
}
