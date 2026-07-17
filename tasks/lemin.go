package tasks

import (
	"fmt"

	"github.com/aarock1234/unicap"
)

// LeminTask represents a Lemin Cropped captcha solving task.
type LeminTask struct {
	WebsiteURL         string
	CaptchaID          string
	DivID              string
	APIServerSubdomain string
	UserAgent          string
	Proxy              *unicap.Proxy
}

// Type returns the SDK task type identifier.
func (t *LeminTask) Type() unicap.TaskType {
	return unicap.TaskTypeLemin
}

// Validate ensures required fields are present.
func (t *LeminTask) Validate() error {
	if t.WebsiteURL == "" {
		return fmt.Errorf("website_url: %w", unicap.ErrInvalidTask)
	}

	if t.CaptchaID == "" {
		return fmt.Errorf("captcha_id: %w", unicap.ErrInvalidTask)
	}

	if t.DivID == "" {
		return fmt.Errorf("div_id: %w", unicap.ErrInvalidTask)
	}

	return nil
}
