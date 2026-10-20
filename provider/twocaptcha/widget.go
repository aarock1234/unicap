package twocaptcha

import (
	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

// siteKeyTask is the shape shared by widgets identified by a URL and a site
// key.
type siteKeyTask struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	WebsiteKey string `json:"websiteKey"`
	UserAgent  string `json:"userAgent,omitempty"`
	solverapi.ProxyFields
}

type friendlyCaptchaTask struct {
	Type           string `json:"type"`
	WebsiteURL     string `json:"websiteURL"`
	WebsiteKey     string `json:"websiteKey"`
	Version        string `json:"version,omitempty"`
	ModuleScript   string `json:"moduleScript,omitempty"`
	NomoduleScript string `json:"nomoduleScript,omitempty"`
	solverapi.ProxyFields
}

type leminTask struct {
	Type                    string `json:"type"`
	WebsiteURL              string `json:"websiteURL"`
	CaptchaID               string `json:"captchaId"`
	DivID                   string `json:"divId"`
	LeminAPIServerSubdomain string `json:"leminApiServerSubdomain,omitempty"`
	UserAgent               string `json:"userAgent,omitempty"`
	solverapi.ProxyFields
}

type cutCaptchaTask struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	MiseryKey  string `json:"miseryKey"`
	APIKey     string `json:"apiKey"`
	solverapi.ProxyFields
}

type altchaTask struct {
	Type          string `json:"type"`
	WebsiteURL    string `json:"websiteURL"`
	ChallengeURL  string `json:"challengeURL,omitempty"`
	ChallengeJSON string `json:"challengeJSON,omitempty"`
	solverapi.ProxyFields
}

type keyCaptchaTask struct {
	Type           string `json:"type"`
	WebsiteURL     string `json:"websiteURL"`
	UserID         string `json:"s_s_c_user_id"`
	SessionID      string `json:"s_s_c_session_id"`
	WebServerSign  string `json:"s_s_c_web_server_sign"`
	WebServerSign2 string `json:"s_s_c_web_server_sign2"`
	solverapi.ProxyFields
}

type cyberSiAraTask struct {
	Type             string `json:"type"`
	WebsiteURL       string `json:"websiteURL"`
	SlideMasterURLID string `json:"SlideMasterUrlId"`
	UserAgent        string `json:"userAgent"`
	solverapi.ProxyFields
}

type atbCaptchaTask struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	AppID      string `json:"appId"`
	APIServer  string `json:"apiServer"`
	solverapi.ProxyFields
}

type tencentTask struct {
	Type          string `json:"type"`
	WebsiteURL    string `json:"websiteURL"`
	AppID         string `json:"appId"`
	CaptchaScript string `json:"captchaScript,omitempty"`
	solverapi.ProxyFields
}

type captchaFoxTask struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	WebsiteKey string `json:"websiteKey"`
	UserAgent  string `json:"userAgent"`
	APIServer  string `json:"apiServer,omitempty"`
	solverapi.ProxyFields
}

type yandexSmartCaptchaTask struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	WebsiteKey string `json:"websiteKey"`
	UserAgent  string `json:"userAgent,omitempty"`
	Cookies    string `json:"cookies,omitempty"`
	solverapi.ProxyFields
}

func mapMTCaptcha(task *tasks.MTCaptchaTask) (siteKeyTask, error) {
	taskType, err := variant("MtCaptchaTask").For(task.Proxy)
	if err != nil {
		return siteKeyTask{}, err
	}

	return siteKeyTask{
		Type:        taskType,
		WebsiteURL:  task.WebsiteURL,
		WebsiteKey:  task.WebsiteKey,
		ProxyFields: solverapi.NewProxyFields(task.Proxy),
	}, nil
}

func mapProsopo(task *tasks.ProsopoTask) (siteKeyTask, error) {
	taskType, err := variant("ProsopoTask").For(task.Proxy)
	if err != nil {
		return siteKeyTask{}, err
	}

	return siteKeyTask{
		Type:        taskType,
		WebsiteURL:  task.WebsiteURL,
		WebsiteKey:  task.WebsiteKey,
		ProxyFields: solverapi.NewProxyFields(task.Proxy),
	}, nil
}

func mapCapy(task *tasks.CapyTask) (siteKeyTask, error) {
	taskType, err := variant("CapyTask").For(task.Proxy)
	if err != nil {
		return siteKeyTask{}, err
	}

	return siteKeyTask{
		Type:        taskType,
		WebsiteURL:  task.WebsiteURL,
		WebsiteKey:  task.WebsiteKey,
		UserAgent:   task.UserAgent,
		ProxyFields: solverapi.NewProxyFields(task.Proxy),
	}, nil
}

func mapFriendlyCaptcha(task *tasks.FriendlyCaptchaTask) (friendlyCaptchaTask, error) {
	taskType, err := variant("FriendlyCaptchaTask").For(task.Proxy)
	if err != nil {
		return friendlyCaptchaTask{}, err
	}

	return friendlyCaptchaTask{
		Type:           taskType,
		WebsiteURL:     task.WebsiteURL,
		WebsiteKey:     task.WebsiteKey,
		Version:        task.Version,
		ModuleScript:   task.ModuleScript,
		NomoduleScript: task.NomoduleScript,
		ProxyFields:    solverapi.NewProxyFields(task.Proxy),
	}, nil
}

