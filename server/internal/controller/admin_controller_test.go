package controller_test

import (
	"context"
	"errors"
	"testing"

	"github.com/tsntt/footballapi/internal/controller"
	"github.com/tsntt/footballapi/internal/dto"
	"github.com/tsntt/footballapi/internal/model"
	"github.com/tsntt/footballapi/pkg/broadcast"
)

func BenchmarkAdminController_GetMatches(b *testing.B) {
	mockAPI := &mockChampionshipAPI{
		getChampionships: func(ctx context.Context) ([]model.Championship, error) {
			return []model.Championship{{ID: 1}}, nil
		},
		getMatches: func(ctx context.Context, championshipID int, team, stage string) ([]model.Match, error) {
			return []model.Match{
				{Status: "SCHEDULED"},
				{Status: "LIVE"},
				{Status: "INVALID"},
			}, nil
		},
	}
	mockFanRepo := &mockFanRepository{
		getAll: func(ctx context.Context) ([]model.Fan, error) {
			return []model.Fan{{ID: 1, TeamID: 1}}, nil
		},
	}

	adminController := controller.NewAdminController(mockAPI, mockFanRepo, nil, nil)

	for i := 0; i < b.N; i++ {
		_, _ = adminController.GetMatches(context.Background())
	}
}

func BenchmarkAdminController_BroadcastMatch(b *testing.B) {
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
				return []model.Fan{{ID: 1, TeamID: 1}}, nil
			}
			if teamID == 2 {
				return []model.Fan{{ID: 2, TeamID: 2}}, nil
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

	broadcastService := broadcast.NewBroadcastService()
	adminController := controller.NewAdminController(mockAPI, mockFanRepo, mockBroadcastRepo, broadcastService)

	for i := 0; i < b.N; i++ {
		_, _ = adminController.BroadcastMatch(context.Background(), 123)
	}
}

func TestAdminController_GetMatches(t *testing.T) {
	type fields struct {
		championshipAPI *mockChampionshipAPI
		fanRepo         *mockFanRepository
		broadcastRepo   *mockBroadcastRepository
		broadcast       *broadcast.BroadcastService
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		wantLen int
	}{
		{
			name: "success",
			fields: fields{
				championshipAPI: &mockChampionshipAPI{
					getChampionships: func(ctx context.Context) ([]model.Championship, error) {
						return []model.Championship{{ID: 1}}, nil
					},
					getMatches: func(ctx context.Context, championshipID int, team, stage string) ([]model.Match, error) {
						return []model.Match{
							{Status: "SCHEDULED", HomeTeam: model.Team{ID: 1}, AwayTeam: model.Team{ID: 2}},
							{Status: "LIVE", HomeTeam: model.Team{ID: 1}, AwayTeam: model.Team{ID: 2}},
							{Status: "INVALID", HomeTeam: model.Team{ID: 1}, AwayTeam: model.Team{ID: 2}},
						}, nil
					},
				},
				fanRepo: &mockFanRepository{
					getAll: func(ctx context.Context) ([]model.Fan, error) {
						return []model.Fan{{ID: 1, TeamID: 1}}, nil
					},
				},
			},
			wantErr: false,
			wantLen: 2,
		},
		{
			name: "error getting championships",
			fields: fields{
				championshipAPI: &mockChampionshipAPI{
					getChampionships: func(ctx context.Context) ([]model.Championship, error) {
						return nil, errors.New("get championships error")
					},
				},
			},
			wantErr: true,
			wantLen: 0,
		},
		{
			name: "error getting matches",
			fields: fields{
				championshipAPI: &mockChampionshipAPI{
					getChampionships: func(ctx context.Context) ([]model.Championship, error) {
						return []model.Championship{{ID: 1}}, nil
					},
					getMatches: func(ctx context.Context, championshipID int, team, stage string) ([]model.Match, error) {
						return nil, errors.New("get matches error")
					},
				},
				fanRepo: &mockFanRepository{ // precisa inicializar
					getAll: func(ctx context.Context) ([]model.Fan, error) {
						return []model.Fan{}, nil
					},
				},
				broadcastRepo: &mockBroadcastRepository{},    // inicializa vazio se não usa
				broadcast:     &broadcast.BroadcastService{}, // idem
			},
			wantErr: true,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := controller.NewAdminController(tt.fields.championshipAPI, tt.fields.fanRepo, tt.fields.broadcastRepo, tt.fields.broadcast)
			got, err := a.GetMatches(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminController.GetMatches() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantLen {
				t.Errorf("AdminController.GetMatches() = %v, want %v", len(got), tt.wantLen)
			}
		})
	}
}

