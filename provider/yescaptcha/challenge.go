package yescaptcha

import (
	"slices"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

// turnstileTask is YesCaptcha's Cloudflare Turnstile task. The proxy fields
// belong to the proxy wire type only.
type turnstileTask struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	WebsiteKey string `json:"websiteKey"`
	PageAction string `json:"pageAction,omitempty"`
	Data       string `json:"data,omitempty"`
	solverapi.ProxyFields
}

// cloudflareChallengeTask is YesCaptcha's Cloudflare challenge page task. It
// takes the proxy as a single URL string.
type cloudflareChallengeTask struct {
	Type            string   `json:"type"`
	WebsiteURL      string   `json:"websiteURL"`
	Proxy           string   `json:"proxy"`
	RequiredCookies []string `json:"requiredCookies"`
}

// challengeCookies are the cookies a Cloudflare challenge task must obtain
// to count as solved. Without the requirement YesCaptcha can report, and
// charge for, a task that did not produce the clearance cookie.
var challengeCookies = []string{
	"cf_clearance",
}

// challengeProxyTypes are the proxy protocols the Cloudflare challenge task
// accepts in its proxy URL.
var challengeProxyTypes = []unicap.ProxyType{
	unicap.ProxyTypeHTTP,
	unicap.ProxyTypeSOCKS5,
}

// mapTurnstile builds YesCaptcha's Cloudflare Turnstile task.
func mapTurnstile(task *tasks.TurnstileTask) (turnstileTask, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("PageData", task.PageData != ""),
	); err != nil {
		return turnstileTask{}, err
	}

	taskType, err := variant("TurnstileTask").For(task.Proxy)
	if err != nil {
		return turnstileTask{}, err
	}

	return turnstileTask{
		Type:        taskType,
		WebsiteURL:  task.WebsiteURL,
		WebsiteKey:  task.WebsiteKey,
		PageAction:  task.Action,
		Data:        task.CData,
		ProxyFields: solverapi.NewProxyFields(task.Proxy),
	}, nil
}

// mapCloudflareChallenge builds YesCaptcha's Cloudflare challenge page task.
// UserAgent is rejected because YesCaptcha silently ignores every value that
// is not a desktop Chrome on Linux.
func mapCloudflareChallenge(task *tasks.CloudflareChallengeTask) (cloudflareChallengeTask, error) {
	if err := solverapi.Check(
		solverapi.Unsupported(
			solverapi.Setting("HTML", task.HTML != ""),
			solverapi.Setting("UserAgent", task.UserAgent != ""),
		),
		unsupportedChallengeProxy(task.Proxy),
	); err != nil {
		return cloudflareChallengeTask{}, err
	}

	taskType, err := proxyOnly("CloudFlareTaskS3").For(task.Proxy)
	if err != nil {
		return cloudflareChallengeTask{}, err
	}

	return cloudflareChallengeTask{
		Type:            taskType,
		WebsiteURL:      task.WebsiteURL,
		Proxy:           solverapi.ProxyURL(task.Proxy),
		RequiredCookies: challengeCookies,
	}, nil
}

// unsupportedChallengeProxy reports a proxy the Cloudflare challenge task
// rejects: a protocol outside challengeProxyTypes, or a SOCKS5 proxy with
// credentials. A nil proxy passes; the wire type lookup reports it as
// missing.
func unsupportedChallengeProxy(proxy *unicap.Proxy) error {
	if proxy == nil {
		return nil
	}

	isSOCKS5WithCredentials := proxy.Type == unicap.ProxyTypeSOCKS5 && proxy.HasCredentials()

	return solverapi.Unsupported(
		solverapi.Setting("Proxy.Type "+string(proxy.Type), !slices.Contains(challengeProxyTypes, proxy.Type)),
		solverapi.Setting("Proxy credentials with socks5", isSOCKS5WithCredentials),
	)
}
