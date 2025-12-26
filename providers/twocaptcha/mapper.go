package twocaptcha

import (
	"fmt"

	"github.com/aarock1234/unicap/unicap"
	"github.com/aarock1234/unicap/unicap/tasks"
)

// Task types for 2Captcha API

// Base proxy fields for all tasks
type twocaptchaProxyFields struct {
	ProxyType     string `json:"proxyType,omitempty"`
	ProxyAddress  string `json:"proxyAddress,omitempty"`
	ProxyPort     int    `json:"proxyPort,omitempty"`
	ProxyLogin    string `json:"proxyLogin,omitempty"`
	ProxyPassword string `json:"proxyPassword,omitempty"`
}

type twocaptchaReCaptchaV2Task struct {
	Type                string `json:"type"`
	WebsiteURL          string `json:"websiteURL"`
	WebsiteKey          string `json:"websiteKey"`
	IsInvisible         bool   `json:"isInvisible,omitempty"`
	RecaptchaDataSValue string `json:"recaptchaDataSValue,omitempty"`
	UserAgent           string `json:"userAgent,omitempty"`
	Cookies             string `json:"cookies,omitempty"`
	APIDomain           string `json:"apiDomain,omitempty"`
	twocaptchaProxyFields
}

type twocaptchaReCaptchaV3Task struct {
	Type         string  `json:"type"`
	WebsiteURL   string  `json:"websiteURL"`
	WebsiteKey   string  `json:"websiteKey"`
	PageAction   string  `json:"pageAction,omitempty"`
	MinScore     float64 `json:"minScore"`
	IsEnterprise bool    `json:"isEnterprise,omitempty"`
	APIDomain    string  `json:"apiDomain,omitempty"`
	twocaptchaProxyFields
}

type twocaptchaReCaptchaV2EnterpriseTask struct {
	Type                string         `json:"type"`
	WebsiteURL          string         `json:"websiteURL"`
	WebsiteKey          string         `json:"websiteKey"`
	IsInvisible         bool           `json:"isInvisible,omitempty"`
	PageAction          string         `json:"pageAction,omitempty"`
	RecaptchaDataSValue string         `json:"recaptchaDataSValue,omitempty"`
	EnterprisePayload   map[string]any `json:"enterprisePayload,omitempty"`
	ApiDomain           string         `json:"apiDomain,omitempty"`
	UserAgent           string         `json:"userAgent,omitempty"`
	Cookies             string         `json:"cookies,omitempty"`
	twocaptchaProxyFields
}

type twocaptchaReCaptchaV3EnterpriseTask struct {
	Type              string         `json:"type"`
	WebsiteURL        string         `json:"websiteURL"`
	WebsiteKey        string         `json:"websiteKey"`
	PageAction        string         `json:"pageAction,omitempty"`
	MinScore          float64        `json:"minScore"`
	EnterprisePayload map[string]any `json:"enterprisePayload,omitempty"`
	ApiDomain         string         `json:"apiDomain,omitempty"`
	twocaptchaProxyFields
}

type twocaptchaFunCaptchaTask struct {
	Type                     string `json:"type"`
	WebsiteURL               string `json:"websiteURL"`
	WebsitePublicKey         string `json:"websitePublicKey"`
	FuncaptchaApiJSSubdomain string `json:"funcaptchaApiJSSubdomain,omitempty"`
	Data                     string `json:"data,omitempty"`
	UserAgent                string `json:"userAgent,omitempty"`
	twocaptchaProxyFields
}

type twocaptchaTurnstileTask struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	WebsiteKey string `json:"websiteKey"`
	Action     string `json:"action,omitempty"`
	Data       string `json:"data,omitempty"`
	PageData   string `json:"pagedata,omitempty"`
	twocaptchaProxyFields
}

