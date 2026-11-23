package easycaptcha

import (
	"fmt"
	"maps"
	"slices"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

// easycaptcha.to's token task types.
const (
	typeHCaptchaToken           = "hcaptcha_token"
	typeHCaptchaEnterpriseToken = "hcaptcha_enterprise_token"
)

// rqdataKey is the only tasks.HCaptchaTask.EnterprisePayload key
// easycaptcha.to has a field for.
const rqdataKey = "rqdata"

// proxySchemes are the proxy types easycaptcha.to accepts. SOCKS4 is not
// among them.
var proxySchemes = []unicap.ProxyType{
	unicap.ProxyTypeHTTP,
	unicap.ProxyTypeHTTPS,
	unicap.ProxyTypeSOCKS5,
}

// tokenRequest is the create-task body of both hCaptcha token task types.
// RQData is sent by the enterprise type only.
type tokenRequest struct {
	Type    string `json:"type"`
	Sitekey string `json:"sitekey"`
	URL     string `json:"url"`
	RQData  string `json:"rqdata,omitempty"`
	Proxy   *proxy `json:"proxy,omitempty"`
}

// proxy is easycaptcha.to's proxy object. The username and password are
// optional.
type proxy struct {
	Scheme   unicap.ProxyType `json:"scheme"`
	Host     string           `json:"host"`
	Port     int              `json:"port"`
	Username string           `json:"username,omitempty"`
	Password string           `json:"password,omitempty"`
}

// mapHCaptcha builds easycaptcha.to's hCaptcha token task.
func mapHCaptcha(task *tasks.HCaptchaTask) (tokenRequest, error) {
	if err := solverapi.Check(
		solverapi.Unsupported(
			solverapi.Setting("IsInvisible", task.IsInvisible),
			solverapi.Setting("UserAgent", task.UserAgent != ""),
			solverapi.Setting("Cookies", task.Cookies != ""),
		),
		checkProxy(task.Proxy),
	); err != nil {
		return tokenRequest{}, err
	}

	rqdata, err := enterpriseRQData(task.EnterprisePayload)
	if err != nil {
		return tokenRequest{}, err
	}

	return tokenRequest{
		Type:    tokenType(rqdata),
		Sitekey: task.WebsiteKey,
		URL:     task.WebsiteURL,
		RQData:  rqdata,
		Proxy:   newProxy(task.Proxy),
	}, nil
}

// tokenType chooses the task type. unicap has no enterprise flag, so a
// challenge that comes with rqdata is taken to be hCaptcha Enterprise.
func tokenType(rqdata string) string {
	if rqdata != "" {
		return typeHCaptchaEnterpriseToken
	}

	return typeHCaptchaToken
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

// checkProxy reports a proxy of a type easycaptcha.to does not accept as
// unsupported. A nil proxy passes.
func checkProxy(source *unicap.Proxy) error {
	if source == nil {
		return nil
	}

	return solverapi.Unsupported(
		solverapi.Setting("Proxy.Type "+string(source.Type), !slices.Contains(proxySchemes, source.Type)),
	)
}

// newProxy converts an optional proxy to easycaptcha.to's proxy object. A
// nil proxy yields nil, which leaves the proxy out of the body.
func newProxy(source *unicap.Proxy) *proxy {
	if source == nil {
		return nil
	}

	return &proxy{
		Scheme:   source.Type,
		Host:     source.Host,
		Port:     source.Port,
		Username: source.Username,
		Password: source.Password,
	}
}
