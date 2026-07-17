package tasks

import (
	"fmt"

	"github.com/aarock1234/unicap"
)

// MTCaptchaTask represents an MTCaptcha solving task.
type MTCaptchaTask struct {
	WebsiteURL string
	WebsiteKey string
	Proxy      *unicap.Proxy
}

// Type returns the SDK task type identifier.
func (t *MTCaptchaTask) Type() unicap.TaskType {
	return unicap.TaskTypeMTCaptcha
}

// Validate ensures required fields are present.
func (t *MTCaptchaTask) Validate() error {
	if t.WebsiteURL == "" {
		return fmt.Errorf("website_url: %w", unicap.ErrInvalidTask)
	}

	if t.WebsiteKey == "" {
		return fmt.Errorf("website_key: %w", unicap.ErrInvalidTask)
	}

	return nil
}