type twocaptchaGeeTestTask struct {
	Type                      string      `json:"type"`
	WebsiteURL                string      `json:"websiteURL"`
	GT                        string      `json:"gt,omitempty"`
	Challenge                 string      `json:"challenge,omitempty"`
	Version                   int         `json:"version,omitempty"`
	InitParameters            geetestInit `json:"initParameters,omitempty"`
	GeetestApiServerSubdomain string      `json:"geetestApiServerSubdomain,omitempty"`
	UserAgent                 string      `json:"userAgent,omitempty"`
	twocaptchaProxyFields
}

type geetestInit struct {
	CaptchaID string `json:"captcha_id,omitempty"`
}

type twocaptchaDataDomeTask struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	CaptchaURL string `json:"captchaUrl"`
	UserAgent  string `json:"userAgent"`
	twocaptchaProxyFields
}

type twocaptchaHCaptchaTask struct {
	Type              string         `json:"type"`
	WebsiteURL        string         `json:"websiteURL"`
	WebsiteKey        string         `json:"websiteKey"`
	IsInvisible       bool           `json:"isInvisible,omitempty"`
	EnterprisePayload map[string]any `json:"enterprisePayload,omitempty"`
	UserAgent         string         `json:"userAgent,omitempty"`
	Cookies           string         `json:"cookies,omitempty"`
	twocaptchaProxyFields
}

type twocaptchaImageToTextTask struct {
	Type            string `json:"type"`
	Body            string `json:"body"`
	Phrase          bool   `json:"phrase,omitempty"`
	Case            bool   `json:"case,omitempty"`
	Numeric         int    `json:"numeric,omitempty"`
	Math            bool   `json:"math,omitempty"`
	MinLength       int    `json:"minLength,omitempty"`
	MaxLength       int    `json:"maxLength,omitempty"`
	Comment         string `json:"comment,omitempty"`
	ImgInstructions string `json:"imgInstructions,omitempty"`
}

