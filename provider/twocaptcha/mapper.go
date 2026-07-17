package twocaptcha

import (
	"fmt"

	"github.com/aarock1234/unicap"
	"github.com/aarock1234/unicap/internal/solverapi"
	"github.com/aarock1234/unicap/tasks"
)

type reCaptchaV2Task struct {
	Type                string `json:"type"`
	WebsiteURL          string `json:"websiteURL"`
	WebsiteKey          string `json:"websiteKey"`
	IsInvisible         bool   `json:"isInvisible,omitempty"`
	RecaptchaDataSValue string `json:"recaptchaDataSValue,omitempty"`
	UserAgent           string `json:"userAgent,omitempty"`
	Cookies             string `json:"cookies,omitempty"`
	APIDomain           string `json:"apiDomain,omitempty"`
	solverapi.ProxyFields
}

type reCaptchaV3Task struct {
	Type         string  `json:"type"`
	WebsiteURL   string  `json:"websiteURL"`
	WebsiteKey   string  `json:"websiteKey"`
	PageAction   string  `json:"pageAction,omitempty"`
	MinScore     float64 `json:"minScore,omitempty"`
	IsEnterprise bool    `json:"isEnterprise,omitempty"`
	APIDomain    string  `json:"apiDomain,omitempty"`
	solverapi.ProxyFields
}

type reCaptchaV2EnterpriseTask struct {
	Type                string         `json:"type"`
	WebsiteURL          string         `json:"websiteURL"`
	WebsiteKey          string         `json:"websiteKey"`
	IsInvisible         bool           `json:"isInvisible,omitempty"`
	PageAction          string         `json:"pageAction,omitempty"`
	RecaptchaDataSValue string         `json:"recaptchaDataSValue,omitempty"`
	EnterprisePayload   map[string]any `json:"enterprisePayload,omitempty"`
	APIDomain           string         `json:"apiDomain,omitempty"`
	UserAgent           string         `json:"userAgent,omitempty"`
	Cookies             string         `json:"cookies,omitempty"`
	solverapi.ProxyFields
}

type reCaptchaV3EnterpriseTask struct {
	Type              string         `json:"type"`
	WebsiteURL        string         `json:"websiteURL"`
	WebsiteKey        string         `json:"websiteKey"`
	PageAction        string         `json:"pageAction,omitempty"`
	MinScore          float64        `json:"minScore,omitempty"`
	EnterprisePayload map[string]any `json:"enterprisePayload,omitempty"`
	APIDomain         string         `json:"apiDomain,omitempty"`
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
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	WebsiteKey string `json:"websiteKey"`
	Action     string `json:"action,omitempty"`
	Data       string `json:"data,omitempty"`
	PageData   string `json:"pagedata,omitempty"`
	solverapi.ProxyFields
}

type geeTestTask struct {
	Type                      string       `json:"type"`
	WebsiteURL                string       `json:"websiteURL"`
	GT                        string       `json:"gt,omitempty"`
	Challenge                 string       `json:"challenge,omitempty"`
	Version                   int          `json:"version,omitempty"`
	InitParameters            *geetestInit `json:"initParameters,omitempty"`
	GeetestAPIServerSubdomain string       `json:"geetestApiServerSubdomain,omitempty"`
	UserAgent                 string       `json:"userAgent,omitempty"`
	solverapi.ProxyFields
}

type geetestInit struct {
	CaptchaID string `json:"captcha_id,omitempty"`
}

type dataDomeTask struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	CaptchaURL string `json:"captchaUrl"`
	UserAgent  string `json:"userAgent"`
	solverapi.ProxyFields
}

type hCaptchaTask struct {
	Type              string         `json:"type"`
	WebsiteURL        string         `json:"websiteURL"`
	WebsiteKey        string         `json:"websiteKey"`
	IsInvisible       bool           `json:"isInvisible,omitempty"`
	EnterprisePayload map[string]any `json:"enterprisePayload,omitempty"`
	UserAgent         string         `json:"userAgent,omitempty"`
	Cookies           string         `json:"cookies,omitempty"`
	solverapi.ProxyFields
}

type imageToTextTask struct {
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

type amazonTask struct {
	Type            string `json:"type"`
	WebsiteURL      string `json:"websiteURL"`
	WebsiteKey      string `json:"websiteKey"`
	IV              string `json:"iv,omitempty"`
	Context         string `json:"context,omitempty"`
	ChallengeScript string `json:"challengeScript,omitempty"`
	CaptchaScript   string `json:"captchaScript,omitempty"`
	solverapi.ProxyFields
}

type mtCaptchaTask struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	WebsiteKey string `json:"websiteKey"`
	solverapi.ProxyFields
}

