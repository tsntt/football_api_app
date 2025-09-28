package broadcast

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/gorilla/websocket"
)

type BroadcastJob struct {
	Subscription Subscription
	Message      Message
}

type BroadcastResult struct {
	Success bool
	Error   error
}

// INFO: in production you should use an strategy to avoid duplicate messages
type BroadcastService struct {
	notifiers     map[NotificationType]IBroadcaster
	subscriptions map[int][]Subscription
	admConn       map[*websocket.Conn]bool
	mu            sync.RWMutex
}

func NewBroadcastService() *BroadcastService {
	return &BroadcastService{
		notifiers:     make(map[NotificationType]IBroadcaster),
		subscriptions: make(map[int][]Subscription),
		admConn:       make(map[*websocket.Conn]bool),
		mu:            sync.RWMutex{},
	}
}

func (s *BroadcastService) RegisterNotifier(nt NotificationType, notifier IBroadcaster) {
	s.notifiers[nt] = notifier
}

func (s *BroadcastService) AddSubscription(subscription Subscription) {
	s.subscriptions[subscription.ChannelID] = append(s.subscriptions[subscription.ChannelID], subscription)
}

func (s *BroadcastService) RegisterAdmConn(conn *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.admConn[conn] = true
}

func (s *BroadcastService) UnregisterAdmConn(conn *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.admConn, conn)
}

func (s *BroadcastService) BroadCastToChannel(ctx context.Context, nWorkers, channelID int, msg Message) {
	subs, ok := s.subscriptions[channelID]
	if !ok || len(subs) == 0 {
		slog.Info("No subscriptions found for channel", slog.Int("channel_id", channelID))
		return
	}

	totalJobs := len(subs)
	jobs := make(chan BroadcastJob, totalJobs)
	results := make(chan BroadcastResult, totalJobs)

	for w := 1; w <= nWorkers; w++ {
		go s.worker(w, jobs, results)
	}

	for _, sub := range subs {
		jobs <- BroadcastJob{Subscription: sub, Message: msg}
	}
	close(jobs)

	go s.aggregateResults(channelID, totalJobs, results)
}

func (s *BroadcastService) worker(id int, jobs <-chan BroadcastJob, results chan<- BroadcastResult) {
	for job := range jobs {
		slog.Info("Worker", slog.Int("id", id), "processing job", slog.Int("channel_id", job.Subscription.ChannelID))

		broadcaster, ok := s.notifiers[job.Subscription.NotificationType]
		if !ok {
			results <- BroadcastResult{
				Success: false,
				Error:   fmt.Errorf("no notifier found for %s", job.Subscription.NotificationType),
			}
			continue
		}

		err := broadcaster.Send(context.Background(), job.Subscription, job.Message)
		if err != nil {
			results <- BroadcastResult{
				Success: false,
				Error:   err,
			}
			continue
		}
		results <- BroadcastResult{
			Success: true,
			Error:   nil,
		}
	}
}

func (s *BroadcastService) aggregateResults(channelID int, totalJobs int, results <-chan BroadcastResult) {
	status := BroadcastStatus{
		ChannelID:    channelID,
		TotalToSend:  totalJobs,
		ErrorDetails: []string{},
	}

	for i := 0; i < totalJobs; i++ {
		result := <-results
		if result.Success {
			status.SentCount++
		} else {
			status.FailedCount++
			if result.Error != nil {
				status.ErrorDetails = append(status.ErrorDetails, result.Error.Error())
			}
		}

		s.broadcastStatusToAdmins(status)
	}

	status.IsCompleted = true
	s.broadcastStatusToAdmins(status)
	slog.Info(fmt.Sprintf("Broadcast completed for channel %d\n with %d sent, %d failed\n", channelID, status.SentCount, status.FailedCount))
}

func (s *BroadcastService) broadcastStatusToAdmins(status BroadcastStatus) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for conn := range s.admConn {
		if err := conn.WriteJSON(status); err != nil {
			slog.Warn("Failed to send broadcast status to admin", slog.String("err", err.Error()))
			conn.Close()
			go s.UnregisterAdmConn(conn)
		}
	}
}
