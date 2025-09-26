package broadcaster_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/tsntt/footballapi/pkg/broadcaster"
)

func TestBroadcastService_NewBroadcastService(t *testing.T) {
	service := broadcaster.NewBroadcastService(2)
	defer service.Stop()

	if service == nil {
		t.Fatal("expected service to not be nil")
	}
}

func TestBroadcastService_SendNotificationWithID_SingleTarget(t *testing.T) {
	service := broadcaster.NewBroadcastService(1)
	defer service.Stop()

	targets := []broadcaster.NotificationTarget{
		{ID: "1", Type: "user", Address: "user1"},
	}

	err := service.SendNotificationWithID(context.Background(), "test-notif", "hello", targets)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestBroadcastService_SendNotificationWithID_MultipleTargets(t *testing.T) {
	service := broadcaster.NewBroadcastService(2)
	defer service.Stop()

	targets := []broadcaster.NotificationTarget{
		{ID: "1", Type: "user", Address: "user1"},
		{ID: "2", Type: "email", Address: "test@example.com"},
	}

	err := service.SendNotificationWithID(context.Background(), "test-notif", "hello", targets)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestBroadcastService_SendNotificationWithID_ContextCanceled(t *testing.T) {
	service := broadcaster.NewBroadcastService(1)
	defer service.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	targets := []broadcaster.NotificationTarget{
		{ID: "1", Type: "user", Address: "user1"},
	}

	err := service.SendNotificationWithID(ctx, "test-notif", "hello", targets)

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestBroadcastService_SendWebhook(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer wg.Done()
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	service := broadcaster.NewBroadcastService(1)
	defer service.Stop()

	targets := []broadcaster.NotificationTarget{
		{ID: "1", Type: "webhook", Address: server.URL},
	}

	err := service.SendNotificationWithID(context.Background(), "test-webhook", "hello", targets)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for webhook to be called")
	}
}
