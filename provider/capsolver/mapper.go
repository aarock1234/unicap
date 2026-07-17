package capsolver

import (
	"fmt"

	"github.com/aarock1234/unicap"
	"github.com/aarock1234/unicap/internal/solverapi"
	"github.com/aarock1234/unicap/tasks"
)

type reCaptchaV2Base struct {
	Type                string `json:"type"`
	WebsiteURL          string `json:"websiteURL"`
	WebsiteKey          string `json:"websiteKey"`
	IsInvisible         bool   `json:"isInvisible,omitempty"`
	PageAction          string `json:"pageAction,omitempty"`
	RecaptchaDataSValue string `json:"recaptchaDataSValue,omitempty"`
	IsSession           bool   `json:"isSession,omitempty"`
	APIDomain           string `json:"apiDomain,omitempty"`
}

// reCaptchaV2Task carries both standard and enterprise reCAPTCHA v2 payloads;
// enterprise mode is distinguished solely by the Type string.
type reCaptchaV2Task struct {
	reCaptchaV2Base
	EnterprisePayload map[string]any `json:"enterprisePayload,omitempty"`
	solverapi.ProxyFields
}

type reCaptchaV3Base struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	WebsiteKey string `json:"websiteKey"`
	PageAction string `json:"pageAction,omitempty"`
	IsSession  bool   `json:"isSession,omitempty"`
	APIDomain  string `json:"apiDomain,omitempty"`
}

// reCaptchaV3Task carries both standard and enterprise reCAPTCHA v3 payloads;
// enterprise mode is distinguished solely by the Type string.
type reCaptchaV3Task struct {
	reCaptchaV3Base
	EnterprisePayload map[string]any `json:"enterprisePayload,omitempty"`
	solverapi.ProxyFields
}

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
	Type       string            `json:"type"`
	WebsiteURL string            `json:"websiteURL"`
	WebsiteKey string            `json:"websiteKey"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	solverapi.ProxyFields
}

type geeTestTask struct {
	Type                      string `json:"type"`
	WebsiteURL                string `json:"websiteURL"`
	GT                        string `json:"gt"`
	Challenge                 string `json:"challenge"`
	GeetestAPIServerSubdomain string `json:"geetestApiServerSubdomain,omitempty"`
	solverapi.ProxyFields
}

type geeTestV4Task struct {
	Type                      string `json:"type"`
	WebsiteURL                string `json:"websiteURL"`
	CaptchaID                 string `json:"captchaId"`
	GeetestAPIServerSubdomain string `json:"geetestApiServerSubdomain,omitempty"`
	solverapi.ProxyFields
}

type hCaptchaTask struct {
	Type              string         `json:"type"`
	WebsiteURL        string         `json:"websiteURL"`
	WebsiteKey        string         `json:"websiteKey"`
	IsInvisible       bool           `json:"isInvisible,omitempty"`
	EnterprisePayload map[string]any `json:"enterprisePayload,omitempty"`
	UserAgent         string         `json:"userAgent,omitempty"`
	solverapi.ProxyFields
}

type cloudflareChallengeTask struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	HTML       string `json:"html,omitempty"`
	UserAgent  string `json:"userAgent,omitempty"`
	Proxy      string `json:"proxy"`
}

type dataDomeTask struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL,omitempty"`
	CaptchaURL string `json:"captchaUrl"`
	UserAgent  string `json:"userAgent"`
	Proxy      string `json:"proxy"`
}

type imageToTextTask struct {
	Type       string `json:"type"`
	Body       string `json:"body"`
	WebsiteURL string `json:"websiteURL,omitempty"`
	Module     string `json:"module,omitempty"`
}

type awsWAFTask struct {
	Type           string `json:"type"`
	WebsiteURL     string `json:"websiteURL"`
	AWSKey         string `json:"awsKey,omitempty"`
	AWSIV          string `json:"awsIv,omitempty"`
	AWSContext     string `json:"awsContext,omitempty"`
	AWSChallengeJS string `json:"awsChallengeJS,omitempty"`
	AWSAPIJS       string `json:"awsApiJs,omitempty"`
	Proxy          string `json:"proxy,omitempty"`
}

type mtCaptchaTask struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	WebsiteKey string `json:"websiteKey"`
	solverapi.ProxyFields
}