// mapToTwoCaptchaTask converts a universal task to 2Captcha format
func mapToTwoCaptchaTask(task unicap.Task) (any, error) {
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

func mapReCaptchaV2(task *tasks.ReCaptchaV2Task) twocaptchaReCaptchaV2Task {
	result := twocaptchaReCaptchaV2Task{
		Type:                "RecaptchaV2TaskProxyless",
		WebsiteURL:          task.WebsiteURL,
		WebsiteKey:          task.WebsiteKey,
		IsInvisible:         task.IsInvisible,
		RecaptchaDataSValue: task.DataS,
		UserAgent:           task.UserAgent,
		Cookies:             task.Cookies,
		APIDomain:           task.APIDomain,
	}

	if task.Proxy.IsSet() {
		result.Type = "RecaptchaV2Task"
		result.twocaptchaProxyFields = twocaptchaProxyFields{
			ProxyType:     string(task.Proxy.Type),
			ProxyAddress:  task.Proxy.Address,
			ProxyPort:     task.Proxy.Port,
			ProxyLogin:    task.Proxy.Login,
			ProxyPassword: task.Proxy.Password,
		}
	}

	return result
}

func mapReCaptchaV3(task *tasks.ReCaptchaV3Task) twocaptchaReCaptchaV3Task {
	result := twocaptchaReCaptchaV3Task{
		Type:         "RecaptchaV3TaskProxyless",
		WebsiteURL:   task.WebsiteURL,
		WebsiteKey:   task.WebsiteKey,
		PageAction:   task.PageAction,
		MinScore:     task.MinScore,
		IsEnterprise: task.IsEnterprise,
		APIDomain:    task.APIDomain,
	}

	if task.Proxy.IsSet() {
		result.Type = "RecaptchaV3Task"
		result.twocaptchaProxyFields = twocaptchaProxyFields{
			ProxyType:     string(task.Proxy.Type),
			ProxyAddress:  task.Proxy.Address,
			ProxyPort:     task.Proxy.Port,
			ProxyLogin:    task.Proxy.Login,
			ProxyPassword: task.Proxy.Password,
		}
	}

	return result
}

func mapReCaptchaV2Enterprise(task *tasks.ReCaptchaV2EnterpriseTask) twocaptchaReCaptchaV2EnterpriseTask {
	result := twocaptchaReCaptchaV2EnterpriseTask{
		Type:                "RecaptchaV2EnterpriseTaskProxyless",
		WebsiteURL:          task.WebsiteURL,
		WebsiteKey:          task.WebsiteKey,
		IsInvisible:         task.IsInvisible,
		PageAction:          task.PageAction,
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

func mapReCaptchaV3Enterprise(task *tasks.ReCaptchaV3EnterpriseTask) twocaptchaReCaptchaV3EnterpriseTask {
	result := twocaptchaReCaptchaV3EnterpriseTask{
		Type:              "RecaptchaV3EnterpriseTaskProxyless",
		WebsiteURL:        task.WebsiteURL,
		WebsiteKey:        task.WebsiteKey,
		PageAction:        task.PageAction,
		MinScore:          task.MinScore,
		EnterprisePayload: task.EnterprisePayload,
		ApiDomain:         task.APIDomain,
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

func mapFunCaptcha(task *tasks.FunCaptchaTask) twocaptchaFunCaptchaTask {
	result := twocaptchaFunCaptchaTask{
		Type:                     "FunCaptchaTaskProxyless",
		WebsiteURL:               task.WebsiteURL,
		WebsitePublicKey:         task.WebsitePublicKey,
		FuncaptchaApiJSSubdomain: task.APIJSSubdomain,
		Data:                     task.Data,
		UserAgent:                task.UserAgent,
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

func mapTurnstile(task *tasks.TurnstileTask) twocaptchaTurnstileTask {
	result := twocaptchaTurnstileTask{
		Type:       "TurnstileTaskProxyless",
		WebsiteURL: task.WebsiteURL,
		WebsiteKey: task.WebsiteKey,
		Action:     task.Action,
		Data:       task.CData,
		PageData:   task.PageData,
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

func mapDataDome(task *tasks.DataDomeTask) twocaptchaDataDomeTask {
	return twocaptchaDataDomeTask{
		Type:       "DataDomeSliderTask",
		WebsiteURL: task.WebsiteURL,
		CaptchaURL: task.CaptchaURL,
		UserAgent:  task.UserAgent,
		twocaptchaProxyFields: twocaptchaProxyFields{
			ProxyType:     string(task.Proxy.Type),
			ProxyAddress:  task.Proxy.Address,
			ProxyPort:     task.Proxy.Port,
			ProxyLogin:    task.Proxy.Login,
			ProxyPassword: task.Proxy.Password,
		},
	}
}

func mapGeeTest(task *tasks.GeeTestTask) twocaptchaGeeTestTask {
	result := twocaptchaGeeTestTask{
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

func mapGeeTestV4(task *tasks.GeeTestV4Task) twocaptchaGeeTestTask {
	result := twocaptchaGeeTestTask{
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

func mapHCaptcha(task *tasks.HCaptchaTask) twocaptchaHCaptchaTask {
	result := twocaptchaHCaptchaTask{
		Type:              "HCaptchaTaskProxyless",
		WebsiteURL:        task.WebsiteURL,
		WebsiteKey:        task.WebsiteKey,
		IsInvisible:       task.IsInvisible,
		EnterprisePayload: task.EnterpriseData,
		UserAgent:         task.UserAgent,
		Cookies:           task.Cookies,
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

func mapImageToText(task *tasks.ImageToTextTask) twocaptchaImageToTextTask {
	return twocaptchaImageToTextTask{
		Type:            "ImageToTextTask",
		Body:            task.Body,
		Phrase:          task.Phrase,
		Case:            task.Case,
		Numeric:         int(task.Numeric),
		Math:            task.Math,
		MinLength:       task.MinLength,
		MaxLength:       task.MaxLength,
		Comment:         task.Comment,
		ImgInstructions: task.ImgInstructions,
	}
}
