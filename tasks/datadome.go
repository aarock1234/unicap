package tasks

import (
	"fmt"

	"github.com/aarock1234/unicap"
)

// DataDomeTask represents a DataDome slider captcha solving task.
type DataDomeTask struct {
	WebsiteURL string
	CaptchaURL string
	UserAgent  string
	Proxy      *unicap.Proxy
}

// Type returns the SDK task type identifier.
func (t *DataDomeTask) Type() unicap.TaskType {
	return unicap.TaskTypeDataDome
}

// Validate reports an error if required fields are missing or no proxy is set.
func (t *DataDomeTask) Validate() error {
	if t.WebsiteURL == "" {
		return fmt.Errorf("website_url: %w", unicap.ErrInvalidTask)
	}

	if t.CaptchaURL == "" {
		return fmt.Errorf("captcha_url: %w", unicap.ErrInvalidTask)
	}

	if t.UserAgent == "" {
		return fmt.Errorf("user_agent: %w", unicap.ErrInvalidTask)
	}

	if !t.Proxy.IsSet() {
		return fmt.Errorf("proxy is required for datadome: %w", unicap.ErrInvalidTask)
	}

	return nil
}
