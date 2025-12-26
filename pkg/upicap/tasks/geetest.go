package tasks

import (
	"encoding/json"
	"fmt"

	"github.com/aarock1234/unicap/pkg/upicap"
)

// GeeTestTask represents a GeeTest v3 solving task
type GeeTestTask struct {
	WebsiteURL         string
	GT                 string
	Challenge          string
	APIServerSubdomain string
	Proxy              *upicap.Proxy
}

func (t *GeeTestTask) Type() upicap.TaskType {
	return upicap.TaskTypeGeeTest
}

func (t *GeeTestTask) MarshalJSON() ([]byte, error) {
	return json.Marshal(t)
}

// Validate ensures required fields are present
func (t *GeeTestTask) Validate() error {
	if t.WebsiteURL == "" {
		return fmt.Errorf("website_url: %w", upicap.ErrInvalidTask)
	}
	if t.GT == "" {
		return fmt.Errorf("gt: %w", upicap.ErrInvalidTask)
	}
	if t.Challenge == "" {
		return fmt.Errorf("challenge: %w", upicap.ErrInvalidTask)
	}
	return nil
}

// GeeTestV4Task represents a GeeTest v4 solving task
type GeeTestV4Task struct {
	WebsiteURL         string
	CaptchaID          string
	APIServerSubdomain string
	Proxy              *upicap.Proxy
}

func (t *GeeTestV4Task) Type() upicap.TaskType {
	return upicap.TaskTypeGeeTestV4
}

func (t *GeeTestV4Task) MarshalJSON() ([]byte, error) {
	return json.Marshal(t)
}

// Validate ensures required fields are present
func (t *GeeTestV4Task) Validate() error {
	if t.WebsiteURL == "" {
		return fmt.Errorf("website_url: %w", upicap.ErrInvalidTask)
	}
	if t.CaptchaID == "" {
		return fmt.Errorf("captcha_id: %w", upicap.ErrInvalidTask)
	}
	return nil
}
