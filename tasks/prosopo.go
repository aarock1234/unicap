package tasks

import (
	"fmt"

	"github.com/aarock1234/unicap"
)

// ProsopoTask represents a Prosopo Procaptcha solving task.
type ProsopoTask struct {
	WebsiteURL string
	WebsiteKey string
	Proxy      *unicap.Proxy
}

// Type returns the SDK task type identifier.
func (t *ProsopoTask) Type() unicap.TaskType {
	return unicap.TaskTypeProsopo
}

// Validate ensures required fields are present.
func (t *ProsopoTask) Validate() error {
	if t.WebsiteURL == "" {
		return fmt.Errorf("website_url: %w", unicap.ErrInvalidTask)
	}

	if t.WebsiteKey == "" {
		return fmt.Errorf("website_key: %w", unicap.ErrInvalidTask)
	}

	return nil
}
