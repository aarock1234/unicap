package capsolver

import (
	"fmt"

	"upicap/pkg/upicap"
	"upicap/pkg/upicap/tasks"
)

// Task types for CapSolver API
type capsolverReCaptchaV2Task struct {
	Type                string `json:"type"`
	WebsiteURL          string `json:"websiteURL"`
	WebsiteKey          string `json:"websiteKey"`
	IsInvisible         bool   `json:"isInvisible,omitempty"`
	RecaptchaDataSValue string `json:"recaptchaDataSValue,omitempty"`
	PageAction          string `json:"pageAction,omitempty"`
	ProxyType           string `json:"proxyType,omitempty"`
	ProxyAddress        string `json:"proxyAddress,omitempty"`
	ProxyPort           int    `json:"proxyPort,omitempty"`
	ProxyLogin          string `json:"proxyLogin,omitempty"`
	ProxyPassword       string `json:"proxyPassword,omitempty"`
}

type capsolverReCaptchaV3Task struct {
	Type          string  `json:"type"`
	WebsiteURL    string  `json:"websiteURL"`
	WebsiteKey    string  `json:"websiteKey"`
	PageAction    string  `json:"pageAction,omitempty"`
	MinScore      float64 `json:"minScore,omitempty"`
	ProxyType     string  `json:"proxyType,omitempty"`
	ProxyAddress  string  `json:"proxyAddress,omitempty"`
	ProxyPort     int     `json:"proxyPort,omitempty"`
	ProxyLogin    string  `json:"proxyLogin,omitempty"`
	ProxyPassword string  `json:"proxyPassword,omitempty"`
}

type capsolverReCaptchaV2EnterpriseTask struct {
	Type              string         `json:"type"`
	WebsiteURL        string         `json:"websiteURL"`
	WebsiteKey        string         `json:"websiteKey"`
	IsInvisible       bool           `json:"isInvisible,omitempty"`
	EnterprisePayload map[string]any `json:"enterprisePayload,omitempty"`
	ApiDomain         string         `json:"apiDomain,omitempty"`
	ProxyType         string         `json:"proxyType,omitempty"`
	ProxyAddress      string         `json:"proxyAddress,omitempty"`
	ProxyPort         int            `json:"proxyPort,omitempty"`
	ProxyLogin        string         `json:"proxyLogin,omitempty"`
	ProxyPassword     string         `json:"proxyPassword,omitempty"`
}

type capsolverReCaptchaV3EnterpriseTask struct {
	Type              string         `json:"type"`
	WebsiteURL        string         `json:"websiteURL"`
	WebsiteKey        string         `json:"websiteKey"`
	PageAction        string         `json:"pageAction,omitempty"`
	MinScore          float64        `json:"minScore,omitempty"`
	EnterprisePayload map[string]any `json:"enterprisePayload,omitempty"`
	ApiDomain         string         `json:"apiDomain,omitempty"`
	ProxyType         string         `json:"proxyType,omitempty"`
	ProxyAddress      string         `json:"proxyAddress,omitempty"`
	ProxyPort         int            `json:"proxyPort,omitempty"`
	ProxyLogin        string         `json:"proxyLogin,omitempty"`
	ProxyPassword     string         `json:"proxyPassword,omitempty"`
}

