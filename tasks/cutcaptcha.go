package tasks

import (
	"fmt"

	"github.com/aarock1234/unicap"
)

// CutCaptchaTask represents a Cutcaptcha solving task.
type CutCaptchaTask struct {
	WebsiteURL string
	MiseryKey  string
	APIKey     string
	Proxy      *unicap.Proxy
}

// Type returns the SDK task type identifier.
func (t *CutCaptchaTask) Type() unicap.TaskType {
	return unicap.TaskTypeCutCaptcha
}

// Validate ensures required fields are present.
func (t *CutCaptchaTask) Validate() error {
	if t.WebsiteURL == "" {
		return fmt.Errorf("website_url: %w", unicap.ErrInvalidTask)
	}

	if t.MiseryKey == "" {
		return fmt.Errorf("misery_key: %w", unicap.ErrInvalidTask)
	}

	if t.APIKey == "" {
		return fmt.Errorf("api_key: %w", unicap.ErrInvalidTask)
	}

	return nil
}
