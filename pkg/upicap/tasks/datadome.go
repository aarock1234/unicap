package tasks

import (
	"encoding/json"
	"fmt"

	"github.com/aarock1234/unicap/pkg/upicap"
)

// DataDomeTask represents a DataDome slider captcha solving task
type DataDomeTask struct {
	WebsiteURL string
	CaptchaURL string
	UserAgent  string
	Proxy      *upicap.Proxy
}

func (t *DataDomeTask) Type() upicap.TaskType {
	return upicap.TaskTypeDataDome
}

func (t *DataDomeTask) MarshalJSON() ([]byte, error) {
	return json.Marshal(t)
}

// Validate ensures required fields are present
func (t *DataDomeTask) Validate() error {
	if t.CaptchaURL == "" {
		return fmt.Errorf("captcha_url: %w", upicap.ErrInvalidTask)
	}
	if t.UserAgent == "" {
		return fmt.Errorf("user_agent: %w", upicap.ErrInvalidTask)
	}
	if !t.Proxy.IsSet() {
		return fmt.Errorf("proxy is required for datadome: %w", upicap.ErrInvalidTask)
	}
	return nil
}
