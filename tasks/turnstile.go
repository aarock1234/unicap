package tasks

import (
	"fmt"

	"github.com/aarock1234/unicap"
)

// TurnstileTask represents a Cloudflare Turnstile solving task
type TurnstileTask struct {
	WebsiteURL string
	WebsiteKey string
	Action     string
	CData      string
	PageData   string
	Proxy      *unicap.Proxy
}

// Type returns the SDK task type identifier.
func (t *TurnstileTask) Type() unicap.TaskType {
	return unicap.TaskTypeTurnstile
}

// Validate ensures required fields are present
func (t *TurnstileTask) Validate() error {
	if t.WebsiteURL == "" {
		return fmt.Errorf("website_url: %w", unicap.ErrInvalidTask)
	}
	if t.WebsiteKey == "" {
		return fmt.Errorf("website_key: %w", unicap.ErrInvalidTask)
	}
	return nil
}