type friendlyCaptchaTask struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	WebsiteKey string `json:"websiteKey"`
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

type textCaptchaTask struct {
	Type    string `json:"type"`
	Comment string `json:"comment"`
}

type prosopoTask struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	WebsiteKey string `json:"websiteKey"`
	solverapi.ProxyFields
}

type altchaTask struct {
	Type          string `json:"type"`
	WebsiteURL    string `json:"websiteURL"`
	ChallengeURL  string `json:"challengeURL,omitempty"`
	ChallengeJSON string `json:"challengeJSON,omitempty"`
	solverapi.ProxyFields
}

// mapTask converts a universal task into the 2Captcha task format.
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
	case *tasks.FriendlyCaptchaTask:
		return mapFriendlyCaptcha(t), nil
	case *tasks.LeminTask:
		return mapLemin(t), nil
	case *tasks.CutCaptchaTask:
		return mapCutCaptcha(t), nil
	case *tasks.TextCaptchaTask:
		return mapText(t), nil
	case *tasks.ProsopoTask:
		return mapProsopo(t), nil
	case *tasks.AltchaTask:
		return mapAltcha(t), nil
	default:
		return nil, fmt.Errorf("%s: %w", task.Type(), unicap.ErrUnsupportedTask)
	}
}

func mapReCaptchaV2(task *tasks.ReCaptchaV2Task) reCaptchaV2Task {
	result := reCaptchaV2Task{
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
		result.ProxyFields = solverapi.ProxyFieldsFrom(task.Proxy)
	}

	return result
}

func mapReCaptchaV3(task *tasks.ReCaptchaV3Task) reCaptchaV3Task {
	result := reCaptchaV3Task{
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
		result.ProxyFields = solverapi.ProxyFieldsFrom(task.Proxy)
	}

	return result
}

func mapReCaptchaV2Enterprise(task *tasks.ReCaptchaV2EnterpriseTask) reCaptchaV2EnterpriseTask {
	result := reCaptchaV2EnterpriseTask{
		Type:                "RecaptchaV2EnterpriseTaskProxyless",
		WebsiteURL:          task.WebsiteURL,
		WebsiteKey:          task.WebsiteKey,
		IsInvisible:         task.IsInvisible,
		PageAction:          task.PageAction,
		RecaptchaDataSValue: task.DataS,
		EnterprisePayload:   task.EnterprisePayload,
		APIDomain:           task.APIDomain,
	}

	if task.Proxy.IsSet() {
		result.Type = "RecaptchaV2EnterpriseTask"
		result.ProxyFields = solverapi.ProxyFieldsFrom(task.Proxy)
	}

	return result
}

func mapReCaptchaV3Enterprise(task *tasks.ReCaptchaV3EnterpriseTask) reCaptchaV3EnterpriseTask {
	result := reCaptchaV3EnterpriseTask{
		Type:              "RecaptchaV3EnterpriseTaskProxyless",
		WebsiteURL:        task.WebsiteURL,
		WebsiteKey:        task.WebsiteKey,
		PageAction:        task.PageAction,
		MinScore:          task.MinScore,
		EnterprisePayload: task.EnterprisePayload,
		APIDomain:         task.APIDomain,
	}

	if task.Proxy.IsSet() {
		result.Type = "RecaptchaV3EnterpriseTask"
		result.ProxyFields = solverapi.ProxyFieldsFrom(task.Proxy)
	}

	return result
}

func mapFunCaptcha(task *tasks.FunCaptchaTask) funCaptchaTask {
	result := funCaptchaTask{
		Type:                     "FunCaptchaTaskProxyless",
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
		Type:       "TurnstileTaskProxyless",
		WebsiteURL: task.WebsiteURL,
		WebsiteKey: task.WebsiteKey,
		Action:     task.Action,
		Data:       task.CData,
		PageData:   task.PageData,
	}

	if task.Proxy.IsSet() {
		result.Type = "TurnstileTask"
		result.ProxyFields = solverapi.ProxyFieldsFrom(task.Proxy)
	}

	return result
}

func mapDataDome(task *tasks.DataDomeTask) dataDomeTask {
	return dataDomeTask{
		Type:        "DataDomeSliderTask",
		WebsiteURL:  task.WebsiteURL,
		CaptchaURL:  task.CaptchaURL,
		UserAgent:   task.UserAgent,
		ProxyFields: solverapi.ProxyFieldsFrom(task.Proxy),
	}
}

func mapGeeTest(task *tasks.GeeTestTask) geeTestTask {
	result := geeTestTask{
		Type:                      "GeeTestTaskProxyless",
		WebsiteURL:                task.WebsiteURL,
		GT:                        task.GT,
		Challenge:                 task.Challenge,
		Version:                   3,
		GeetestAPIServerSubdomain: task.APIServerSubdomain,
	}

	if task.Proxy.IsSet() {
		result.Type = "GeeTestTask"
		result.ProxyFields = solverapi.ProxyFieldsFrom(task.Proxy)
	}

	return result
}

