package tasks

import (
	"fmt"

	"github.com/aarock1234/unicap"
)

// AWSWAFTask represents an AWS WAF (Amazon) captcha solving task.
type AWSWAFTask struct {
	WebsiteURL      string
	Key             string
	IV              string
	Context         string
	ChallengeScript string
	CaptchaScript   string
	Proxy           *unicap.Proxy
}

// Type returns the SDK task type identifier.
func (t *AWSWAFTask) Type() unicap.TaskType {
	return unicap.TaskTypeAWSWAF
}

// Validate ensures required fields are present.
func (t *AWSWAFTask) Validate() error {
	if t.WebsiteURL == "" {
		return fmt.Errorf("website_url: %w", unicap.ErrInvalidTask)
	}

	if t.Key == "" {
		return fmt.Errorf("key: %w", unicap.ErrInvalidTask)
	}

	return nil
}