func TestAdminController_BroadcastMatch(t *testing.T) {
	type fields struct {
		championshipAPI model.IChampionshipAPI
		fanRepo         model.IFanRepository
		broadcastRepo   model.IBroadcastRepository
		broadcast       *broadcast.BroadcastService
	}
	type args struct {
		matchID int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *dto.APIResponse
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				championshipAPI: &mockChampionshipAPI{
					getMatch: func(ctx context.Context, matchID int) (*model.Match, error) {
						return &model.Match{
							HomeTeam: model.Team{ID: 1, Name: "Home"},
							AwayTeam: model.Team{ID: 2, Name: "Away"},
							Status:   "SCHEDULED",
						}, nil
					},
				},
				fanRepo: &mockFanRepository{
					getByTeamID: func(ctx context.Context, teamID int) ([]model.Fan, error) {
						if teamID == 1 {
							return []model.Fan{{ID: 1, TeamID: 1}}, nil
						}
						if teamID == 2 {
							return []model.Fan{{ID: 2, TeamID: 2}, {ID: 3, TeamID: 2}}, nil
						}
						return nil, nil
					},
				},
				broadcastRepo: &mockBroadcastRepository{
					getByMatchID: func(ctx context.Context, matchID int) (*model.BroadcastMessage, error) {
						return nil, errors.New("not found")
					},
					create: func(ctx context.Context, broadcast *model.BroadcastMessage) error {
						return nil
					},
				},
				broadcast: broadcast.NewBroadcastService(),
			},
			args: args{matchID: 123},
			want: &dto.APIResponse{Message: "Broadcast started! notifying 1 Home fans and 2 Away fans."},
		},
		{
			name: "broadcast already sent",
			fields: fields{
				broadcastRepo: &mockBroadcastRepository{
					getByMatchID: func(ctx context.Context, matchID int) (*model.BroadcastMessage, error) {
						return &model.BroadcastMessage{}, nil
					},
				},
			},
			args: args{matchID: 123},
			want: &dto.APIResponse{Message: "Broadcast already sent for this match"},
		},
		{
			name: "no fans found",
			fields: fields{
				championshipAPI: &mockChampionshipAPI{
					getMatch: func(ctx context.Context, matchID int) (*model.Match, error) {
						return &model.Match{
							HomeTeam: model.Team{ID: 1, Name: "Home"},
							AwayTeam: model.Team{ID: 2, Name: "Away"},
						}, nil
					},
				},
				fanRepo: &mockFanRepository{
					getByTeamID: func(ctx context.Context, teamID int) ([]model.Fan, error) {
						return nil, nil
					},
				},
				broadcastRepo: &mockBroadcastRepository{
					getByMatchID: func(ctx context.Context, matchID int) (*model.BroadcastMessage, error) {
						return nil, errors.New("not found")
					},
				},
			},
			args: args{matchID: 123},
			want: &dto.APIResponse{Message: "No fans found for this match"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := controller.NewAdminController(tt.fields.championshipAPI, tt.fields.fanRepo, tt.fields.broadcastRepo, tt.fields.broadcast)
			got, err := a.BroadcastMatch(context.Background(), tt.args.matchID)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminController.BroadcastMatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Message != tt.want.Message {
				t.Errorf("AdminController.BroadcastMatch() = %v, want %v", got.Message, tt.want.Message)
			}
		})
	}
}

func TestAdminController_RegisterWS(t *testing.T) {
	broadcastService := broadcast.NewBroadcastService()
	adminController := controller.NewAdminController(nil, nil, nil, broadcastService)
	adminController.RegisterWS(nil)
}

func TestAdminController_UnregisterWS(t *testing.T) {
	broadcastService := broadcast.NewBroadcastService()
	adminController := controller.NewAdminController(nil, nil, nil, broadcastService)
	adminController.UnregisterWS(nil)
}
