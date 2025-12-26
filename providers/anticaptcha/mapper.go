package anticaptcha

import (
	"fmt"

	"github.com/aarock1234/unicap/unicap"
	"github.com/aarock1234/unicap/unicap/tasks"
)

// Task types for AntiCaptcha API

// Base proxy fields for all tasks
type anticaptchaProxyFields struct {
	ProxyType     string `json:"proxyType,omitempty"`
	ProxyAddress  string `json:"proxyAddress,omitempty"`
	ProxyPort     int    `json:"proxyPort,omitempty"`
	ProxyLogin    string `json:"proxyLogin,omitempty"`
	ProxyPassword string `json:"proxyPassword,omitempty"`
}

type anticaptchaReCaptchaV2Task struct {
	Type                string `json:"type"`
	WebsiteURL          string `json:"websiteURL"`
	WebsiteKey          string `json:"websiteKey"`
	IsInvisible         bool   `json:"isInvisible,omitempty"`
	RecaptchaDataSValue string `json:"recaptchaDataSValue,omitempty"`
	anticaptchaProxyFields
}

type anticaptchaReCaptchaV3Task struct {
	Type         string  `json:"type"`
	WebsiteURL   string  `json:"websiteURL"`
	WebsiteKey   string  `json:"websiteKey"`
	PageAction   string  `json:"pageAction,omitempty"`
	MinScore     float64 `json:"minScore,omitempty"`
	IsEnterprise bool    `json:"isEnterprise,omitempty"`
	ApiDomain    string  `json:"apiDomain,omitempty"`
	anticaptchaProxyFields
}

type anticaptchaReCaptchaV2EnterpriseTask struct {
	Type                string         `json:"type"`
	WebsiteURL          string         `json:"websiteURL"`
	WebsiteKey          string         `json:"websiteKey"`
	IsInvisible         bool           `json:"isInvisible,omitempty"`
	RecaptchaDataSValue string         `json:"recaptchaDataSValue,omitempty"`
	EnterprisePayload   map[string]any `json:"enterprisePayload,omitempty"`
	ApiDomain           string         `json:"apiDomain,omitempty"`
	anticaptchaProxyFields
}

type anticaptchaReCaptchaV3EnterpriseTask struct {
	Type              string         `json:"type"`
	WebsiteURL        string         `json:"websiteURL"`
	WebsiteKey        string         `json:"websiteKey"`
	PageAction        string         `json:"pageAction,omitempty"`
	MinScore          float64        `json:"minScore,omitempty"`
	EnterprisePayload map[string]any `json:"enterprisePayload,omitempty"`
	ApiDomain         string         `json:"apiDomain,omitempty"`
	IsEnterprise      bool           `json:"isEnterprise,omitempty"`
	anticaptchaProxyFields
}

type anticaptchaFunCaptchaTask struct {
	Type                     string `json:"type"`
	WebsiteURL               string `json:"websiteURL"`
	WebsitePublicKey         string `json:"websitePublicKey"`
	FuncaptchaApiJSSubdomain string `json:"funcaptchaApiJSSubdomain,omitempty"`
	Data                     string `json:"data,omitempty"`
	anticaptchaProxyFields
}

type anticaptchaTurnstileTask struct {
	Type        string `json:"type"`
	WebsiteURL  string `json:"websiteURL"`
	WebsiteKey  string `json:"websiteKey"`
	Action      string `json:"action,omitempty"`
	CData       string `json:"cData,omitempty"`
	ChlPageData string `json:"chlPageData,omitempty"`
	anticaptchaProxyFields
}

type anticaptchaGeeTestTask struct {
	Type                      string      `json:"type"`
	WebsiteURL                string      `json:"websiteURL"`
	GT                        string      `json:"gt,omitempty"`
	Challenge                 string      `json:"challenge,omitempty"`
	Version                   int         `json:"version,omitempty"`
	InitParameters            geetestInit `json:"initParameters,omitempty"`
	GeetestApiServerSubdomain string      `json:"geetestApiServerSubdomain,omitempty"`
	anticaptchaProxyFields
}

type geetestInit struct {
	CaptchaID string `json:"captcha_id,omitempty"`
}

type anticaptchaHCaptchaTask struct {
	Type              string         `json:"type"`
	WebsiteURL        string         `json:"websiteURL"`
	WebsiteKey        string         `json:"websiteKey"`
	IsInvisible       bool           `json:"isInvisible,omitempty"`
	EnterprisePayload map[string]any `json:"enterprisePayload,omitempty"`
	anticaptchaProxyFields
}

