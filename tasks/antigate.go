package tasks

import (
	"encoding/json"

	"github.com/aarock1234/unicap/v3"
)

// AntiGateTask runs an Anti-Captcha AntiGate template, where a worker follows
// a scripted scenario in a real browser. Decode the result with
// unicap.DecodeSolution[AntiGateSolution]. Variables set to "_WAIT_FOR_IT_"
// can be supplied later with the Anti-Captcha client's PushVariable.
type AntiGateTask struct {
	// WebsiteURL is the page the scenario starts on.
	WebsiteURL string

	// TemplateName is the AntiGate template to run.
	TemplateName string

	// Variables are the template variables. Templates define their own
	// variable names and types, so the map is intentionally dynamic.
	Variables map[string]any

	// DomainsOfInterest lists domains whose cookies and storage the worker
	// should capture.
	DomainsOfInterest []string

	// Proxy is the optional proxy the worker loads the page through.
	// Anti-Captcha accepts HTTPS proxies only for AntiGate.
	Proxy *unicap.Proxy
}

// Type returns unicap.TaskTypeAntiGate.
func (t *AntiGateTask) Type() unicap.TaskType {
	return unicap.TaskTypeAntiGate
}

// Validate reports missing required fields and malformed proxies.
func (t *AntiGateTask) Validate() error {
	return validate(
		required("WebsiteURL", t.WebsiteURL),
		required("TemplateName", t.TemplateName),
		optionalProxy(t.Proxy),
	)
}

// AntiGateSolution is the decoded payload of an AntiGateTask. Templates
// decide which fields are filled.
type AntiGateSolution struct {
	Cookies        map[string]string `json:"cookies"`
	LocalStorage   map[string]string `json:"localStorage"`
	SessionStorage map[string]string `json:"sessionStorage"`
	URL            string            `json:"url"`
	Domain         string            `json:"domain"`
	HTMLsInBase64  []string          `json:"HTMLsInBase64"`
	Screenshots    []string          `json:"screenshots"`

	// Fingerprint, RequestHeaders, and ResponseHeaders have template-defined
	// shapes, so they are kept raw.
	Fingerprint     json.RawMessage `json:"fingerprint"`
	RequestHeaders  json.RawMessage `json:"requestHeaders"`
	ResponseHeaders json.RawMessage `json:"responseHeaders"`
}
