package capsolver

import (
	"fmt"

	"github.com/aarock1234/unicap/unicap"
	"github.com/aarock1234/unicap/unicap/tasks"
)

// Task types for CapSolver API

// Base proxy fields for all tasks
type capsolverProxyFields struct {
	ProxyType     string `json:"proxyType,omitempty"`
	ProxyAddress  string `json:"proxyAddress,omitempty"`
	ProxyPort     int    `json:"proxyPort,omitempty"`
	ProxyLogin    string `json:"proxyLogin,omitempty"`
	ProxyPassword string `json:"proxyPassword,omitempty"`
}

// Base ReCaptcha fields
type capsolverReCaptchaV2Base struct {
	Type                string `json:"type"`
	WebsiteURL          string `json:"websiteURL"`
	WebsiteKey          string `json:"websiteKey"`
	IsInvisible         bool   `json:"isInvisible,omitempty"`
	PageAction          string `json:"pageAction,omitempty"`
	RecaptchaDataSValue string `json:"recaptchaDataSValue,omitempty"`
	IsSession           bool   `json:"isSession,omitempty"`
	ApiDomain           string `json:"apiDomain,omitempty"`
}

type capsolverReCaptchaV2Task struct {
	capsolverReCaptchaV2Base
	EnterprisePayload map[string]any `json:"enterprisePayload,omitempty"`
	capsolverProxyFields
}

type capsolverReCaptchaV2EnterpriseTask struct {
	capsolverReCaptchaV2Base
	EnterprisePayload map[string]any `json:"enterprisePayload,omitempty"`
	capsolverProxyFields
}

type capsolverReCaptchaV3Base struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	WebsiteKey string `json:"websiteKey"`
	PageAction string `json:"pageAction,omitempty"`
	IsSession  bool   `json:"isSession,omitempty"`
	ApiDomain  string `json:"apiDomain,omitempty"`
}

type capsolverReCaptchaV3Task struct {
	capsolverReCaptchaV3Base
	EnterprisePayload map[string]any `json:"enterprisePayload,omitempty"`
	capsolverProxyFields
}

type capsolverReCaptchaV3EnterpriseTask struct {
	capsolverReCaptchaV3Base
	EnterprisePayload map[string]any `json:"enterprisePayload,omitempty"`
	capsolverProxyFields
}

type capsolverFunCaptchaTask struct {
	Type                     string `json:"type"`
	WebsiteURL               string `json:"websiteURL"`
	WebsitePublicKey         string `json:"websitePublicKey"`
	FuncaptchaApiJSSubdomain string `json:"funcaptchaApiJSSubdomain,omitempty"`
	Data                     string `json:"data,omitempty"`
	capsolverProxyFields
}

