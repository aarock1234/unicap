package anticaptcha

import (
	"fmt"

	"upicap/pkg/upicap"
	"upicap/pkg/upicap/tasks"
)

// Task types for AntiCaptcha API
type anticaptchaReCaptchaV2Task struct {
	Type                string `json:"type"`
	WebsiteURL          string `json:"websiteURL"`
	WebsiteKey          string `json:"websiteKey"`
	IsInvisible         bool   `json:"isInvisible,omitempty"`
	RecaptchaDataSValue string `json:"recaptchaDataSValue,omitempty"`
}

type anticaptchaReCaptchaV3Task struct {
	Type       string  `json:"type"`
	WebsiteURL string  `json:"websiteURL"`
	WebsiteKey string  `json:"websiteKey"`
	PageAction string  `json:"pageAction,omitempty"`
	MinScore   float64 `json:"minScore,omitempty"`
}

type anticaptchaReCaptchaV2EnterpriseTask struct {
	Type              string         `json:"type"`
	WebsiteURL        string         `json:"websiteURL"`
	WebsiteKey        string         `json:"websiteKey"`
	IsInvisible       bool           `json:"isInvisible,omitempty"`
	EnterprisePayload map[string]any `json:"enterprisePayload,omitempty"`
	ApiDomain         string         `json:"apiDomain,omitempty"`
}

type anticaptchaReCaptchaV3EnterpriseTask struct {
	Type              string         `json:"type"`
	WebsiteURL        string         `json:"websiteURL"`
	WebsiteKey        string         `json:"websiteKey"`
	PageAction        string         `json:"pageAction,omitempty"`
	MinScore          float64        `json:"minScore,omitempty"`
	EnterprisePayload map[string]any `json:"enterprisePayload,omitempty"`
	ApiDomain         string         `json:"apiDomain,omitempty"`
}

type anticaptchaFunCaptchaTask struct {
	Type                     string `json:"type"`
	WebsiteURL               string `json:"websiteURL"`
	WebsitePublicKey         string `json:"websitePublicKey"`
	FuncaptchaApiJSSubdomain string `json:"funcaptchaApiJSSubdomain,omitempty"`
	Data                     string `json:"data,omitempty"`
	ProxyType                string `json:"proxyType,omitempty"`
	ProxyAddress             string `json:"proxyAddress,omitempty"`
	ProxyPort                int    `json:"proxyPort,omitempty"`
	ProxyLogin               string `json:"proxyLogin,omitempty"`
	ProxyPassword            string `json:"proxyPassword,omitempty"`
}

type anticaptchaTurnstileTask struct {
	Type          string `json:"type"`
	WebsiteURL    string `json:"websiteURL"`
	WebsiteKey    string `json:"websiteKey"`
	Action        string `json:"action,omitempty"`
	CData         string `json:"cData,omitempty"`
	ChlPageData   string `json:"chlPageData,omitempty"`
	ProxyType     string `json:"proxyType,omitempty"`
	ProxyAddress  string `json:"proxyAddress,omitempty"`
	ProxyPort     int    `json:"proxyPort,omitempty"`
	ProxyLogin    string `json:"proxyLogin,omitempty"`
	ProxyPassword string `json:"proxyPassword,omitempty"`
}

type anticaptchaGeeTestTask struct {
	Type                      string      `json:"type"`
	WebsiteURL                string      `json:"websiteURL"`
	GT                        string      `json:"gt,omitempty"`
	Challenge                 string      `json:"challenge,omitempty"`
	Version                   int         `json:"version,omitempty"`
	InitParameters            geetestInit `json:"initParameters,omitempty"`
	GeetestApiServerSubdomain string      `json:"geetestApiServerSubdomain,omitempty"`
	ProxyType                 string      `json:"proxyType,omitempty"`
	ProxyAddress              string      `json:"proxyAddress,omitempty"`
	ProxyPort                 int         `json:"proxyPort,omitempty"`
	ProxyLogin                string      `json:"proxyLogin,omitempty"`
	ProxyPassword             string      `json:"proxyPassword,omitempty"`
}

type geetestInit struct {
	CaptchaID string `json:"captcha_id,omitempty"`
}

// mapToAntiCaptchaTask converts a universal task to AntiCaptcha format
func mapToAntiCaptchaTask(task upicap.Task) (any, error) {
	switch t := task.(type) {
	case *tasks.ReCaptchaV2Task:
		return mapReCaptchaV2(t), nil
	case *tasks.ReCaptchaV3Task:
		return mapReCaptchaV3(t), nil
	case *tasks.ReCaptchaV2EnterpriseTask:
		return mapReCaptchaV2Enterprise(t), nil
	case *tasks.ReCaptchaV3EnterpriseTask:
		return mapReCaptchaV3Enterprise(t), nil
	case *tasks.FunCaptchaTask:
		return mapFunCaptcha(t), nil
	case *tasks.TurnstileTask:
		return mapTurnstile(t), nil
	case *tasks.CloudflareChallengeTask:
		return nil, fmt.Errorf("cloudflare challenge not supported by anticaptcha")
	case *tasks.DataDomeTask:
		return nil, fmt.Errorf("datadome not supported by anticaptcha")
	case *tasks.GeeTestTask:
		return mapGeeTest(t), nil
	case *tasks.GeeTestV4Task:
		return mapGeeTestV4(t), nil
	default:
		return nil, fmt.Errorf("unsupported task type: %s", task.Type())
	}
}

