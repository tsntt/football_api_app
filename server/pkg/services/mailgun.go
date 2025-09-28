package services

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"time"

	"github.com/mailgun/mailgun-go/v5"
)

type MailgunService struct {
	mg     mailgun.Mailgun
	from   string
	domain string
}

func NewMailgunService(domain, apiKey, from string) *MailgunService {
	mg := mailgun.NewMailgun(apiKey)

	return &MailgunService{
		mg:     mg,
		from:   from,
		domain: domain,
	}
}

func (m *MailgunService) SendEmail(to, subject, message string, metadata map[string]string) error {
	t, err := template.ParseFiles("internals/libs/email/templates/code.html")
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	data := struct {
		Subject string
		Message string
		AppURL  string
		Year    int
	}{
		Subject: subject,
		Message: message,
		AppURL:  m.domain,
		Year:    time.Now().Year(),
	}

	buf := new(bytes.Buffer)
	if err := t.Execute(buf, data); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	mg := m.mg

	msg := mailgun.NewMessage(m.domain, m.from, subject, message, to)
	msg.SetHTML(buf.String())

	if tag, ok := metadata["tag"]; ok {
		msg.AddTag(tag)
	} else {
		msg.AddTag("football-api")
	}

	// Add metadata
	if matchID, ok := metadata["match_id"]; ok {
		msg.AddVariable("match_id", matchID)
	}

	if team, ok := metadata["team"]; ok {
		msg.AddVariable("team", team)
	}

	if trackClicks, ok := metadata["track_clicks"]; ok && trackClicks == "true" {
		msg.SetTracking(true)
	}

	if trackOpens, ok := metadata["track_opens"]; ok && trackOpens == "true" {
		msg.SetTrackingOpens(true)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	resp, err := mg.Send(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to send email via Mailgun: %w", err)
	}

	// Log da resposta (opcional)
	_ = resp // Response message from Mailgun

	return nil
}
