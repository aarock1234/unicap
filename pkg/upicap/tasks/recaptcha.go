package tasks

import (
	"encoding/json"
	"fmt"

	"upicap/pkg/upicap"
)

// ReCaptchaV2Task represents a ReCaptcha V2 solving task
type ReCaptchaV2Task struct {
	WebsiteURL  string
	WebsiteKey  string
	IsInvisible bool
	DataS       string
	PageAction  string
	Proxy       *upicap.Proxy
}

func (t *ReCaptchaV2Task) Type() upicap.TaskType {
	return upicap.TaskTypeReCaptchaV2
}

func (t *ReCaptchaV2Task) MarshalJSON() ([]byte, error) {
	return json.Marshal(t)
}

// Validate ensures required fields are present
func (t *ReCaptchaV2Task) Validate() error {
	if t.WebsiteURL == "" {
		return fmt.Errorf("website_url: %w", upicap.ErrInvalidTask)
	}
	if t.WebsiteKey == "" {
		return fmt.Errorf("website_key: %w", upicap.ErrInvalidTask)
	}
	return nil
}

// ReCaptchaV3Task represents a ReCaptcha V3 solving task
type ReCaptchaV3Task struct {
	WebsiteURL string
	WebsiteKey string
	PageAction string
	MinScore   float64
	Proxy      *upicap.Proxy
}

func (t *ReCaptchaV3Task) Type() upicap.TaskType {
	return upicap.TaskTypeReCaptchaV3
}

func (t *ReCaptchaV3Task) MarshalJSON() ([]byte, error) {
	return json.Marshal(t)
}

// Validate ensures required fields are present
func (t *ReCaptchaV3Task) Validate() error {
	if t.WebsiteURL == "" {
		return fmt.Errorf("website_url: %w", upicap.ErrInvalidTask)
	}
	if t.WebsiteKey == "" {
		return fmt.Errorf("website_key: %w", upicap.ErrInvalidTask)
	}
	if t.MinScore < 0 || t.MinScore > 1 {
		return fmt.Errorf("min_score must be between 0 and 1: %w", upicap.ErrInvalidTask)
	}
	return nil
}

// ReCaptchaV2EnterpriseTask represents a ReCaptcha V2 Enterprise solving task
type ReCaptchaV2EnterpriseTask struct {
	WebsiteURL     string
	WebsiteKey     string
	EnterpriseData map[string]any
	IsInvisible    bool
	ApiDomain      string
	Proxy          *upicap.Proxy
}

func (t *ReCaptchaV2EnterpriseTask) Type() upicap.TaskType {
	return upicap.TaskTypeReCaptchaV2Enterprise
}

func (t *ReCaptchaV2EnterpriseTask) MarshalJSON() ([]byte, error) {
	return json.Marshal(t)
}

// Validate ensures required fields are present
func (t *ReCaptchaV2EnterpriseTask) Validate() error {
	if t.WebsiteURL == "" {
		return fmt.Errorf("website_url: %w", upicap.ErrInvalidTask)
	}
	if t.WebsiteKey == "" {
		return fmt.Errorf("website_key: %w", upicap.ErrInvalidTask)
	}
	return nil
}

// ReCaptchaV3EnterpriseTask represents a ReCaptcha V3 Enterprise solving task
type ReCaptchaV3EnterpriseTask struct {
	WebsiteURL     string
	WebsiteKey     string
	PageAction     string
	EnterpriseData map[string]any
	MinScore       float64
	ApiDomain      string
	Proxy          *upicap.Proxy
}

func (t *ReCaptchaV3EnterpriseTask) Type() upicap.TaskType {
	return upicap.TaskTypeReCaptchaV3Enterprise
}

func (t *ReCaptchaV3EnterpriseTask) MarshalJSON() ([]byte, error) {
	return json.Marshal(t)
}

// Validate ensures required fields are present
func (t *ReCaptchaV3EnterpriseTask) Validate() error {
	if t.WebsiteURL == "" {
		return fmt.Errorf("website_url: %w", upicap.ErrInvalidTask)
	}
	if t.WebsiteKey == "" {
		return fmt.Errorf("website_key: %w", upicap.ErrInvalidTask)
	}
	if t.MinScore < 0 || t.MinScore > 1 {
		return fmt.Errorf("min_score must be between 0 and 1: %w", upicap.ErrInvalidTask)
	}
	return nil
}