func mapGeeTestV4(task *tasks.GeeTestV4Task) geeTestTask {
	result := geeTestTask{
		Type:                      "GeeTestTaskProxyless",
		WebsiteURL:                task.WebsiteURL,
		Version:                   4,
		InitParameters:            &geetestInit{CaptchaID: task.CaptchaID},
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
		result.ProxyFields = solverapi.ProxyFieldsFrom(task.Proxy)
	}

	return result
}

func mapImageToText(task *tasks.ImageToTextTask) imageToTextTask {
	return imageToTextTask{
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

func mapAWSWAF(task *tasks.AWSWAFTask) amazonTask {
	result := amazonTask{
		Type:            "AmazonTaskProxyless",
		WebsiteURL:      task.WebsiteURL,
		WebsiteKey:      task.Key,
		IV:              task.IV,
		Context:         task.Context,
		ChallengeScript: task.ChallengeScript,
		CaptchaScript:   task.CaptchaScript,
	}

	if task.Proxy.IsSet() {
		result.Type = "AmazonTask"
		result.ProxyFields = solverapi.ProxyFieldsFrom(task.Proxy)
	}

	return result
}

func mapMTCaptcha(task *tasks.MTCaptchaTask) mtCaptchaTask {
	result := mtCaptchaTask{
		Type:       "MtCaptchaTaskProxyless",
		WebsiteURL: task.WebsiteURL,
		WebsiteKey: task.WebsiteKey,
	}

	if task.Proxy.IsSet() {
		result.Type = "MtCaptchaTask"
		result.ProxyFields = solverapi.ProxyFieldsFrom(task.Proxy)
	}

	return result
}

func mapFriendlyCaptcha(task *tasks.FriendlyCaptchaTask) friendlyCaptchaTask {
	result := friendlyCaptchaTask{
		Type:       "FriendlyCaptchaTaskProxyless",
		WebsiteURL: task.WebsiteURL,
		WebsiteKey: task.WebsiteKey,
	}

	if task.Proxy.IsSet() {
		result.Type = "FriendlyCaptchaTask"
		result.ProxyFields = solverapi.ProxyFieldsFrom(task.Proxy)
	}

	return result
}

func mapLemin(task *tasks.LeminTask) leminTask {
	result := leminTask{
		Type:                    "LeminTaskProxyless",
		WebsiteURL:              task.WebsiteURL,
		CaptchaID:               task.CaptchaID,
		DivID:                   task.DivID,
		LeminAPIServerSubdomain: task.APIServerSubdomain,
		UserAgent:               task.UserAgent,
	}

	if task.Proxy.IsSet() {
		result.Type = "LeminTask"
		result.ProxyFields = solverapi.ProxyFieldsFrom(task.Proxy)
	}

	return result
}

func mapCutCaptcha(task *tasks.CutCaptchaTask) cutCaptchaTask {
	result := cutCaptchaTask{
		Type:       "CutCaptchaTaskProxyless",
		WebsiteURL: task.WebsiteURL,
		MiseryKey:  task.MiseryKey,
		APIKey:     task.APIKey,
	}

	if task.Proxy.IsSet() {
		result.Type = "CutCaptchaTask"
		result.ProxyFields = solverapi.ProxyFieldsFrom(task.Proxy)
	}

	return result
}

func mapText(task *tasks.TextCaptchaTask) textCaptchaTask {
	return textCaptchaTask{
		Type:    "TextCaptchaTask",
		Comment: task.Question,
	}
}

func mapProsopo(task *tasks.ProsopoTask) prosopoTask {
	result := prosopoTask{
		Type:       "ProsopoTaskProxyless",
		WebsiteURL: task.WebsiteURL,
		WebsiteKey: task.WebsiteKey,
	}

	if task.Proxy.IsSet() {
		result.Type = "ProsopoTask"
		result.ProxyFields = solverapi.ProxyFieldsFrom(task.Proxy)
	}

	return result
}

func mapAltcha(task *tasks.AltchaTask) altchaTask {
	result := altchaTask{
		Type:          "AltchaTaskProxyless",
		WebsiteURL:    task.WebsiteURL,
		ChallengeURL:  task.ChallengeURL,
		ChallengeJSON: task.ChallengeJSON,
	}

	if task.Proxy.IsSet() {
		result.Type = "AltchaTask"
		result.ProxyFields = solverapi.ProxyFieldsFrom(task.Proxy)
	}

	return result
}
