package tasks

import (
	"encoding/json"
	"fmt"

	"upicap/pkg/upicap"
)

// HCaptchaTask represents an hCaptcha solving task
type HCaptchaTask struct {
	WebsiteURL     string
	WebsiteKey     string
	IsInvisible    bool
	EnterpriseData map[string]any
	Proxy          *upicap.Proxy
}

func (t *HCaptchaTask) Type() upicap.TaskType {
	return upicap.TaskTypeHCaptcha
}

func (t *HCaptchaTask) MarshalJSON() ([]byte, error) {
	return json.Marshal(t)
}

// Validate ensures required fields are present
func (t *HCaptchaTask) Validate() error {
	if t.WebsiteURL == "" {
		return fmt.Errorf("website_url: %w", upicap.ErrInvalidTask)
	}
	if t.WebsiteKey == "" {
		return fmt.Errorf("website_key: %w", upicap.ErrInvalidTask)
	}
	return nil
}
