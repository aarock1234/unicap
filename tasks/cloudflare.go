package tasks

import (
	"fmt"

	"github.com/aarock1234/unicap"
)

// CloudflareChallengeTask represents a Cloudflare challenge solving task. This
// is for the "Just a moment" challenge page, not Turnstile.
type CloudflareChallengeTask struct {
	WebsiteURL string
	HTML       string
	UserAgent  string
	Proxy      *unicap.Proxy
}

// Type returns the SDK task type identifier.
func (t *CloudflareChallengeTask) Type() unicap.TaskType {
	return unicap.TaskTypeCloudflareChallenge
}

// Validate reports an error if required fields are missing or no proxy is set.
func (t *CloudflareChallengeTask) Validate() error {
	if t.WebsiteURL == "" {
		return fmt.Errorf("website_url: %w", unicap.ErrInvalidTask)
	}

	if !t.Proxy.IsSet() {
		return fmt.Errorf("proxy is required for cloudflare challenge: %w", unicap.ErrInvalidTask)
	}

	return nil
}