type anticaptchaImageToTextTask struct {
	Type         string `json:"type"`
	Body         string `json:"body"`
	Phrase       bool   `json:"phrase,omitempty"`
	Case         bool   `json:"case,omitempty"`
	Numeric      int    `json:"numeric,omitempty"`
	Math         bool   `json:"math,omitempty"`
	MinLength    int    `json:"minLength,omitempty"`
	MaxLength    int    `json:"maxLength,omitempty"`
	Comment      string `json:"comment,omitempty"`
	WebsiteURL   string `json:"websiteURL,omitempty"`
	LanguagePool string `json:"languagePool,omitempty"`
}

// mapToAntiCaptchaTask converts a universal task to AntiCaptcha format
func mapToAntiCaptchaTask(task unicap.Task) (any, error) {
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
	case *tasks.HCaptchaTask:
		return mapHCaptcha(t), nil
	case *tasks.ImageToTextTask:
		return mapImageToText(t), nil
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
	result := anticaptchaReCaptchaV2EnterpriseTask{
		Type:                "RecaptchaV2EnterpriseTaskProxyless",
		WebsiteURL:          task.WebsiteURL,
		WebsiteKey:          task.WebsiteKey,
		IsInvisible:         task.IsInvisible,
		RecaptchaDataSValue: task.DataS,
		EnterprisePayload:   task.EnterprisePayload,
		ApiDomain:           task.APIDomain,
	}

	if task.Proxy.IsSet() {
		result.Type = "RecaptchaV2EnterpriseTask"
		result.ProxyType = string(task.Proxy.Type)
		result.ProxyAddress = task.Proxy.Address
		result.ProxyPort = task.Proxy.Port
		result.ProxyLogin = task.Proxy.Login
		result.ProxyPassword = task.Proxy.Password
	}

	return result
}

func mapReCaptchaV3Enterprise(task *tasks.ReCaptchaV3EnterpriseTask) anticaptchaReCaptchaV3EnterpriseTask {
	result := anticaptchaReCaptchaV3EnterpriseTask{
		Type:              "RecaptchaV3EnterpriseTaskProxyless",
		WebsiteURL:        task.WebsiteURL,
		WebsiteKey:        task.WebsiteKey,
		PageAction:        task.PageAction,
		MinScore:          task.MinScore,
		EnterprisePayload: task.EnterprisePayload,
		ApiDomain:         task.APIDomain,
		IsEnterprise:      true,
	}

	if task.Proxy.IsSet() {
		result.Type = "RecaptchaV3EnterpriseTask"
		result.ProxyType = string(task.Proxy.Type)
		result.ProxyAddress = task.Proxy.Address
		result.ProxyPort = task.Proxy.Port
		result.ProxyLogin = task.Proxy.Login
		result.ProxyPassword = task.Proxy.Password
	}

	return result
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

func mapHCaptcha(task *tasks.HCaptchaTask) anticaptchaHCaptchaTask {
	result := anticaptchaHCaptchaTask{
		Type:              "HCaptchaTaskProxyless",
		WebsiteURL:        task.WebsiteURL,
		WebsiteKey:        task.WebsiteKey,
		IsInvisible:       task.IsInvisible,
		EnterprisePayload: task.EnterpriseData,
	}

	if task.Proxy.IsSet() {
		result.Type = "HCaptchaTask"
		result.ProxyType = string(task.Proxy.Type)
		result.ProxyAddress = task.Proxy.Address
		result.ProxyPort = task.Proxy.Port
		result.ProxyLogin = task.Proxy.Login
		result.ProxyPassword = task.Proxy.Password
	}

	return result
}

func mapImageToText(task *tasks.ImageToTextTask) anticaptchaImageToTextTask {
	return anticaptchaImageToTextTask{
		Type:         "ImageToTextTask",
		Body:         task.Body,
		Phrase:       task.Phrase,
		Case:         task.Case,
		Numeric:      int(task.Numeric),
		Math:         task.Math,
		MinLength:    task.MinLength,
		MaxLength:    task.MaxLength,
		Comment:      task.Comment,
		WebsiteURL:   task.WebsiteURL,
		LanguagePool: task.LanguagePool,
	}
}
