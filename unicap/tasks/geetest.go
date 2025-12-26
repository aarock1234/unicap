package tasks

import (
	"encoding/json"
	"fmt"

	"github.com/aarock1234/unicap/unicap"
)

// GeeTestTask represents a GeeTest v3 solving task
type GeeTestTask struct {
	WebsiteURL         string
	GT                 string
	Challenge          string
	APIServerSubdomain string
	Proxy              *unicap.Proxy
}

func (t *GeeTestTask) Type() unicap.TaskType {
	return unicap.TaskTypeGeeTest
}

func (t *GeeTestTask) MarshalJSON() ([]byte, error) {
	return json.Marshal(t)
}

// Validate ensures required fields are present
func (t *GeeTestTask) Validate() error {
	if t.WebsiteURL == "" {
		return fmt.Errorf("website_url: %w", unicap.ErrInvalidTask)
	}
	if t.GT == "" {
		return fmt.Errorf("gt: %w", unicap.ErrInvalidTask)
	}
	if t.Challenge == "" {
		return fmt.Errorf("challenge: %w", unicap.ErrInvalidTask)
	}
	return nil
}

// GeeTestV4Task represents a GeeTest v4 solving task
type GeeTestV4Task struct {
	WebsiteURL         string
	CaptchaID          string
	APIServerSubdomain string
	Proxy              *unicap.Proxy
}

func (t *GeeTestV4Task) Type() unicap.TaskType {
	return unicap.TaskTypeGeeTestV4
}

func (t *GeeTestV4Task) MarshalJSON() ([]byte, error) {
	return json.Marshal(t)
}

// Validate ensures required fields are present
func (t *GeeTestV4Task) Validate() error {
	if t.WebsiteURL == "" {
		return fmt.Errorf("website_url: %w", unicap.ErrInvalidTask)
	}
	if t.CaptchaID == "" {
		return fmt.Errorf("captcha_id: %w", unicap.ErrInvalidTask)
	}
	return nil
}
