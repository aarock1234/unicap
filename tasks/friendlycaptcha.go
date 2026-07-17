package tasks

import (
	"fmt"

	"github.com/aarock1234/unicap"
)

// FriendlyCaptchaTask represents a Friendly Captcha solving task.
type FriendlyCaptchaTask struct {
	WebsiteURL string
	WebsiteKey string
	Proxy      *unicap.Proxy
}

// Type returns the SDK task type identifier.
func (t *FriendlyCaptchaTask) Type() unicap.TaskType {
	return unicap.TaskTypeFriendlyCaptcha
}

// Validate ensures required fields are present.
func (t *FriendlyCaptchaTask) Validate() error {
	if t.WebsiteURL == "" {
		return fmt.Errorf("website_url: %w", unicap.ErrInvalidTask)
	}

	if t.WebsiteKey == "" {
		return fmt.Errorf("website_key: %w", unicap.ErrInvalidTask)
	}

	return nil
}
