package anticaptcha

import (
	"fmt"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

// siteKeyTask is the shape shared by widgets identified by a URL and a site
// key.
type siteKeyTask struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	WebsiteKey string `json:"websiteKey"`
	solverapi.ProxyFields
}

type altchaTask struct {
	Type          string `json:"type"`
	WebsiteURL    string `json:"websiteURL"`
	ChallengeURL  string `json:"challengeURL,omitempty"`
	ChallengeJSON string `json:"challengeJSON,omitempty"`
	solverapi.ProxyFields
}

type antiGateTask struct {
	Type              string         `json:"type"`
	WebsiteURL        string         `json:"websiteURL"`
	TemplateName      string         `json:"templateName"`
	Variables         map[string]any `json:"variables"`
	DomainsOfInterest []string       `json:"domainsOfInterest,omitempty"`
	ProxyAddress      string         `json:"proxyAddress,omitempty"`
	ProxyPort         int            `json:"proxyPort,omitempty"`
	ProxyLogin        string         `json:"proxyLogin,omitempty"`
	ProxyPassword     string         `json:"proxyPassword,omitempty"`
}

func mapFriendlyCaptcha(task *tasks.FriendlyCaptchaTask) (siteKeyTask, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("Version", task.Version != ""),
		solverapi.Setting("ModuleScript", task.ModuleScript != ""),
		solverapi.Setting("NomoduleScript", task.NomoduleScript != ""),
	); err != nil {
		return siteKeyTask{}, err
	}

	taskType, err := variant("FriendlyCaptchaTask").For(task.Proxy)
	if err != nil {
		return siteKeyTask{}, err
	}

	proxy, err := proxyFields(task.Proxy)
	if err != nil {
		return siteKeyTask{}, err
	}

	return siteKeyTask{
		Type:        taskType,
		WebsiteURL:  task.WebsiteURL,
		WebsiteKey:  task.WebsiteKey,
		ProxyFields: proxy,
	}, nil
}

func mapProsopo(task *tasks.ProsopoTask) (siteKeyTask, error) {
	taskType, err := variant("ProsopoTask").For(task.Proxy)
	if err != nil {
		return siteKeyTask{}, err
	}

	proxy, err := proxyFields(task.Proxy)
	if err != nil {
		return siteKeyTask{}, err
	}

	return siteKeyTask{
		Type:        taskType,
		WebsiteURL:  task.WebsiteURL,
		WebsiteKey:  task.WebsiteKey,
		ProxyFields: proxy,
	}, nil
}

func mapAltcha(task *tasks.AltchaTask) (altchaTask, error) {
	taskType, err := variant("AltchaTask").For(task.Proxy)
	if err != nil {
		return altchaTask{}, err
	}

	proxy, err := proxyFields(task.Proxy)
	if err != nil {
		return altchaTask{}, err
	}

	return altchaTask{
		Type:          taskType,
		WebsiteURL:    task.WebsiteURL,
		ChallengeURL:  task.ChallengeURL,
		ChallengeJSON: task.ChallengeJSON,
		ProxyFields:   proxy,
	}, nil
}

// mapAntiGate sends AntiGate's proxy without a proxyType, because AntiGate
// accepts HTTP(S) proxies only.
func mapAntiGate(task *tasks.AntiGateTask) (antiGateTask, error) {
	if err := solverapi.RequireIPProxy(task.Proxy); err != nil {
		return antiGateTask{}, err
	}

	if task.Proxy != nil && task.Proxy.Type != unicap.ProxyTypeHTTP && task.Proxy.Type != unicap.ProxyTypeHTTPS {
		return antiGateTask{}, fmt.Errorf("Proxy.Type %q, AntiGate accepts https only: %w", task.Proxy.Type, unicap.ErrUnsupportedOption)
	}

	variables := task.Variables
	if variables == nil {
		variables = map[string]any{}
	}

	wire := antiGateTask{
		Type:              "AntiGateTask",
		WebsiteURL:        task.WebsiteURL,
		TemplateName:      task.TemplateName,
		Variables:         variables,
		DomainsOfInterest: task.DomainsOfInterest,
	}

	if task.Proxy != nil {
		wire.ProxyAddress = task.Proxy.Host
		wire.ProxyPort = task.Proxy.Port
		wire.ProxyLogin = task.Proxy.Username
		wire.ProxyPassword = task.Proxy.Password
	}

	return wire, nil
}
