package broadcaster

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type BroadcastService struct {
	// Canais para comunicação entre goroutines
	notificationChan chan NotificationJob
	workers          int
	wg               sync.WaitGroup
}

type NotificationJob struct {
	ID      string                     `json:"id"`      // ID único da notificação (opcional)
	Message string                     `json:"message"` // Mensagem a ser enviada
	Targets []NotificationTarget       `json:"targets"` // Alvos da notificação
	Done    chan NotificationJobResult `json:"-"`       // Canal para resultado
}

type NotificationJobResult struct {
	Success bool                 `json:"success"`
	Error   error                `json:"error"`
	Results []NotificationResult `json:"results"`
}

type NotificationResult struct {
	TargetID string `json:"target_id"`
	Status   string `json:"status"` // "sent", "failed", "pending"
	Error    error  `json:"error"`
}

func NewBroadcastService(workers int) *BroadcastService {
	service := &BroadcastService{
		notificationChan: make(chan NotificationJob, 100),
		workers:          workers,
	}

	// Iniciar workers
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
		result := s.processNotification(job)

		// Enviar resultado de volta
		select {
		case job.Done <- result:
		case <-time.After(5 * time.Second):
			log.Printf("Worker %d: timeout sending result for notification %s", id, job.ID)
		}
	}

	log.Printf("Broadcast worker %d stopped", id)
}

func (s *BroadcastService) processNotification(job NotificationJob) NotificationJobResult {
	log.Printf("Processing notification %s to %d targets", job.ID, len(job.Targets))

	var wg sync.WaitGroup
	results := make(chan NotificationResult, len(job.Targets))

	// Send notifications to all targets in parallel
	for _, target := range job.Targets {
		wg.Add(1)
		go func(t NotificationTarget) {
			defer wg.Done()

			var err error

			switch t.Type {
			case "user":
				err = s.sendToUser(t.Address, job.Message)
			case "email":
				err = s.sendEmail(t.Address, job.Message, t.Metadata)
			case "sms":
				err = s.sendSMS(t.Address, job.Message)
			case "webhook":
				err = s.sendWebhook(t.Address, job.Message, t.Metadata)
			default:
				err = fmt.Errorf("unsupported notification type: %s", t.Type)
			}

			status := "sent"
			if err != nil {
				status = "failed"
				log.Printf("Failed to send notification to %s (%s): %v", t.Address, t.Type, err)
			}

			results <- NotificationResult{
				TargetID: t.ID,
				Status:   status,
				Error:    err,
			}
		}(target)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var errors []error
	var allResults []NotificationResult
	successCount := 0

	for result := range results {
		allResults = append(allResults, result)
		if result.Error != nil {
			errors = append(errors, result.Error)
		} else {
			successCount++
		}
	}

	log.Printf("Notification %s completed: %d success, %d errors",
		job.ID, successCount, len(errors))

	success := len(errors) == 0
	var jobError error
	if !success {
		jobError = fmt.Errorf("notification completed with %d errors", len(errors))
	}

	return NotificationJobResult{
		Success: success,
		Error:   jobError,
		Results: allResults,
	}
}

// SendNotification without ID
func (s *BroadcastService) SendNotification(ctx context.Context, message string, targets []NotificationTarget) error {
	notificationID := fmt.Sprintf("notif_%d", time.Now().UnixNano())
	return s.SendNotificationWithID(ctx, notificationID, message, targets)
}

func (s *BroadcastService) SendNotificationWithID(ctx context.Context, notificationID, message string, targets []NotificationTarget) error {
	if len(targets) == 0 {
		return nil
	}

	job := NotificationJob{
		ID:      notificationID,
		Message: message,
		Targets: targets,
		Done:    make(chan NotificationJobResult, 1),
	}

	select {
	case s.notificationChan <- job:
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(10 * time.Second):
		return fmt.Errorf("timeout queuing notification job")
	}

	select {
	case result := <-job.Done:
		return result.Error
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(30 * time.Second):
		return fmt.Errorf("timeout waiting for notification completion")
	}
}

// Push with fallback to websocket
func (s *BroadcastService) sendToUser(userID string, message string) error {
	if err := s.sendPushNotification(userID, message); err != nil {
		return s.sendWebSocketMessage(userID, message)
	}
	return nil
}

func (s *BroadcastService) sendPushNotification(userID string, message string) error {
	// Simular envio de push notification
	// Em uma implementação real, você usaria Firebase, APNs, etc.
	log.Printf("Sending push notification to user %s: %s", userID, message)

	// Simular possível falha (10% de chance)
	if time.Now().UnixNano()%10 == 0 {
		return fmt.Errorf("push notification service unavailable")
	}

	// Simular latência
	time.Sleep(100 * time.Millisecond)
	return nil
}

func (s *BroadcastService) sendWebSocketMessage(userID string, message string) error {
	// Simular envio via WebSocket como fallback
	log.Printf("Sending WebSocket message to user %s: %s", userID, message)

	// Simular possível falha (5% de chance)
	if time.Now().UnixNano()%20 == 0 {
		return fmt.Errorf("websocket connection not available")
	}

	// Simular latência menor que push notification
	time.Sleep(50 * time.Millisecond)
	return nil
}

func (s *BroadcastService) sendEmail(email string, message string, metadata map[string]int) error {
	// Simular envio de email
	log.Printf("Sending email to %s: %s", email, message)

	subject := "Notification from Football"

	log.Printf("Email subject: %s", subject)

	// Simular latência
	time.Sleep(200 * time.Millisecond)
	return nil
}

func (s *BroadcastService) sendSMS(phone string, message string) error {
	// Simular envio de SMS
	log.Printf("Sending SMS to %s: %s", phone, message)

	// Simular possível falha (8% de chance)
	if time.Now().UnixNano()%12 == 0 {
		return fmt.Errorf("SMS service temporarily unavailable")
	}

	// Simular latência
	time.Sleep(150 * time.Millisecond)
	return nil
}

func (s *BroadcastService) sendWebhook(url string, message string, metadata map[string]int) error {
	// INFO: in production the payload should be signed
	log.Printf("Sending webhook to %s: %s", url, message)

	payload := map[string]interface{}{
		"message":   message,
		"metadata":  metadata,
		"timestamp": time.Now(),
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Fatal("Failed to marshal payload:", err)
	}

	_, err = http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))

	return err
}

func (s *BroadcastService) Stop() {
	close(s.notificationChan)
	s.wg.Wait()
	log.Println("Broadcast service stopped")
}
