package tasks

import (
	"encoding/json"
	"fmt"

	"github.com/aarock1234/unicap/pkg/unicap"
)

// CloudflareChallengeTask represents a Cloudflare Challenge solving task
// This is for the "Just a moment" challenge page, not Turnstile
type CloudflareChallengeTask struct {
	WebsiteURL string
	HTML       string
	UserAgent  string
	Proxy      *unicap.Proxy
}

func (t *CloudflareChallengeTask) Type() unicap.TaskType {
	return unicap.TaskTypeCloudflareChallenge
}

func (t *CloudflareChallengeTask) MarshalJSON() ([]byte, error) {
	return json.Marshal(t)
}

// Validate ensures required fields are present
func (t *CloudflareChallengeTask) Validate() error {
	if t.WebsiteURL == "" {
		return fmt.Errorf("website_url: %w", unicap.ErrInvalidTask)
	}
	if !t.Proxy.IsSet() {
		return fmt.Errorf("proxy is required for cloudflare challenge: %w", unicap.ErrInvalidTask)
	}
	return nil
}