func mapReCaptchaV2(task *tasks.ReCaptchaV2Task) anticaptchaReCaptchaV2Task {
	return anticaptchaReCaptchaV2Task{
		Type:                "RecaptchaV2TaskProxyless",
		WebsiteURL:          task.WebsiteURL,
		WebsiteKey:          task.WebsiteKey,
		IsInvisible:         task.IsInvisible,
		RecaptchaDataSValue: task.DataS,
	}
}

func mapReCaptchaV3(task *tasks.ReCaptchaV3Task) anticaptchaReCaptchaV3Task {
	return anticaptchaReCaptchaV3Task{
		Type:       "RecaptchaV3TaskProxyless",
		WebsiteURL: task.WebsiteURL,
		WebsiteKey: task.WebsiteKey,
		PageAction: task.PageAction,
		MinScore:   task.MinScore,
	}
}

func mapReCaptchaV2Enterprise(task *tasks.ReCaptchaV2EnterpriseTask) anticaptchaReCaptchaV2EnterpriseTask {
	return anticaptchaReCaptchaV2EnterpriseTask{
		Type:              "RecaptchaV2EnterpriseTaskProxyless",
		WebsiteURL:        task.WebsiteURL,
		WebsiteKey:        task.WebsiteKey,
		IsInvisible:       task.IsInvisible,
		EnterprisePayload: task.EnterpriseData,
		ApiDomain:         task.ApiDomain,
	}
}

func mapReCaptchaV3Enterprise(task *tasks.ReCaptchaV3EnterpriseTask) anticaptchaReCaptchaV3EnterpriseTask {
	return anticaptchaReCaptchaV3EnterpriseTask{
		Type:              "RecaptchaV3EnterpriseTaskProxyless",
		WebsiteURL:        task.WebsiteURL,
		WebsiteKey:        task.WebsiteKey,
		PageAction:        task.PageAction,
		MinScore:          task.MinScore,
		EnterprisePayload: task.EnterpriseData,
		ApiDomain:         task.ApiDomain,
	}
}

func mapFunCaptcha(task *tasks.FunCaptchaTask) anticaptchaFunCaptchaTask {
	result := anticaptchaFunCaptchaTask{
		Type:                     "FunCaptchaTaskProxyless",
		WebsiteURL:               task.WebsiteURL,
		WebsitePublicKey:         task.WebsitePublicKey,
		FuncaptchaApiJSSubdomain: task.APIJSSubdomain,
		Data:                     task.Data,
	}

	if task.Proxy.IsSet() {
		result.Type = "FunCaptchaTask"
		result.ProxyType = string(task.Proxy.Type)
		result.ProxyAddress = task.Proxy.Address
		result.ProxyPort = task.Proxy.Port
		result.ProxyLogin = task.Proxy.Login
		result.ProxyPassword = task.Proxy.Password
	}

	return result
}

func mapTurnstile(task *tasks.TurnstileTask) anticaptchaTurnstileTask {
	result := anticaptchaTurnstileTask{
		Type:        "TurnstileTaskProxyless",
		WebsiteURL:  task.WebsiteURL,
		WebsiteKey:  task.WebsiteKey,
		Action:      task.Action,
		CData:       task.CData,
		ChlPageData: task.PageData,
	}

	if task.Proxy.IsSet() {
		result.Type = "TurnstileTask"
		result.ProxyType = string(task.Proxy.Type)
		result.ProxyAddress = task.Proxy.Address
		result.ProxyPort = task.Proxy.Port
		result.ProxyLogin = task.Proxy.Login
		result.ProxyPassword = task.Proxy.Password
	}

	return result
}

func mapGeeTest(task *tasks.GeeTestTask) anticaptchaGeeTestTask {
	result := anticaptchaGeeTestTask{
		Type:                      "GeeTestTaskProxyless",
		WebsiteURL:                task.WebsiteURL,
		GT:                        task.GT,
		Challenge:                 task.Challenge,
		Version:                   3,
		GeetestApiServerSubdomain: task.APIServerSubdomain,
	}

	if task.Proxy.IsSet() {
		result.Type = "GeeTestTask"
		result.ProxyType = string(task.Proxy.Type)
		result.ProxyAddress = task.Proxy.Address
		result.ProxyPort = task.Proxy.Port
		result.ProxyLogin = task.Proxy.Login
		result.ProxyPassword = task.Proxy.Password
	}

	return result
}

func mapGeeTestV4(task *tasks.GeeTestV4Task) anticaptchaGeeTestTask {
	result := anticaptchaGeeTestTask{
		Type:                      "GeeTestTaskProxyless",
		WebsiteURL:                task.WebsiteURL,
		Version:                   4,
		InitParameters:            geetestInit{CaptchaID: task.CaptchaID},
		GeetestApiServerSubdomain: task.APIServerSubdomain,
	}

	if task.Proxy.IsSet() {
		result.Type = "GeeTestTask"
		result.ProxyType = string(task.Proxy.Type)
		result.ProxyAddress = task.Proxy.Address
		result.ProxyPort = task.Proxy.Port
		result.ProxyLogin = task.Proxy.Login
		result.ProxyPassword = task.Proxy.Password
	}

	return result
}
