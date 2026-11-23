package noptcha

import (
	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

// hCaptchaParams are the params of Noptcha's hcaptcha task. The enterprise
// payload, such as {"rqdata": "..."}, is passed on as is.
type hCaptchaParams struct {
	WebsiteURL        string         `json:"website_url"`
	WebsiteKey        string         `json:"website_key"`
	IsInvisible       bool           `json:"is_invisible,omitempty"`
	EnterprisePayload map[string]any `json:"enterprise_payload,omitempty"`
	UserAgent         string         `json:"user_agent,omitempty"`
	Proxy             *proxy         `json:"proxy,omitempty"`
}

// mapHCaptcha builds the params of Noptcha's hcaptcha task.
func mapHCaptcha(task *tasks.HCaptchaTask) (hCaptchaParams, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("Cookies", task.Cookies != ""),
	); err != nil {
		return hCaptchaParams{}, err
	}

	return hCaptchaParams{
		WebsiteURL:        task.WebsiteURL,
		WebsiteKey:        task.WebsiteKey,
		IsInvisible:       task.IsInvisible,
		EnterprisePayload: task.EnterprisePayload,
		UserAgent:         task.UserAgent,
		Proxy:             newProxy(task.Proxy),
	}, nil
}
