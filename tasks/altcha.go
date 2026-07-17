package tasks

import (
	"fmt"

	"github.com/aarock1234/unicap"
)

// AltchaTask represents an Altcha captcha solving task. Provide exactly one of
// ChallengeURL or ChallengeJSON.
type AltchaTask struct {
	WebsiteURL    string
	ChallengeURL  string
	ChallengeJSON string
	Proxy         *unicap.Proxy
}

// Type returns the SDK task type identifier.
func (t *AltchaTask) Type() unicap.TaskType {
	return unicap.TaskTypeAltcha
}

// Validate ensures required fields are present and exactly one challenge source
// is provided.
func (t *AltchaTask) Validate() error {
	if t.WebsiteURL == "" {
		return fmt.Errorf("website_url: %w", unicap.ErrInvalidTask)
	}

	if t.ChallengeURL == "" && t.ChallengeJSON == "" {
		return fmt.Errorf("challenge_url or challenge_json: %w", unicap.ErrInvalidTask)
	}

	if t.ChallengeURL != "" && t.ChallengeJSON != "" {
		return fmt.Errorf("challenge_url and challenge_json are mutually exclusive: %w", unicap.ErrInvalidTask)
	}

	return nil
}
