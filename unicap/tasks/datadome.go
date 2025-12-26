package tasks

import (
	"encoding/json"
	"fmt"

	"github.com/aarock1234/unicap/unicap"
)

// DataDomeTask represents a DataDome slider captcha solving task
type DataDomeTask struct {
	WebsiteURL string
	CaptchaURL string
	UserAgent  string
	Proxy      *unicap.Proxy
}

func (t *DataDomeTask) Type() unicap.TaskType {
	return unicap.TaskTypeDataDome
}

func (t *DataDomeTask) MarshalJSON() ([]byte, error) {
	return json.Marshal(t)
}

// Validate ensures required fields are present
func (t *DataDomeTask) Validate() error {
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
