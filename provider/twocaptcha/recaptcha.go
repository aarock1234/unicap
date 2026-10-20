package twocaptcha

import (
	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

type reCaptchaV2Task struct {
	Type                string         `json:"type"`
	WebsiteURL          string         `json:"websiteURL"`
	WebsiteKey          string         `json:"websiteKey"`
	IsInvisible         bool           `json:"isInvisible,omitempty"`
	RecaptchaDataSValue string         `json:"recaptchaDataSValue,omitempty"`
	EnterprisePayload   map[string]any `json:"enterprisePayload,omitempty"`
	UserAgent           string         `json:"userAgent,omitempty"`
	Cookies             string         `json:"cookies,omitempty"`
	APIDomain           string         `json:"apiDomain,omitempty"`
	solverapi.ProxyFields
}

type reCaptchaV3Task struct {
	Type         string  `json:"type"`
	WebsiteURL   string  `json:"websiteURL"`
	WebsiteKey   string  `json:"websiteKey"`
	MinScore     float64 `json:"minScore"`
	PageAction   string  `json:"pageAction,omitempty"`
	IsEnterprise bool    `json:"isEnterprise,omitempty"`
	APIDomain    string  `json:"apiDomain,omitempty"`
}

// reCaptchaV3Variant is 2Captcha's only reCAPTCHA v3 wire type; it covers
// Enterprise through isEnterprise and has no proxy variant.
var reCaptchaV3Variant = proxyless("RecaptchaV3TaskProxyless")

func mapReCaptchaV2(task *tasks.ReCaptchaV2Task) (reCaptchaV2Task, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("PageAction", task.PageAction != ""),
		solverapi.Setting("IsSession", task.IsSession),
	); err != nil {
		return reCaptchaV2Task{}, err
	}

	taskType, err := variant("RecaptchaV2Task").For(task.Proxy)
	if err != nil {
		return reCaptchaV2Task{}, err
	}

	return reCaptchaV2Task{
		Type:                taskType,
		WebsiteURL:          task.WebsiteURL,
		WebsiteKey:          task.WebsiteKey,
		IsInvisible:         task.IsInvisible,
		RecaptchaDataSValue: task.DataS,
		UserAgent:           task.UserAgent,
		Cookies:             task.Cookies,
		APIDomain:           task.APIDomain,
		ProxyFields:         solverapi.NewProxyFields(task.Proxy),
	}, nil
}

func mapReCaptchaV2Enterprise(task *tasks.ReCaptchaV2EnterpriseTask) (reCaptchaV2Task, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("PageAction", task.PageAction != ""),
		solverapi.Setting("IsSession", task.IsSession),
	); err != nil {
		return reCaptchaV2Task{}, err
	}

	taskType, err := variant("RecaptchaV2EnterpriseTask").For(task.Proxy)
	if err != nil {
		return reCaptchaV2Task{}, err
	}

	return reCaptchaV2Task{
		Type:              taskType,
		WebsiteURL:        task.WebsiteURL,
		WebsiteKey:        task.WebsiteKey,
		IsInvisible:       task.IsInvisible,
		EnterprisePayload: task.EnterprisePayload,
		UserAgent:         task.UserAgent,
		Cookies:           task.Cookies,
		APIDomain:         task.APIDomain,
		ProxyFields:       solverapi.NewProxyFields(task.Proxy),
	}, nil
}

func mapReCaptchaV3(task *tasks.ReCaptchaV3Task) (reCaptchaV3Task, error) {
	if err := solverapi.Check(
		solverapi.Required(
			solverapi.Setting("MinScore", task.MinScore != 0),
		),
		solverapi.Unsupported(
			solverapi.Setting("HighScore", task.HighScore),
			solverapi.Setting("IsSession", task.IsSession),
		),
	); err != nil {
		return reCaptchaV3Task{}, err
	}

	taskType, err := reCaptchaV3Variant.For(task.Proxy)
	if err != nil {
		return reCaptchaV3Task{}, err
	}

	return reCaptchaV3Task{
		Type:       taskType,
		WebsiteURL: task.WebsiteURL,
		WebsiteKey: task.WebsiteKey,
		MinScore:   task.MinScore,
		PageAction: task.PageAction,
		APIDomain:  task.APIDomain,
	}, nil
}

func mapReCaptchaV3Enterprise(task *tasks.ReCaptchaV3EnterpriseTask) (reCaptchaV3Task, error) {
	if err := solverapi.Check(
		solverapi.Required(
			solverapi.Setting("MinScore", task.MinScore != 0),
		),
		solverapi.Unsupported(
			solverapi.Setting("HighScore", task.HighScore),
			solverapi.Setting("IsSession", task.IsSession),
			solverapi.Setting("EnterprisePayload", len(task.EnterprisePayload) > 0),
		),
	); err != nil {
		return reCaptchaV3Task{}, err
	}

	taskType, err := reCaptchaV3Variant.For(task.Proxy)
	if err != nil {
		return reCaptchaV3Task{}, err
	}

	return reCaptchaV3Task{
		Type:         taskType,
		WebsiteURL:   task.WebsiteURL,
		WebsiteKey:   task.WebsiteKey,
		MinScore:     task.MinScore,
		PageAction:   task.PageAction,
		IsEnterprise: true,
		APIDomain:    task.APIDomain,
	}, nil
}
