package noptcha

import (
	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

// turnstileParams are the params of Noptcha's turnstile task.
type turnstileParams struct {
	WebsiteURL string `json:"website_url"`
	WebsiteKey string `json:"website_key"`
	Action     string `json:"action,omitempty"`
	CData      string `json:"cdata,omitempty"`
	PageData   string `json:"page_data,omitempty"`
	Proxy      *proxy `json:"proxy,omitempty"`
}

// cloudflareChallengeParams are the params of Noptcha's cloudflare.challenge
// task, which is solved only through the caller's proxy.
type cloudflareChallengeParams struct {
	WebsiteURL string `json:"website_url"`
	UserAgent  string `json:"user_agent,omitempty"`
	Proxy      *proxy `json:"proxy"`
}

// mapTurnstile builds the params of Noptcha's turnstile task. Every task
// field has a counterpart.
func mapTurnstile(task *tasks.TurnstileTask) (turnstileParams, error) {
	return turnstileParams{
		WebsiteURL: task.WebsiteURL,
		WebsiteKey: task.WebsiteKey,
		Action:     task.Action,
		CData:      task.CData,
		PageData:   task.PageData,
		Proxy:      newProxy(task.Proxy),
	}, nil
}

// mapCloudflareChallenge builds the params of Noptcha's cloudflare.challenge
// task. Noptcha has no field for the challenge page's HTML.
func mapCloudflareChallenge(task *tasks.CloudflareChallengeTask) (cloudflareChallengeParams, error) {
	if err := solverapi.Check(
		solverapi.Required(
			solverapi.Setting("Proxy", task.Proxy != nil),
		),
		solverapi.Unsupported(
			solverapi.Setting("HTML", task.HTML != ""),
		),
	); err != nil {
		return cloudflareChallengeParams{}, err
	}

	return cloudflareChallengeParams{
		WebsiteURL: task.WebsiteURL,
		UserAgent:  task.UserAgent,
		Proxy:      newProxy(task.Proxy),
	}, nil
}
