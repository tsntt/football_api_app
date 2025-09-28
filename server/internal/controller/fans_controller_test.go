package controller_test

import (
	"context"
	"errors"
	"testing"

	"github.com/tsntt/footballapi/internal/controller"
	"github.com/tsntt/footballapi/internal/dto"
	"github.com/tsntt/footballapi/internal/model"
)

func BenchmarkFanController_Subscribe(b *testing.B) {
	mockFanRepo := &mockFanRepository{
		create: func(ctx context.Context, fan *model.Fan) error {
			return nil
		},
	}

	fanController := controller.NewFanController(mockFanRepo)
	req := &dto.FanRequest{
		UserID:   1,
		TeamID:   1,
		TeamName: "Test Team",
	}

	for i := 0; i < b.N; i++ {
		_, _ = fanController.Subscribe(context.Background(), req)
	}
}

func BenchmarkFanController_Unsubscribe(b *testing.B) {
	mockFanRepo := &mockFanRepository{
		deleteByUserIDAndTeam: func(ctx context.Context, userID int, team string) error {
			return nil
		},
	}

	fanController := controller.NewFanController(mockFanRepo)
	req := &dto.UnsubscribeRequest{
		TeamID: "1",
	}

	for i := 0; i < b.N; i++ {
		_, _ = fanController.Unsubscribe(context.Background(), 1, req)
	}
}

func BenchmarkFanController_GetSubscriptions(b *testing.B) {
	mockFanRepo := &mockFanRepository{
		getByUserID: func(ctx context.Context, userID int) ([]model.Fan, error) {
			return []model.Fan{{ID: 1, UserID: 1, TeamID: 1}}, nil
		},
	}

	fanController := controller.NewFanController(mockFanRepo)

	for i := 0; i < b.N; i++ {
		_, _ = fanController.GetSubscriptions(context.Background(), 1)
	}
}

func TestFanController_Subscribe(t *testing.T) {
	type fields struct {
		fanRepo model.IFanRepository
	}
	type args struct {
		req *dto.FanRequest
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
				fanRepo: &mockFanRepository{
					create: func(ctx context.Context, fan *model.Fan) error {
						return nil
					},
				},
			},
			args: args{
				req: &dto.FanRequest{
					UserID:   1,
					TeamID:   1,
					TeamName: "Test Team",
				},
			},
			want:    &dto.APIResponse{Message: "Subscribed to Test Team"},
			wantErr: false,
		},
		{
			name: "validation error",
			fields: fields{
				fanRepo: &mockFanRepository{},
			},
			args: args{
				req: &dto.FanRequest{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "create error",
			fields: fields{
				fanRepo: &mockFanRepository{
					create: func(ctx context.Context, fan *model.Fan) error {
						return errors.New("create error")
					},
				},
			},
			args: args{
				req: &dto.FanRequest{
					UserID:   1,
					TeamID:   1,
					TeamName: "Test Team",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := controller.NewFanController(tt.fields.fanRepo)
			got, err := f.Subscribe(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("FanController.Subscribe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil && got.Message != tt.want.Message {
				t.Errorf("FanController.Subscribe() = %v, want %v", got.Message, tt.want.Message)
			}
		})
	}
}