func mapLemin(task *tasks.LeminTask) (leminTask, error) {
	taskType, err := variant("LeminTask").For(task.Proxy)
	if err != nil {
		return leminTask{}, err
	}

	return leminTask{
		Type:                    taskType,
		WebsiteURL:              task.WebsiteURL,
		CaptchaID:               task.CaptchaID,
		DivID:                   task.DivID,
		LeminAPIServerSubdomain: task.APIServerSubdomain,
		UserAgent:               task.UserAgent,
		ProxyFields:             solverapi.NewProxyFields(task.Proxy),
	}, nil
}

func mapCutCaptcha(task *tasks.CutCaptchaTask) (cutCaptchaTask, error) {
	taskType, err := variant("CutCaptchaTask").For(task.Proxy)
	if err != nil {
		return cutCaptchaTask{}, err
	}

	return cutCaptchaTask{
		Type:        taskType,
		WebsiteURL:  task.WebsiteURL,
		MiseryKey:   task.MiseryKey,
		APIKey:      task.APIKey,
		ProxyFields: solverapi.NewProxyFields(task.Proxy),
	}, nil
}

func mapAltcha(task *tasks.AltchaTask) (altchaTask, error) {
	taskType, err := variant("AltchaTask").For(task.Proxy)
	if err != nil {
		return altchaTask{}, err
	}

	return altchaTask{
		Type:          taskType,
		WebsiteURL:    task.WebsiteURL,
		ChallengeURL:  task.ChallengeURL,
		ChallengeJSON: task.ChallengeJSON,
		ProxyFields:   solverapi.NewProxyFields(task.Proxy),
	}, nil
}

func mapKeyCaptcha(task *tasks.KeyCaptchaTask) (keyCaptchaTask, error) {
	taskType, err := variant("KeyCaptchaTask").For(task.Proxy)
	if err != nil {
		return keyCaptchaTask{}, err
	}

	return keyCaptchaTask{
		Type:           taskType,
		WebsiteURL:     task.WebsiteURL,
		UserID:         task.UserID,
		SessionID:      task.SessionID,
		WebServerSign:  task.WebServerSign,
		WebServerSign2: task.WebServerSign2,
		ProxyFields:    solverapi.NewProxyFields(task.Proxy),
	}, nil
}

func mapCyberSiAra(task *tasks.CyberSiAraTask) (cyberSiAraTask, error) {
	taskType, err := variant("AntiCyberSiAraTask").For(task.Proxy)
	if err != nil {
		return cyberSiAraTask{}, err
	}

	return cyberSiAraTask{
		Type:             taskType,
		WebsiteURL:       task.WebsiteURL,
		SlideMasterURLID: task.SlideMasterURLID,
		UserAgent:        task.UserAgent,
		ProxyFields:      solverapi.NewProxyFields(task.Proxy),
	}, nil
}

func mapAtbCaptcha(task *tasks.AtbCaptchaTask) (atbCaptchaTask, error) {
	taskType, err := variant("AtbCaptchaTask").For(task.Proxy)
	if err != nil {
		return atbCaptchaTask{}, err
	}

	return atbCaptchaTask{
		Type:        taskType,
		WebsiteURL:  task.WebsiteURL,
		AppID:       task.AppID,
		APIServer:   task.APIServer,
		ProxyFields: solverapi.NewProxyFields(task.Proxy),
	}, nil
}

func mapTencent(task *tasks.TencentTask) (tencentTask, error) {
	taskType, err := variant("TencentTask").For(task.Proxy)
	if err != nil {
		return tencentTask{}, err
	}

	return tencentTask{
		Type:          taskType,
		WebsiteURL:    task.WebsiteURL,
		AppID:         task.AppID,
		CaptchaScript: task.CaptchaScript,
		ProxyFields:   solverapi.NewProxyFields(task.Proxy),
	}, nil
}

func mapCaptchaFox(task *tasks.CaptchaFoxTask) (captchaFoxTask, error) {
	taskType, err := proxyOnly("CaptchaFoxTask").For(task.Proxy)
	if err != nil {
		return captchaFoxTask{}, err
	}

	return captchaFoxTask{
		Type:        taskType,
		WebsiteURL:  task.WebsiteURL,
		WebsiteKey:  task.WebsiteKey,
		UserAgent:   task.UserAgent,
		APIServer:   task.APIServer,
		ProxyFields: solverapi.NewProxyFields(task.Proxy),
	}, nil
}

func mapYandexSmartCaptcha(task *tasks.YandexSmartCaptchaTask) (yandexSmartCaptchaTask, error) {
	taskType, err := variant("YandexSmartCaptchaTask").For(task.Proxy)
	if err != nil {
		return yandexSmartCaptchaTask{}, err
	}

	return yandexSmartCaptchaTask{
		Type:        taskType,
		WebsiteURL:  task.WebsiteURL,
		WebsiteKey:  task.WebsiteKey,
		UserAgent:   task.UserAgent,
		Cookies:     task.Cookies,
		ProxyFields: solverapi.NewProxyFields(task.Proxy),
	}, nil
}
