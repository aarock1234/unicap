package capmonster

import (
	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

// mtCaptchaTask is CapMonster Cloud's MTCaptcha task.
type mtCaptchaTask struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	WebsiteKey string `json:"websiteKey"`
	solverapi.ProxyFields
}

// prosopoTask is CapMonster Cloud's Prosopo Procaptcha task.
type prosopoTask struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	WebsiteKey string `json:"websiteKey"`
	solverapi.ProxyFields
}

// yidunTask is CapMonster Cloud's NetEase Yidun task.
type yidunTask struct {
	Type                    string `json:"type"`
	WebsiteURL              string `json:"websiteURL"`
	WebsiteKey              string `json:"websiteKey"`
	UserAgent               string `json:"userAgent,omitempty"`
	YidunGetLib             string `json:"yidunGetLib,omitempty"`
	YidunAPIServerSubdomain string `json:"yidunApiServerSubdomain,omitempty"`
	Challenge               string `json:"challenge,omitempty"`
	HCG                     string `json:"hcg,omitempty"`
	HCT                     int64  `json:"hct,omitempty"`
	solverapi.ProxyFields
}

// binanceTask is CapMonster Cloud's Binance captcha task.
type binanceTask struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	WebsiteKey string `json:"websiteKey"`
	ValidateID string `json:"validateId"`
	UserAgent  string `json:"userAgent,omitempty"`
	solverapi.ProxyFields
}

// friendlyCaptchaTask is CapMonster Cloud's Friendly Captcha task, a
// CustomTask of class "friendly".
type friendlyCaptchaTask struct {
	Type       string                  `json:"type"`
	Class      string                  `json:"class"`
	WebsiteURL string                  `json:"websiteURL"`
	WebsiteKey string                  `json:"websiteKey"`
	Metadata   friendlyCaptchaMetadata `json:"metadata"`
	solverapi.ProxyFields
}

// friendlyCaptchaMetadata is the metadata object of a Friendly Captcha task.
// APIGetLib is the URL of the widget script.
type friendlyCaptchaMetadata struct {
	APIGetLib string `json:"apiGetLib"`
}

// mapMTCaptcha builds CapMonster Cloud's MTCaptcha task.
func mapMTCaptcha(task *tasks.MTCaptchaTask) (mtCaptchaTask, error) {
	proxy, err := proxyFields(task.Proxy)
	if err != nil {
		return mtCaptchaTask{}, err
	}

	return mtCaptchaTask{
		Type:        "MTCaptchaTask",
		WebsiteURL:  task.WebsiteURL,
		WebsiteKey:  task.WebsiteKey,
		ProxyFields: proxy,
	}, nil
}

// mapProsopo builds CapMonster Cloud's Prosopo Procaptcha task.
func mapProsopo(task *tasks.ProsopoTask) (prosopoTask, error) {
	proxy, err := proxyFields(task.Proxy)
	if err != nil {
		return prosopoTask{}, err
	}

	return prosopoTask{
		Type:        "ProsopoTask",
		WebsiteURL:  task.WebsiteURL,
		WebsiteKey:  task.WebsiteKey,
		ProxyFields: proxy,
	}, nil
}

// mapYidun builds CapMonster Cloud's NetEase Yidun task.
func mapYidun(task *tasks.YidunTask) (yidunTask, error) {
	proxy, err := proxyFields(task.Proxy)
	if err != nil {
		return yidunTask{}, err
	}

	return yidunTask{
		Type:                    "YidunTask",
		WebsiteURL:              task.WebsiteURL,
		WebsiteKey:              task.WebsiteKey,
		UserAgent:               task.UserAgent,
		YidunGetLib:             task.GetLibURL,
		YidunAPIServerSubdomain: task.APIServerSubdomain,
		Challenge:               task.Challenge,
		HCG:                     task.HCG,
		HCT:                     task.HCT,
		ProxyFields:             proxy,
	}, nil
}

// mapBinance builds CapMonster Cloud's Binance captcha task.
func mapBinance(task *tasks.BinanceTask) (binanceTask, error) {
	proxy, err := proxyFields(task.Proxy)
	if err != nil {
		return binanceTask{}, err
	}

	return binanceTask{
		Type:        "BinanceTask",
		WebsiteURL:  task.WebsiteURL,
		WebsiteKey:  task.WebsiteKey,
		ValidateID:  task.ValidateID,
		UserAgent:   task.UserAgent,
		ProxyFields: proxy,
	}, nil
}

// mapFriendlyCaptcha builds CapMonster Cloud's Friendly Captcha task. The
// module script URL is required. Version is rejected because the task has no
// field for it: each widget generation is documented with its own script
// URL.
func mapFriendlyCaptcha(task *tasks.FriendlyCaptchaTask) (friendlyCaptchaTask, error) {
	if err := solverapi.Check(
		solverapi.Required(
			solverapi.Setting("ModuleScript", task.ModuleScript != ""),
		),
		solverapi.Unsupported(
			solverapi.Setting("Version", task.Version != ""),
			solverapi.Setting("NomoduleScript", task.NomoduleScript != ""),
		),
	); err != nil {
		return friendlyCaptchaTask{}, err
	}

	proxy, err := proxyFields(task.Proxy)
	if err != nil {
		return friendlyCaptchaTask{}, err
	}

	return friendlyCaptchaTask{
		Type:       customTaskType,
		Class:      "friendly",
		WebsiteURL: task.WebsiteURL,
		WebsiteKey: task.WebsiteKey,
		Metadata: friendlyCaptchaMetadata{
			APIGetLib: task.ModuleScript,
		},
		ProxyFields: proxy,
	}, nil
}
