package services

import (
	"fmt"
	"strings"

	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwilioService struct {
	client    *twilio.RestClient
	fromPhone string
	webhook   string // optional
}

func NewTwilioService(accountSID, authToken, fromPhone, webhook string) *TwilioService {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})

	ts := &TwilioService{
		client:    client,
		fromPhone: fromPhone,
		webhook:   "",
	}

	if webhook != "" {
		ts.webhook = webhook
	}

	return ts
}

func (t *TwilioService) SendSMS(to, message string) error {
	if !t.isValidPhoneNumber(to) {
		return fmt.Errorf("invalid phone number format: %s", to)
	}

	formattedMessage := t.formatMessage(message)

	params := &api.CreateMessageParams{}
	params.SetFrom(t.fromPhone)
	params.SetTo(to)
	params.SetBody(formattedMessage)

	// INFO: webhook is optional
	if t.webhook != "" {
		params.SetStatusCallback(t.webhook)
	}

	resp, err := t.client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send SMS via Twilio: %w", err)
	}

	if resp.Status == nil {
		return fmt.Errorf("twilio returned empty status")
	}

	status := *resp.Status
	if status == "failed" || status == "undelivered" {
		errorMessage := "unknown error"
		if resp.ErrorMessage != nil {
			errorMessage = *resp.ErrorMessage
		}
		return fmt.Errorf("twilio SMS failed with status '%s': %s", status, errorMessage)
	}

	// TODO: add slog
	if resp.Sid != nil {
		fmt.Printf("SMS sent successfully. SID: %s, Status: %s\n", *resp.Sid, status)
	}

	return nil
}

func (t *TwilioService) SendBulkSMS(recipients []string, message string) error {
	if len(recipients) == 0 {
		return fmt.Errorf("no recipients provided")
	}

	var errors []error
	successCount := 0

	for _, recipient := range recipients {
		if err := t.SendSMS(recipient, message); err != nil {
			errors = append(errors, fmt.Errorf("failed to send to %s: %w", recipient, err))
		} else {
			successCount++
		}
	}

	if len(errors) > 0 {
		if successCount == 0 {
			return fmt.Errorf("all SMS failed: %v", errors)
		}
		return fmt.Errorf("some SMS failed (%d/%d successful): %v", successCount, len(recipients), errors)
	}

	return nil
}

// TODO: move to utils
// INFO: in production you should also use twilio lookup api to check if phone number is real
func (t *TwilioService) isValidPhoneNumber(phone string) bool {
	phone = strings.TrimSpace(phone)

	if len(phone) < 10 {
		return false
	}

	if !strings.HasPrefix(phone, "+") {
		return false
	}

	phoneDigits := phone[1:]

	if len(phoneDigits) < 7 || len(phoneDigits) > 15 {
		return false
	}

	for _, char := range phoneDigits {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}

func (t *TwilioService) formatMessage(message string) string {
	prefix := "⚽ Football API: "

	maxLength := 160

	hasUnicode := false
	for _, r := range message {
		if r > 127 {
			hasUnicode = true
			break
		}
	}

	if hasUnicode {
		maxLength = 70
	}

	availableLength := maxLength - len(prefix)

	if len(message) <= availableLength {
		return prefix + message
	}

	truncated := message[:availableLength-3]

	lastSpace := strings.LastIndex(truncated, " ")
	if lastSpace > availableLength/2 {
		truncated = truncated[:lastSpace]
	}

	return prefix + truncated + "..."
}