// mapTask converts a universal task into the CapSolver task format.
func mapTask(task unicap.Task) (any, error) {
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
	case *tasks.AWSWAFTask:
		return mapAWSWAF(t), nil
	case *tasks.MTCaptchaTask:
		return mapMTCaptcha(t), nil
	default:
		return nil, fmt.Errorf("%s: %w", task.Type(), unicap.ErrUnsupportedTask)
	}
}

func mapReCaptchaV2(task *tasks.ReCaptchaV2Task) reCaptchaV2Task {
	result := reCaptchaV2Task{
		reCaptchaV2Base: reCaptchaV2Base{
			Type:                "ReCaptchaV2TaskProxyLess",
			WebsiteURL:          task.WebsiteURL,
			WebsiteKey:          task.WebsiteKey,
			IsInvisible:         task.IsInvisible,
			PageAction:          task.PageAction,
			RecaptchaDataSValue: task.DataS,
			IsSession:           task.IsSession,
			APIDomain:           task.APIDomain,
		},
		EnterprisePayload: task.EnterprisePayload,
	}

	if task.Proxy.IsSet() {
		result.Type = "ReCaptchaV2Task"
		result.ProxyFields = solverapi.ProxyFieldsFrom(task.Proxy)
	}

	return result
}

func mapReCaptchaV3(task *tasks.ReCaptchaV3Task) reCaptchaV3Task {
	result := reCaptchaV3Task{
		reCaptchaV3Base: reCaptchaV3Base{
			Type:       "ReCaptchaV3TaskProxyLess",
			WebsiteURL: task.WebsiteURL,
			WebsiteKey: task.WebsiteKey,
			PageAction: task.PageAction,
			IsSession:  task.IsSession,
			APIDomain:  task.APIDomain,
		},
		EnterprisePayload: task.EnterprisePayload,
	}

	if task.Proxy.IsSet() {
		result.Type = "ReCaptchaV3Task"
		result.ProxyFields = solverapi.ProxyFieldsFrom(task.Proxy)
	}

	return result
}

func mapReCaptchaV2Enterprise(task *tasks.ReCaptchaV2EnterpriseTask) reCaptchaV2Task {
	result := reCaptchaV2Task{
		reCaptchaV2Base: reCaptchaV2Base{
			Type:                "ReCaptchaV2EnterpriseTaskProxyLess",
			WebsiteURL:          task.WebsiteURL,
			WebsiteKey:          task.WebsiteKey,
			IsInvisible:         task.IsInvisible,
			PageAction:          task.PageAction,
			RecaptchaDataSValue: task.DataS,
			IsSession:           task.IsSession,
			APIDomain:           task.APIDomain,
		},
		EnterprisePayload: task.EnterprisePayload,
	}

	if task.Proxy.IsSet() {
		result.Type = "ReCaptchaV2EnterpriseTask"
		result.ProxyFields = solverapi.ProxyFieldsFrom(task.Proxy)
	}

	return result
}

func mapReCaptchaV3Enterprise(task *tasks.ReCaptchaV3EnterpriseTask) reCaptchaV3Task {
	result := reCaptchaV3Task{
		reCaptchaV3Base: reCaptchaV3Base{
			Type:       "ReCaptchaV3EnterpriseTaskProxyLess",
			WebsiteURL: task.WebsiteURL,
			WebsiteKey: task.WebsiteKey,
			PageAction: task.PageAction,
			IsSession:  task.IsSession,
			APIDomain:  task.APIDomain,
		},
		EnterprisePayload: task.EnterprisePayload,
	}

	if task.Proxy.IsSet() {
		result.Type = "ReCaptchaV3EnterpriseTask"
		result.ProxyFields = solverapi.ProxyFieldsFrom(task.Proxy)
	}

	return result
}

func mapFunCaptcha(task *tasks.FunCaptchaTask) funCaptchaTask {
	result := funCaptchaTask{
		Type:                     "FunCaptchaTaskProxyLess",
		WebsiteURL:               task.WebsiteURL,
		WebsitePublicKey:         task.WebsitePublicKey,
		FuncaptchaAPIJSSubdomain: task.APIJSSubdomain,
		Data:                     task.Data,
		UserAgent:                task.UserAgent,
	}

	if task.Proxy.IsSet() {
		result.Type = "FunCaptchaTask"
		result.ProxyFields = solverapi.ProxyFieldsFrom(task.Proxy)
	}

	return result
}