type capsolverFunCaptchaTask struct {
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

type capsolverTurnstileTask struct {
	Type          string            `json:"type"`
	WebsiteURL    string            `json:"websiteURL"`
	WebsiteKey    string            `json:"websiteKey"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	ProxyType     string            `json:"proxyType,omitempty"`
	ProxyAddress  string            `json:"proxyAddress,omitempty"`
	ProxyPort     int               `json:"proxyPort,omitempty"`
	ProxyLogin    string            `json:"proxyLogin,omitempty"`
	ProxyPassword string            `json:"proxyPassword,omitempty"`
}

type capsolverGeeTestTask struct {
	Type                      string `json:"type"`
	WebsiteURL                string `json:"websiteURL"`
	GT                        string `json:"gt"`
	Challenge                 string `json:"challenge"`
	GeetestApiServerSubdomain string `json:"geetestApiServerSubdomain,omitempty"`
	ProxyType                 string `json:"proxyType,omitempty"`
	ProxyAddress              string `json:"proxyAddress,omitempty"`
	ProxyPort                 int    `json:"proxyPort,omitempty"`
	ProxyLogin                string `json:"proxyLogin,omitempty"`
	ProxyPassword             string `json:"proxyPassword,omitempty"`
}

type capsolverGeeTestV4Task struct {
	Type                      string `json:"type"`
	WebsiteURL                string `json:"websiteURL"`
	CaptchaID                 string `json:"captchaId"`
	GeetestApiServerSubdomain string `json:"geetestApiServerSubdomain,omitempty"`
	ProxyType                 string `json:"proxyType,omitempty"`
	ProxyAddress              string `json:"proxyAddress,omitempty"`
	ProxyPort                 int    `json:"proxyPort,omitempty"`
	ProxyLogin                string `json:"proxyLogin,omitempty"`
	ProxyPassword             string `json:"proxyPassword,omitempty"`
}

type capsolverCloudflareChallengeTask struct {
	Type          string `json:"type"`
	WebsiteURL    string `json:"websiteURL"`
	HTML          string `json:"html,omitempty"`
	UserAgent     string `json:"userAgent,omitempty"`
	Proxy         string `json:"proxy"`
	ProxyType     string `json:"proxyType"`
	ProxyAddress  string `json:"proxyAddress"`
	ProxyPort     int    `json:"proxyPort"`
	ProxyLogin    string `json:"proxyLogin,omitempty"`
	ProxyPassword string `json:"proxyPassword,omitempty"`
}

type capsolverDataDomeTask struct {
	Type          string `json:"type"`
	CaptchaURL    string `json:"captchaUrl"`
	UserAgent     string `json:"userAgent"`
	Proxy         string `json:"proxy"`
	ProxyType     string `json:"proxyType"`
	ProxyAddress  string `json:"proxyAddress"`
	ProxyPort     int    `json:"proxyPort"`
	ProxyLogin    string `json:"proxyLogin,omitempty"`
	ProxyPassword string `json:"proxyPassword,omitempty"`
}

// mapToCapSolverTask converts a universal task to CapSolver format
func mapToCapSolverTask(task upicap.Task) (any, error) {
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
		return mapCloudflareChallenge(t), nil
	case *tasks.DataDomeTask:
		return mapDataDome(t), nil
	case *tasks.GeeTestTask:
		return mapGeeTest(t), nil
	case *tasks.GeeTestV4Task:
		return mapGeeTestV4(t), nil
	default:
		return nil, fmt.Errorf("unsupported task type: %s", task.Type())
	}
}

func mapReCaptchaV2(task *tasks.ReCaptchaV2Task) capsolverReCaptchaV2Task {
	taskType := "ReCaptchaV2TaskProxyLess"
	result := capsolverReCaptchaV2Task{
		Type:                taskType,
		WebsiteURL:          task.WebsiteURL,
		WebsiteKey:          task.WebsiteKey,
		IsInvisible:         task.IsInvisible,
		RecaptchaDataSValue: task.DataS,
		PageAction:          task.PageAction,
	}

	if task.Proxy.IsSet() {
		result.Type = "ReCaptchaV2Task"
		result.ProxyType = string(task.Proxy.Type)
		result.ProxyAddress = task.Proxy.Address
		result.ProxyPort = task.Proxy.Port
		result.ProxyLogin = task.Proxy.Login
		result.ProxyPassword = task.Proxy.Password
	}

	return result
}

func mapReCaptchaV3(task *tasks.ReCaptchaV3Task) capsolverReCaptchaV3Task {
	result := capsolverReCaptchaV3Task{
		Type:       "ReCaptchaV3TaskProxyLess",
		WebsiteURL: task.WebsiteURL,
		WebsiteKey: task.WebsiteKey,
		PageAction: task.PageAction,
		MinScore:   task.MinScore,
	}

	if task.Proxy.IsSet() {
		result.Type = "ReCaptchaV3Task"
		result.ProxyType = string(task.Proxy.Type)
		result.ProxyAddress = task.Proxy.Address
		result.ProxyPort = task.Proxy.Port
		result.ProxyLogin = task.Proxy.Login
		result.ProxyPassword = task.Proxy.Password
	}

	return result
}

func mapReCaptchaV2Enterprise(task *tasks.ReCaptchaV2EnterpriseTask) capsolverReCaptchaV2EnterpriseTask {
	result := capsolverReCaptchaV2EnterpriseTask{
		Type:              "ReCaptchaV2EnterpriseTaskProxyLess",
		WebsiteURL:        task.WebsiteURL,
		WebsiteKey:        task.WebsiteKey,
		IsInvisible:       task.IsInvisible,
		EnterprisePayload: task.EnterpriseData,
		ApiDomain:         task.ApiDomain,
	}

	if task.Proxy.IsSet() {
		result.Type = "ReCaptchaV2EnterpriseTask"
		result.ProxyType = string(task.Proxy.Type)
		result.ProxyAddress = task.Proxy.Address
		result.ProxyPort = task.Proxy.Port
		result.ProxyLogin = task.Proxy.Login
		result.ProxyPassword = task.Proxy.Password
	}

	return result
}

func mapReCaptchaV3Enterprise(task *tasks.ReCaptchaV3EnterpriseTask) capsolverReCaptchaV3EnterpriseTask {
	result := capsolverReCaptchaV3EnterpriseTask{
		Type:              "ReCaptchaV3EnterpriseTaskProxyLess",
		WebsiteURL:        task.WebsiteURL,
		WebsiteKey:        task.WebsiteKey,
		PageAction:        task.PageAction,
		MinScore:          task.MinScore,
		EnterprisePayload: task.EnterpriseData,
		ApiDomain:         task.ApiDomain,
	}

	if task.Proxy.IsSet() {
		result.Type = "ReCaptchaV3EnterpriseTask"
		result.ProxyType = string(task.Proxy.Type)
		result.ProxyAddress = task.Proxy.Address
		result.ProxyPort = task.Proxy.Port
		result.ProxyLogin = task.Proxy.Login
		result.ProxyPassword = task.Proxy.Password
	}

	return result
}

func mapFunCaptcha(task *tasks.FunCaptchaTask) capsolverFunCaptchaTask {
	result := capsolverFunCaptchaTask{
		Type:                     "FunCaptchaTaskProxyLess",
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

func mapTurnstile(task *tasks.TurnstileTask) capsolverTurnstileTask {
	result := capsolverTurnstileTask{
		Type:       "AntiTurnstileTaskProxyLess",
		WebsiteURL: task.WebsiteURL,
		WebsiteKey: task.WebsiteKey,
	}

	if task.Action != "" || task.CData != "" {
		result.Metadata = make(map[string]string)
		if task.Action != "" {
			result.Metadata["action"] = task.Action
		}
		if task.CData != "" {
			result.Metadata["cdata"] = task.CData
		}
	}

	if task.Proxy.IsSet() {
		result.Type = "AntiTurnstileTask"
		result.ProxyType = string(task.Proxy.Type)
		result.ProxyAddress = task.Proxy.Address
		result.ProxyPort = task.Proxy.Port
		result.ProxyLogin = task.Proxy.Login
		result.ProxyPassword = task.Proxy.Password
	}

	return result
}

func mapCloudflareChallenge(task *tasks.CloudflareChallengeTask) capsolverCloudflareChallengeTask {
	proxyString := fmt.Sprintf("%s:%s:%d:%s:%s",
		task.Proxy.Type,
		task.Proxy.Address,
		task.Proxy.Port,
		task.Proxy.Login,
		task.Proxy.Password,
	)

	return capsolverCloudflareChallengeTask{
		Type:          "AntiCloudflareTask",
		WebsiteURL:    task.WebsiteURL,
		HTML:          task.HTML,
		UserAgent:     task.UserAgent,
		Proxy:         proxyString,
		ProxyType:     string(task.Proxy.Type),
		ProxyAddress:  task.Proxy.Address,
		ProxyPort:     task.Proxy.Port,
		ProxyLogin:    task.Proxy.Login,
		ProxyPassword: task.Proxy.Password,
	}
}

func mapDataDome(task *tasks.DataDomeTask) capsolverDataDomeTask {
	proxyString := fmt.Sprintf("%s:%s:%d:%s:%s",
		task.Proxy.Type,
		task.Proxy.Address,
		task.Proxy.Port,
		task.Proxy.Login,
		task.Proxy.Password,
	)

	return capsolverDataDomeTask{
		Type:          "DatadomeSliderTask",
		CaptchaURL:    task.CaptchaURL,
		UserAgent:     task.UserAgent,
		Proxy:         proxyString,
		ProxyType:     string(task.Proxy.Type),
		ProxyAddress:  task.Proxy.Address,
		ProxyPort:     task.Proxy.Port,
		ProxyLogin:    task.Proxy.Login,
		ProxyPassword: task.Proxy.Password,
	}
}

func mapGeeTest(task *tasks.GeeTestTask) capsolverGeeTestTask {
	result := capsolverGeeTestTask{
		Type:                      "GeeTestTaskProxyLess",
		WebsiteURL:                task.WebsiteURL,
		GT:                        task.GT,
		Challenge:                 task.Challenge,
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

func mapGeeTestV4(task *tasks.GeeTestV4Task) capsolverGeeTestV4Task {
	result := capsolverGeeTestV4Task{
		Type:                      "GeeTestTaskProxyLess",
		WebsiteURL:                task.WebsiteURL,
		CaptchaID:                 task.CaptchaID,
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
