package broadcaster

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

type BroadcastService struct {
	// Communication channels between goroutines
	notificationChan chan NotificationJob
	workers          int
	wg               sync.WaitGroup
}

type NotificationJob struct {
	MatchID int
	Message string
	to      []interface{}
	Done    chan error
}

type NotificationResult struct {
	UserID int
	Status string
	Error  error
}

func NewBroadcastService(workers int) *BroadcastService {
	service := &BroadcastService{
		notificationChan: make(chan NotificationJob, 100),
		workers:          workers,
	}

	service.startWorkers()

	return service
}

func (s *BroadcastService) startWorkers() {
	for i := 0; i < s.workers; i++ {
		s.wg.Add(1)
		go s.worker(i)
	}
}

func (s *BroadcastService) worker(id int) {
	defer s.wg.Done()
	log.Printf("Broadcast worker %d started", id)

	for job := range s.notificationChan {
		err := s.processNotification(job)

		// Send result back
		select {
		case job.Done <- err:
		case <-time.After(5 * time.Second):
			log.Printf("Worker %d: timeout sending result for match %d", id, job.MatchID)
		}
	}

	log.Printf("Broadcast worker %d stopped", id)
}

func (s *BroadcastService) processNotification(job NotificationJob) error {
	log.Printf("Processing notification for match %d to %d fans", job.MatchID, len(job.to))

	var wg sync.WaitGroup
	results := make(chan NotificationResult, len(job.to))

	// Send notifications in parallel
	for _, fan := range job.to {
		wg.Add(1)
		go func(f interface{}) {
			defer wg.Done()

			err := s.sendPushNotification(f.UserID, job.Message)
			if err != nil {
				// Fallback
				err = s.sendWebSocketMessage(f.UserID, job.Message)
			}

			status := "sent"
			if err != nil {
				status = "failed"
				log.Printf("Failed to send notification to user %d: %v", f.UserID, err)
			}

			results <- NotificationResult{
				UserID: f.UserID,
				Status: status,
				Error:  err,
			}
		}(fan)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var errors []error
	successCount := 0

	for result := range results {
		if result.Error != nil {
			errors = append(errors, result.Error)
		} else {
			successCount++
		}
	}

	log.Printf("Notification for match %d completed: %d success, %d errors",
		job.MatchID, successCount, len(errors))

	if len(errors) > 0 {
		return fmt.Errorf("notification completed with %d errors", len(errors))
	}

	return nil
}

func (s *BroadcastService) SendNotification(ctx context.Context, matchID int, message string, to []interface{}) error {
	if len(to) == 0 {
		return nil
	}

	job := NotificationJob{
		MatchID: matchID,
		Message: message,
		to:      to,
		Done:    make(chan error, 1),
	}

	select {
	case s.notificationChan <- job:
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(10 * time.Second):
		return fmt.Errorf("timeout queuing notification job")
	}

	select {
	case err := <-job.Done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(30 * time.Second):
		return fmt.Errorf("timeout waiting for notification completion")
	}
}

// Push notifications are simulated in this example
func (s *BroadcastService) sendPushNotification(userID int, message string) error {
	log.Printf("Sending push notification to user %d: %s", userID, message)

	// Simulate possible failure (10% chance)
	if time.Now().UnixNano()%10 == 0 {
		return fmt.Errorf("push notification service unavailable")
	}

	// latency simulation
	time.Sleep(100 * time.Millisecond)

	return nil
}

// TODO: implement websockets
func (s *BroadcastService) sendWebSocketMessage(userID int, message string) error {
	log.Printf("Sending WebSocket message to user %d: %s", userID, message)

	if time.Now().UnixNano()%20 == 0 {
		return fmt.Errorf("websocket connection not available")
	}

	time.Sleep(50 * time.Millisecond)

	return nil
}

func (s *BroadcastService) Stop() {
	close(s.notificationChan)
	s.wg.Wait()
	log.Println("Broadcast service stopped")
}
