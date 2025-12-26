package tasks

import (
	"encoding/json"
	"fmt"

	"upicap/pkg/upicap"
)

// TurnstileTask represents a Cloudflare Turnstile solving task
type TurnstileTask struct {
	WebsiteURL string
	WebsiteKey string
	Action     string
	CData      string
	PageData   string
	Proxy      *upicap.Proxy
}

func (t *TurnstileTask) Type() upicap.TaskType {
	return upicap.TaskTypeTurnstile
}

func (t *TurnstileTask) MarshalJSON() ([]byte, error) {
	return json.Marshal(t)
}

// Validate ensures required fields are present
func (t *TurnstileTask) Validate() error {
	if t.WebsiteURL == "" {
		return fmt.Errorf("website_url: %w", upicap.ErrInvalidTask)
	}
	if t.WebsiteKey == "" {
		return fmt.Errorf("website_key: %w", upicap.ErrInvalidTask)
	}
	return nil
}
