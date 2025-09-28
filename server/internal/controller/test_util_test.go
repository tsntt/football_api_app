package controller_test

import (
	"context"

	"github.com/tsntt/footballapi/internal/model"
)

// Mocks
type mockChampionshipAPI struct {
	getChampionships func(ctx context.Context) ([]model.Championship, error)
	getMatches       func(ctx context.Context, championshipID int, team, stage string) ([]model.Match, error)
	getMatch         func(ctx context.Context, matchID int) (*model.Match, error)
}

func (m *mockChampionshipAPI) GetChampionships(ctx context.Context) ([]model.Championship, error) {
	return m.getChampionships(ctx)
}

func (m *mockChampionshipAPI) GetMatches(ctx context.Context, championshipID int, team, stage string) ([]model.Match, error) {
	return m.getMatches(ctx, championshipID, team, stage)
}

func (m *mockChampionshipAPI) GetMatch(ctx context.Context, matchID int) (*model.Match, error) {
	return m.getMatch(ctx, matchID)
}

type mockFanRepository struct {
	create                func(ctx context.Context, fan *model.Fan) error
	getAll                func(ctx context.Context) ([]model.Fan, error)
	getByTeamID           func(ctx context.Context, teamID int) ([]model.Fan, error)
	getByUserID           func(ctx context.Context, userID int) ([]model.Fan, error)
	deleteByUserIDAndTeam func(ctx context.Context, userID int, team string) error
}

func (m *mockFanRepository) Create(ctx context.Context, fan *model.Fan) error {
	return m.create(ctx, fan)
}

func (m *mockFanRepository) GetAll(ctx context.Context) ([]model.Fan, error) {
	return m.getAll(ctx)
}

func (m *mockFanRepository) GetByTeamID(ctx context.Context, teamID int) ([]model.Fan, error) {
	return m.getByTeamID(ctx, teamID)
}

func (m *mockFanRepository) GetByUserID(ctx context.Context, userID int) ([]model.Fan, error) {
	return m.getByUserID(ctx, userID)
}

func (m *mockFanRepository) DeleteByUserIDAndTeam(ctx context.Context, userID int, team string) error {
	return m.deleteByUserIDAndTeam(ctx, userID, team)
}

type mockBroadcastRepository struct {
	create       func(ctx context.Context, broadcast *model.BroadcastMessage) error
	getByMatchID func(ctx context.Context, matchID int) (*model.BroadcastMessage, error)
	update       func(ctx context.Context, broadcast *model.BroadcastMessage) error
}

func (m *mockBroadcastRepository) Create(ctx context.Context, broadcast *model.BroadcastMessage) error {
	return m.create(ctx, broadcast)
}

func (m *mockBroadcastRepository) GetByMatchID(ctx context.Context, matchID int) (*model.BroadcastMessage, error) {
	return m.getByMatchID(ctx, matchID)
}

func (m *mockBroadcastRepository) Update(ctx context.Context, broadcast *model.BroadcastMessage) error {
	return m.update(ctx, broadcast)
}

type mockUserRepository struct {
	create    func(ctx context.Context, user *model.User) error
	getByName func(ctx context.Context, name string) (*model.User, error)
	getByID   func(ctx context.Context, id int) (*model.User, error)
}

func (m *mockUserRepository) Create(ctx context.Context, user *model.User) error {
	return m.create(ctx, user)
}

func (m *mockUserRepository) GetByName(ctx context.Context, name string) (*model.User, error) {
	return m.getByName(ctx, name)
}

func (m *mockUserRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	return m.getByID(ctx, id)
}
