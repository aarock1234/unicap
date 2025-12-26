package tasks

import (
	"encoding/json"
	"fmt"

	"github.com/aarock1234/unicap/pkg/upicap"
)

// CloudflareChallengeTask represents a Cloudflare Challenge solving task
// This is for the "Just a moment" challenge page, not Turnstile
type CloudflareChallengeTask struct {
	WebsiteURL string
	HTML       string
	UserAgent  string
	Proxy      *upicap.Proxy
}

func (t *CloudflareChallengeTask) Type() upicap.TaskType {
	return upicap.TaskTypeCloudflareChallenge
}

func (t *CloudflareChallengeTask) MarshalJSON() ([]byte, error) {
	return json.Marshal(t)
}

// Validate ensures required fields are present
func (t *CloudflareChallengeTask) Validate() error {
	if t.WebsiteURL == "" {
		return fmt.Errorf("website_url: %w", upicap.ErrInvalidTask)
	}
	if !t.Proxy.IsSet() {
		return fmt.Errorf("proxy is required for cloudflare challenge: %w", upicap.ErrInvalidTask)
	}
	return nil
}
