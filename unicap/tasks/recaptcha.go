package tasks

import (
	"encoding/json"
	"fmt"

	"github.com/aarock1234/unicap/unicap"
)

// ReCaptchaV2Task represents a ReCaptcha V2 solving task
type ReCaptchaV2Task struct {
	WebsiteURL        string
	WebsiteKey        string
	IsInvisible       bool
	PageAction        string
	DataS             string
	EnterprisePayload map[string]any
	IsSession         bool
	APIDomain         string
	UserAgent         string
	Cookies           string
	Proxy             *unicap.Proxy
}

func (t *ReCaptchaV2Task) Type() unicap.TaskType {
	return unicap.TaskTypeReCaptchaV2
}

func (t *ReCaptchaV2Task) MarshalJSON() ([]byte, error) {
	return json.Marshal(t)
}

// Validate ensures required fields are present
func (t *ReCaptchaV2Task) Validate() error {
	if t.WebsiteURL == "" {
		return fmt.Errorf("website_url: %w", unicap.ErrInvalidTask)
	}
	if t.WebsiteKey == "" {
		return fmt.Errorf("website_key: %w", unicap.ErrInvalidTask)
	}
	return nil
}

// ReCaptchaV3Task represents a ReCaptcha V3 solving task
type ReCaptchaV3Task struct {
	WebsiteURL        string
	WebsiteKey        string
	PageAction        string
	MinScore          float64
	EnterprisePayload map[string]any
	IsSession         bool
	APIDomain         string
	IsEnterprise      bool
	Proxy             *unicap.Proxy
}

func (t *ReCaptchaV3Task) Type() unicap.TaskType {
	return unicap.TaskTypeReCaptchaV3
}

func (t *ReCaptchaV3Task) MarshalJSON() ([]byte, error) {
	return json.Marshal(t)
}

// Validate ensures required fields are present
func (t *ReCaptchaV3Task) Validate() error {
	if t.WebsiteURL == "" {
		return fmt.Errorf("website_url: %w", unicap.ErrInvalidTask)
	}
	if t.WebsiteKey == "" {
		return fmt.Errorf("website_key: %w", unicap.ErrInvalidTask)
	}
	if t.MinScore < 0 || t.MinScore > 1 {
		return fmt.Errorf("min_score must be between 0 and 1: %w", unicap.ErrInvalidTask)
	}
	return nil
}

// ReCaptchaV2EnterpriseTask represents a ReCaptcha V2 Enterprise solving task
type ReCaptchaV2EnterpriseTask struct {
	WebsiteURL        string
	WebsiteKey        string
	IsInvisible       bool
	PageAction        string
	DataS             string
	EnterprisePayload map[string]any
	IsSession         bool
	APIDomain         string
	Proxy             *unicap.Proxy
}

func (t *ReCaptchaV2EnterpriseTask) Type() unicap.TaskType {
	return unicap.TaskTypeReCaptchaV2Enterprise
}

func (t *ReCaptchaV2EnterpriseTask) MarshalJSON() ([]byte, error) {
	return json.Marshal(t)
}

// Validate ensures required fields are present
func (t *ReCaptchaV2EnterpriseTask) Validate() error {
	if t.WebsiteURL == "" {
		return fmt.Errorf("website_url: %w", unicap.ErrInvalidTask)
	}
	if t.WebsiteKey == "" {
		return fmt.Errorf("website_key: %w", unicap.ErrInvalidTask)
	}
	return nil
}

// ReCaptchaV3EnterpriseTask represents a ReCaptcha V3 Enterprise solving task
type ReCaptchaV3EnterpriseTask struct {
	WebsiteURL        string
	WebsiteKey        string
	PageAction        string
	MinScore          float64
	EnterprisePayload map[string]any
	IsSession         bool
	APIDomain         string
	Proxy             *unicap.Proxy
}

func (t *ReCaptchaV3EnterpriseTask) Type() unicap.TaskType {
	return unicap.TaskTypeReCaptchaV3Enterprise
}

func (t *ReCaptchaV3EnterpriseTask) MarshalJSON() ([]byte, error) {
	return json.Marshal(t)
}

// Validate ensures required fields are present
func (t *ReCaptchaV3EnterpriseTask) Validate() error {
	if t.WebsiteURL == "" {
		return fmt.Errorf("website_url: %w", unicap.ErrInvalidTask)
	}
	if t.WebsiteKey == "" {
		return fmt.Errorf("website_key: %w", unicap.ErrInvalidTask)
	}
	if t.MinScore < 0 || t.MinScore > 1 {
		return fmt.Errorf("min_score must be between 0 and 1: %w", unicap.ErrInvalidTask)
	}
	return nil
}
