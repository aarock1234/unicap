package nopecha

import (
	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

// turnstileRequest is the body NopeCHA's Turnstile token endpoint takes. The
// proxy is required: the token is bound to the IP address it was solved
// from, so the caller submits it through the same proxy.
type turnstileRequest struct {
	Sitekey string        `json:"sitekey"`
	URL     string        `json:"url"`
	Data    turnstileData `json:"data,omitzero"`
	Proxy   *proxy        `json:"proxy"`
}

// turnstileData holds the widget parameters of a Turnstile job.
type turnstileData struct {
	Action string `json:"action,omitempty"`
	CData  string `json:"cdata,omitempty"`
}

// mapTurnstile builds the body of NopeCHA's Turnstile token endpoint.
// NopeCHA has no field for a challenge page's data.
func mapTurnstile(task *tasks.TurnstileTask) (turnstileRequest, error) {
	if err := solverapi.Check(
		solverapi.Required(
			solverapi.Setting("Proxy", task.Proxy != nil),
		),
		solverapi.Unsupported(
			solverapi.Setting("PageData", task.PageData != ""),
		),
		checkProxy(task.Proxy),
	); err != nil {
		return turnstileRequest{}, err
	}

	return turnstileRequest{
		Sitekey: task.WebsiteKey,
		URL:     task.WebsiteURL,
		Data: turnstileData{
			Action: task.Action,
			CData:  task.CData,
		},
		Proxy: newProxy(task.Proxy),
	}, nil
}
