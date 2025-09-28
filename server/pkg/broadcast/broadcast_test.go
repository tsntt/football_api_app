package broadcast_test

import (
	"context"
	"testing"

	"github.com/tsntt/footballapi/pkg/broadcast"
)

// MockBroadcaster is a mock implementation of the IBroadcaster interface for testing purposes.
type MockBroadcaster struct {
	send func(ctx context.Context, subscription broadcast.Subscription, message broadcast.Message) error
}

// Send calls the underlying SendFunc.
func (m *MockBroadcaster) Send(ctx context.Context, subscription broadcast.Subscription, message broadcast.Message) error {
	if m.send != nil {
		return m.send(ctx, subscription, message)
	}
	return nil
}

func TestBroadcastService_RegisterNotifier(t *testing.T) {
	service := broadcast.NewBroadcastService()
	mockNotifier := &MockBroadcaster{}

	service.RegisterNotifier(broadcast.Email, mockNotifier)

	// Use reflection or a dedicated method to inspect internal state if necessary.
	// For this example, we'll assume the registration is successful if no panic occurs.
}

func TestBroadcastService_AddSubscription(t *testing.T) {
	service := broadcast.NewBroadcastService()
	sub := broadcast.Subscription{
		ChannelID:        1,
		NotificationType: broadcast.Email,
	}

	service.AddSubscription(sub)

	// Again, without access to internal state, we assume success if no panic.
	// A more thorough test would involve checking if the subscription was actually added.
}

func TestBroadcastService_BroadCastToChannel_NoSubscriptions(t *testing.T) {
	service := broadcast.NewBroadcastService()
	msg := broadcast.Message{Content: "Test message"}

	// This should not panic or block.
	service.BroadCastToChannel(context.Background(), 2, 1, msg)
}

func TestBroadcastService_BroadCastToChannel_SuccessfulBroadcast(t *testing.T) {
	service := broadcast.NewBroadcastService()
	mockNotifier := &MockBroadcaster{
		send: func(ctx context.Context, sub broadcast.Subscription, msg broadcast.Message) error {
			// Assert that the correct message is being sent to the correct subscriber.
			if msg.Content != "Test message" {
				t.Errorf("Expected message body 'Test message', got '%s'", msg.Content)
			}
			return nil
		},
	}

	service.RegisterNotifier(broadcast.Email, mockNotifier)

	sub := broadcast.Subscription{
		ChannelID:        1,
		NotificationType: broadcast.Email,
	}
	service.AddSubscription(sub)

	msg := broadcast.Message{Content: "Test message"}
	service.BroadCastToChannel(context.Background(), 1, 1, msg)

	// In a real-world scenario, you might need to wait for the broadcast to complete.
	// For this example, we assume synchronous or fast execution.
}
