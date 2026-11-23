package nonecap

import (
	"fmt"
	"maps"
	"slices"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

// NoneCap's solve types.
const (
	typeHCaptcha           = "hcaptcha"
	typeHCaptchaEnterprise = "hcaptcha_enterprise"
)

// rqdataKey is the only tasks.HCaptchaTask.EnterprisePayload key NoneCap has
// a field for.
const rqdataKey = "rqdata"

// solveRequest is the create-solve body of both hCaptcha solve types. Proxy
// is a proxy URL, credentials included.
type solveRequest struct {
	Type      string `json:"type"`
	Sitekey   string `json:"sitekey"`
	URL       string `json:"url"`
	RQData    string `json:"rqdata,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	Proxy     string `json:"proxy,omitempty"`
}

// mapHCaptcha builds NoneCap's hCaptcha solve. IsInvisible is accepted and
// not sent: NoneCap solves invisible sitekeys without a flag, so the option
// is honoured as it stands.
func mapHCaptcha(task *tasks.HCaptchaTask) (solveRequest, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("Cookies", task.Cookies != ""),
	); err != nil {
		return solveRequest{}, err
	}

	rqdata, err := enterpriseRQData(task.EnterprisePayload)
	if err != nil {
		return solveRequest{}, err
	}

	return solveRequest{
		Type:      solveType(rqdata),
		Sitekey:   task.WebsiteKey,
		URL:       task.WebsiteURL,
		RQData:    rqdata,
		UserAgent: task.UserAgent,
		Proxy:     proxyURL(task.Proxy),
	}, nil
}

// solveType chooses the solve type. unicap has no enterprise flag, so a
// challenge that comes with rqdata is taken to be hCaptcha Enterprise.
func solveType(rqdata string) string {
	if rqdata != "" {
		return typeHCaptchaEnterprise
	}

	return typeHCaptcha
}

// enterpriseRQData returns the rqdata string of an hCaptcha Enterprise
// payload. Every other key, and an rqdata that is not a string, is reported
// as unsupported, in key order.
func enterpriseRQData(payload map[string]any) (string, error) {
	var (
		rqdata      string
		unsupported []solverapi.Field
	)

	for _, key := range slices.Sorted(maps.Keys(payload)) {
		value, isString := payload[key].(string)
		if key == rqdataKey && isString {
			rqdata = value

			continue
		}

		unsupported = append(unsupported, solverapi.Setting(fmt.Sprintf("EnterprisePayload[%s]", key), true))
	}

	if err := solverapi.Unsupported(unsupported...); err != nil {
		return "", err
	}

	return rqdata, nil
}

// proxyURL converts an optional proxy to the URL string NoneCap accepts,
// such as "socks5://user:pass@host:port". A nil proxy yields an empty
// string, which leaves the proxy out of the body.
func proxyURL(proxy *unicap.Proxy) string {
	if proxy == nil {
		return ""
	}

	return proxy.URL().String()
}
