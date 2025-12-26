package tasks

import (
	"encoding/json"
	"fmt"

	"github.com/aarock1234/unicap/pkg/upicap"
)

// FunCaptchaTask represents a FunCaptcha (Arkose Labs) solving task
type FunCaptchaTask struct {
	WebsiteURL       string
	WebsitePublicKey string
	APIJSSubdomain   string
	Data             string
	UserAgent        string
	Proxy            *upicap.Proxy
}

func (t *FunCaptchaTask) Type() upicap.TaskType {
	return upicap.TaskTypeFunCaptcha
}

func (t *FunCaptchaTask) MarshalJSON() ([]byte, error) {
	return json.Marshal(t)
}

// Validate ensures required fields are present
func (t *FunCaptchaTask) Validate() error {
	if t.WebsiteURL == "" {
		return fmt.Errorf("website_url: %w", upicap.ErrInvalidTask)
	}
	if t.WebsitePublicKey == "" {
		return fmt.Errorf("website_public_key: %w", upicap.ErrInvalidTask)
	}
	return nil
}
