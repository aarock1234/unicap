package anticaptcha

import (
	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

type funCaptchaTask struct {
	Type                     string `json:"type"`
	WebsiteURL               string `json:"websiteURL"`
	WebsitePublicKey         string `json:"websitePublicKey"`
	FuncaptchaAPIJSSubdomain string `json:"funcaptchaApiJSSubdomain,omitempty"`
	Data                     string `json:"data,omitempty"`
	UserAgent                string `json:"userAgent,omitempty"`
	solverapi.ProxyFields
}

type turnstileTask struct {
	Type        string `json:"type"`
	WebsiteURL  string `json:"websiteURL"`
	WebsiteKey  string `json:"websiteKey"`
	Action      string `json:"action,omitempty"`
	CData       string `json:"cData,omitempty"`
	ChlPageData string `json:"chlPageData,omitempty"`
	solverapi.ProxyFields
}

type geeTestTask struct {
	Type                      string         `json:"type"`
	WebsiteURL                string         `json:"websiteURL"`
	GT                        string         `json:"gt"`
	Challenge                 string         `json:"challenge,omitempty"`
	Version                   int            `json:"version"`
	InitParameters            *geeTestParams `json:"initParameters,omitempty"`
	GeetestAPIServerSubdomain string         `json:"geetestApiServerSubdomain,omitempty"`
	UserAgent                 string         `json:"userAgent,omitempty"`
	solverapi.ProxyFields
}

type geeTestParams struct {
	CaptchaID string `json:"captcha_id"`
	RiskType  string `json:"risk_type,omitempty"`
}

type amazonTask struct {
	Type            string `json:"type"`
	WebsiteURL      string `json:"websiteURL"`
	WebsiteKey      string `json:"websiteKey"`
	IV              string `json:"iv"`
	Context         string `json:"context"`
	ChallengeScript string `json:"challengeScript,omitempty"`
	CaptchaScript   string `json:"captchaScript,omitempty"`
	JSAPIScript     string `json:"jsapiScript,omitempty"`
	solverapi.ProxyFields
}

func mapFunCaptcha(task *tasks.FunCaptchaTask) (funCaptchaTask, error) {
	if err := requireUserAgentWithProxy(task.Proxy, task.UserAgent); err != nil {
		return funCaptchaTask{}, err
	}

	taskType, err := variant("FunCaptchaTask").For(task.Proxy)
	if err != nil {
		return funCaptchaTask{}, err
	}

	proxy, err := proxyFields(task.Proxy)
	if err != nil {
		return funCaptchaTask{}, err
	}

	return funCaptchaTask{
		Type:                     taskType,
		WebsiteURL:               task.WebsiteURL,
		WebsitePublicKey:         task.WebsitePublicKey,
		FuncaptchaAPIJSSubdomain: task.APIJSSubdomain,
		Data:                     task.Data,
		UserAgent:                task.UserAgent,
		ProxyFields:              proxy,
	}, nil
}

func mapTurnstile(task *tasks.TurnstileTask) (turnstileTask, error) {
	taskType, err := variant("TurnstileTask").For(task.Proxy)
	if err != nil {
		return turnstileTask{}, err
	}

	proxy, err := proxyFields(task.Proxy)
	if err != nil {
		return turnstileTask{}, err
	}

	return turnstileTask{
		Type:        taskType,
		WebsiteURL:  task.WebsiteURL,
		WebsiteKey:  task.WebsiteKey,
		Action:      task.Action,
		CData:       task.CData,
		ChlPageData: task.PageData,
		ProxyFields: proxy,
	}, nil
}

func mapGeeTest(task *tasks.GeeTestTask) (geeTestTask, error) {
	if err := requireUserAgentWithProxy(task.Proxy, task.UserAgent); err != nil {
		return geeTestTask{}, err
	}

	taskType, err := variant("GeeTestTask").For(task.Proxy)
	if err != nil {
		return geeTestTask{}, err
	}

	proxy, err := proxyFields(task.Proxy)
	if err != nil {
		return geeTestTask{}, err
	}

	return geeTestTask{
		Type:                      taskType,
		WebsiteURL:                task.WebsiteURL,
		GT:                        task.GT,
		Challenge:                 task.Challenge,
		Version:                   3,
		GeetestAPIServerSubdomain: task.APIServerSubdomain,
		UserAgent:                 task.UserAgent,
		ProxyFields:               proxy,
	}, nil
}

func mapGeeTestV4(task *tasks.GeeTestV4Task) (geeTestTask, error) {
	if err := requireUserAgentWithProxy(task.Proxy, task.UserAgent); err != nil {
		return geeTestTask{}, err
	}

	taskType, err := variant("GeeTestTask").For(task.Proxy)
	if err != nil {
		return geeTestTask{}, err
	}

	proxy, err := proxyFields(task.Proxy)
	if err != nil {
		return geeTestTask{}, err
	}

	return geeTestTask{
		Type:       taskType,
		WebsiteURL: task.WebsiteURL,
		GT:         task.CaptchaID,
		Version:    4,
		InitParameters: &geeTestParams{
			CaptchaID: task.CaptchaID,
			RiskType:  task.RiskType,
		},
		GeetestAPIServerSubdomain: task.APIServerSubdomain,
		UserAgent:                 task.UserAgent,
		ProxyFields:               proxy,
	}, nil
}

func mapAWSWAF(task *tasks.AWSWAFTask) (amazonTask, error) {
	if err := solverapi.Check(
		solverapi.Required(
			solverapi.Setting("Key", task.Key != ""),
			solverapi.Setting("IV", task.IV != ""),
			solverapi.Setting("Context", task.Context != ""),
		),
		solverapi.Unsupported(
			solverapi.Setting("ProblemURL", task.ProblemURL != ""),
			solverapi.Setting("APIKey", task.APIKey != ""),
			solverapi.Setting("ExistingToken", task.ExistingToken != ""),
		),
	); err != nil {
		return amazonTask{}, err
	}

	taskType, err := variant("AmazonTask").For(task.Proxy)
	if err != nil {
		return amazonTask{}, err
	}

	proxy, err := proxyFields(task.Proxy)
	if err != nil {
		return amazonTask{}, err
	}

	return amazonTask{
		Type:            taskType,
		WebsiteURL:      task.WebsiteURL,
		WebsiteKey:      task.Key,
		IV:              task.IV,
		Context:         task.Context,
		ChallengeScript: task.ChallengeScript,
		CaptchaScript:   task.CaptchaScript,
		JSAPIScript:     task.JSAPIScript,
		ProxyFields:     proxy,
	}, nil
}