type capsolverTurnstileTask struct {
	Type       string            `json:"type"`
	WebsiteURL string            `json:"websiteURL"`
	WebsiteKey string            `json:"websiteKey"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	capsolverProxyFields
}

type capsolverGeeTestTask struct {
	Type                      string `json:"type"`
	WebsiteURL                string `json:"websiteURL"`
	GT                        string `json:"gt"`
	Challenge                 string `json:"challenge"`
	GeetestApiServerSubdomain string `json:"geetestApiServerSubdomain,omitempty"`
	capsolverProxyFields
}

type capsolverGeeTestV4Task struct {
	Type                      string `json:"type"`
	WebsiteURL                string `json:"websiteURL"`
	CaptchaID                 string `json:"captchaId"`
	GeetestApiServerSubdomain string `json:"geetestApiServerSubdomain,omitempty"`
	capsolverProxyFields
}

type capsolverHCaptchaTask struct {
	Type              string         `json:"type"`
	WebsiteURL        string         `json:"websiteURL"`
	WebsiteKey        string         `json:"websiteKey"`
	IsInvisible       bool           `json:"isInvisible,omitempty"`
	EnterprisePayload map[string]any `json:"enterprisePayload,omitempty"`
	capsolverProxyFields
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

type capsolverImageToTextTask struct {
	Type       string `json:"type"`
	Body       string `json:"body"`
	WebsiteURL string `json:"websiteURL,omitempty"`
	Module     string `json:"module,omitempty"`
}

// mapToCapSolverTask converts a universal task to CapSolver format
func mapToCapSolverTask(task unicap.Task) (any, error) {
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
	case *tasks.HCaptchaTask:
		return mapHCaptcha(t), nil
	case *tasks.ImageToTextTask:
		return mapImageToText(t), nil
	default:
		return nil, fmt.Errorf("unsupported task type: %s", task.Type())
	}
}

func mapReCaptchaV2(task *tasks.ReCaptchaV2Task) capsolverReCaptchaV2Task {
	result := capsolverReCaptchaV2Task{
		capsolverReCaptchaV2Base: capsolverReCaptchaV2Base{
			Type:                "ReCaptchaV2TaskProxyLess",
			WebsiteURL:          task.WebsiteURL,
			WebsiteKey:          task.WebsiteKey,
			IsInvisible:         task.IsInvisible,
			PageAction:          task.PageAction,
			RecaptchaDataSValue: task.DataS,
			IsSession:           task.IsSession,
			ApiDomain:           task.APIDomain,
		},
		EnterprisePayload: task.EnterprisePayload,
	}

	if task.Proxy.IsSet() {
		result.Type = "ReCaptchaV2Task"
		result.capsolverProxyFields = capsolverProxyFields{
			ProxyType:     string(task.Proxy.Type),
			ProxyAddress:  task.Proxy.Address,
			ProxyPort:     task.Proxy.Port,
			ProxyLogin:    task.Proxy.Login,
			ProxyPassword: task.Proxy.Password,
		}
	}

	return result
}

func mapReCaptchaV3(task *tasks.ReCaptchaV3Task) capsolverReCaptchaV3Task {
	result := capsolverReCaptchaV3Task{
		capsolverReCaptchaV3Base: capsolverReCaptchaV3Base{
			Type:       "ReCaptchaV3TaskProxyLess",
			WebsiteURL: task.WebsiteURL,
			WebsiteKey: task.WebsiteKey,
			PageAction: task.PageAction,
			IsSession:  task.IsSession,
			ApiDomain:  task.APIDomain,
		},
		EnterprisePayload: task.EnterprisePayload,
	}

	if task.Proxy.IsSet() {
		result.Type = "ReCaptchaV3Task"
		result.capsolverProxyFields = capsolverProxyFields{
			ProxyType:     string(task.Proxy.Type),
			ProxyAddress:  task.Proxy.Address,
			ProxyPort:     task.Proxy.Port,
			ProxyLogin:    task.Proxy.Login,
			ProxyPassword: task.Proxy.Password,
		}
	}

	return result
}

func mapReCaptchaV2Enterprise(task *tasks.ReCaptchaV2EnterpriseTask) capsolverReCaptchaV2EnterpriseTask {
	result := capsolverReCaptchaV2EnterpriseTask{
		capsolverReCaptchaV2Base: capsolverReCaptchaV2Base{
			Type:                "ReCaptchaV2EnterpriseTaskProxyLess",
			WebsiteURL:          task.WebsiteURL,
			WebsiteKey:          task.WebsiteKey,
			IsInvisible:         task.IsInvisible,
			PageAction:          task.PageAction,
			RecaptchaDataSValue: task.DataS,
			IsSession:           task.IsSession,
			ApiDomain:           task.APIDomain,
		},
		EnterprisePayload: task.EnterprisePayload,
	}

	if task.Proxy.IsSet() {
		result.Type = "ReCaptchaV2EnterpriseTask"
		result.capsolverProxyFields = capsolverProxyFields{
			ProxyType:     string(task.Proxy.Type),
			ProxyAddress:  task.Proxy.Address,
			ProxyPort:     task.Proxy.Port,
			ProxyLogin:    task.Proxy.Login,
			ProxyPassword: task.Proxy.Password,
		}
	}

	return result
}

func mapReCaptchaV3Enterprise(task *tasks.ReCaptchaV3EnterpriseTask) capsolverReCaptchaV3EnterpriseTask {
	result := capsolverReCaptchaV3EnterpriseTask{
		capsolverReCaptchaV3Base: capsolverReCaptchaV3Base{
			Type:       "ReCaptchaV3EnterpriseTaskProxyLess",
			WebsiteURL: task.WebsiteURL,
			WebsiteKey: task.WebsiteKey,
			PageAction: task.PageAction,
			IsSession:  task.IsSession,
			ApiDomain:  task.APIDomain,
		},
		EnterprisePayload: task.EnterprisePayload,
	}

	if task.Proxy.IsSet() {
		result.Type = "ReCaptchaV3EnterpriseTask"
		result.capsolverProxyFields = capsolverProxyFields{
			ProxyType:     string(task.Proxy.Type),
			ProxyAddress:  task.Proxy.Address,
			ProxyPort:     task.Proxy.Port,
			ProxyLogin:    task.Proxy.Login,
			ProxyPassword: task.Proxy.Password,
		}
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
		result.capsolverProxyFields = capsolverProxyFields{
			ProxyType:     string(task.Proxy.Type),
			ProxyAddress:  task.Proxy.Address,
			ProxyPort:     task.Proxy.Port,
			ProxyLogin:    task.Proxy.Login,
			ProxyPassword: task.Proxy.Password,
		}
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
		result.capsolverProxyFields = capsolverProxyFields{
			ProxyType:     string(task.Proxy.Type),
			ProxyAddress:  task.Proxy.Address,
			ProxyPort:     task.Proxy.Port,
			ProxyLogin:    task.Proxy.Login,
			ProxyPassword: task.Proxy.Password,
		}
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
		result.capsolverProxyFields = capsolverProxyFields{
			ProxyType:     string(task.Proxy.Type),
			ProxyAddress:  task.Proxy.Address,
			ProxyPort:     task.Proxy.Port,
			ProxyLogin:    task.Proxy.Login,
			ProxyPassword: task.Proxy.Password,
		}
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
		result.capsolverProxyFields = capsolverProxyFields{
			ProxyType:     string(task.Proxy.Type),
			ProxyAddress:  task.Proxy.Address,
			ProxyPort:     task.Proxy.Port,
			ProxyLogin:    task.Proxy.Login,
			ProxyPassword: task.Proxy.Password,
		}
	}

	return result
}

func mapHCaptcha(task *tasks.HCaptchaTask) capsolverHCaptchaTask {
	result := capsolverHCaptchaTask{
		Type:              "HCaptchaTaskProxyLess",
		WebsiteURL:        task.WebsiteURL,
		WebsiteKey:        task.WebsiteKey,
		IsInvisible:       task.IsInvisible,
		EnterprisePayload: task.EnterpriseData,
	}

	if task.Proxy.IsSet() {
		result.Type = "HCaptchaTask"
		result.capsolverProxyFields = capsolverProxyFields{
			ProxyType:     string(task.Proxy.Type),
			ProxyAddress:  task.Proxy.Address,
			ProxyPort:     task.Proxy.Port,
			ProxyLogin:    task.Proxy.Login,
			ProxyPassword: task.Proxy.Password,
		}
	}

	return result
}

func mapImageToText(task *tasks.ImageToTextTask) capsolverImageToTextTask {
	return capsolverImageToTextTask{
		Type:       "ImageToTextTask",
		Body:       task.Body,
		WebsiteURL: task.WebsiteURL,
		Module:     task.Module,
	}
}
