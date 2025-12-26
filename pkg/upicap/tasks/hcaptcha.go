package tasks

import (
	"encoding/json"
	"fmt"

	"github.com/aarock1234/unicap/pkg/unicap"
)

// HCaptchaTask represents an hCaptcha solving task
type HCaptchaTask struct {
	WebsiteURL     string
	WebsiteKey     string
	IsInvisible    bool
	EnterpriseData map[string]any
	UserAgent      string
	Cookies        string
	Proxy          *unicap.Proxy
}

func (t *HCaptchaTask) Type() unicap.TaskType {
	return unicap.TaskTypeHCaptcha
}

func (t *HCaptchaTask) MarshalJSON() ([]byte, error) {
	return json.Marshal(t)
}

// Validate ensures required fields are present
func (t *HCaptchaTask) Validate() error {
	if t.WebsiteURL == "" {
		return fmt.Errorf("website_url: %w", unicap.ErrInvalidTask)
	}
	if t.WebsiteKey == "" {
		return fmt.Errorf("website_key: %w", unicap.ErrInvalidTask)
	}
	return nil
}