func mapTurnstile(task *tasks.TurnstileTask) turnstileTask {
	result := turnstileTask{
		Type:       "AntiTurnstileTaskProxyLess",
		WebsiteURL: task.WebsiteURL,
		WebsiteKey: task.WebsiteKey,
	}

	if task.Action != "" || task.CData != "" {
		result.Metadata = make(map[string]string, 2)
		if task.Action != "" {
			result.Metadata["action"] = task.Action
		}
		if task.CData != "" {
			result.Metadata["cdata"] = task.CData
		}
	}

	if task.Proxy.IsSet() {
		result.Type = "AntiTurnstileTask"
		result.ProxyFields = solverapi.ProxyFieldsFrom(task.Proxy)
	}

	return result
}

func mapCloudflareChallenge(task *tasks.CloudflareChallengeTask) cloudflareChallengeTask {
	return cloudflareChallengeTask{
		Type:       "AntiCloudflareTask",
		WebsiteURL: task.WebsiteURL,
		HTML:       task.HTML,
		UserAgent:  task.UserAgent,
		Proxy:      solverapi.ProxyString(task.Proxy),
	}
}

func mapDataDome(task *tasks.DataDomeTask) dataDomeTask {
	return dataDomeTask{
		Type:       "DatadomeSliderTask",
		WebsiteURL: task.WebsiteURL,
		CaptchaURL: task.CaptchaURL,
		UserAgent:  task.UserAgent,
		Proxy:      solverapi.ProxyString(task.Proxy),
	}
}

func mapGeeTest(task *tasks.GeeTestTask) geeTestTask {
	result := geeTestTask{
		Type:                      "GeeTestTaskProxyLess",
		WebsiteURL:                task.WebsiteURL,
		GT:                        task.GT,
		Challenge:                 task.Challenge,
		GeetestAPIServerSubdomain: task.APIServerSubdomain,
	}

	if task.Proxy.IsSet() {
		result.Type = "GeeTestTask"
		result.ProxyFields = solverapi.ProxyFieldsFrom(task.Proxy)
	}

	return result
}

func mapGeeTestV4(task *tasks.GeeTestV4Task) geeTestV4Task {
	result := geeTestV4Task{
		Type:                      "GeeTestTaskProxyLess",
		WebsiteURL:                task.WebsiteURL,
		CaptchaID:                 task.CaptchaID,
		GeetestAPIServerSubdomain: task.APIServerSubdomain,
	}

	if task.Proxy.IsSet() {
		result.Type = "GeeTestTask"
		result.ProxyFields = solverapi.ProxyFieldsFrom(task.Proxy)
	}

	return result
}

func mapHCaptcha(task *tasks.HCaptchaTask) hCaptchaTask {
	result := hCaptchaTask{
		Type:              "HCaptchaTaskProxyLess",
		WebsiteURL:        task.WebsiteURL,
		WebsiteKey:        task.WebsiteKey,
		IsInvisible:       task.IsInvisible,
		EnterprisePayload: task.EnterpriseData,
		UserAgent:         task.UserAgent,
	}

	if task.Proxy.IsSet() {
		result.Type = "HCaptchaTask"
		result.ProxyFields = solverapi.ProxyFieldsFrom(task.Proxy)
	}

	return result
}

func mapImageToText(task *tasks.ImageToTextTask) imageToTextTask {
	return imageToTextTask{
		Type:       "ImageToTextTask",
		Body:       task.Body,
		WebsiteURL: task.WebsiteURL,
		Module:     task.Module,
	}
}

func mapAWSWAF(task *tasks.AWSWAFTask) awsWAFTask {
	result := awsWAFTask{
		Type:           "AntiAwsWafTaskProxyLess",
		WebsiteURL:     task.WebsiteURL,
		AWSKey:         task.Key,
		AWSIV:          task.IV,
		AWSContext:     task.Context,
		AWSChallengeJS: task.ChallengeScript,
		AWSAPIJS:       task.CaptchaScript,
	}

	if task.Proxy.IsSet() {
		result.Type = "AntiAwsWafTask"
		result.Proxy = solverapi.ProxyString(task.Proxy)
	}

	return result
}

func mapMTCaptcha(task *tasks.MTCaptchaTask) mtCaptchaTask {
	result := mtCaptchaTask{
		Type:       "MtCaptchaTaskProxyLess",
		WebsiteURL: task.WebsiteURL,
		WebsiteKey: task.WebsiteKey,
	}

	if task.Proxy.IsSet() {
		result.Type = "MtCaptchaTask"
		result.ProxyFields = solverapi.ProxyFieldsFrom(task.Proxy)
	}

	return result
}
