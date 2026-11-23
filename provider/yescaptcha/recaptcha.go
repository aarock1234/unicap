package yescaptcha

import (
	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

// reCaptchaV2Task is YesCaptcha's reCAPTCHA v2 task, standard or Enterprise.
// The user agent and the proxy fields belong to the proxy wire types only.
type reCaptchaV2Task struct {
	Type              string         `json:"type"`
	WebsiteURL        string         `json:"websiteURL"`
	WebsiteKey        string         `json:"websiteKey"`
	IsInvisible       bool           `json:"isInvisible,omitempty"`
	EnterprisePayload map[string]any `json:"enterprisePayload,omitempty"`
	UserAgent         string         `json:"userAgent,omitempty"`
	solverapi.ProxyFields
}

// reCaptchaV3Task is YesCaptcha's reCAPTCHA v3 task, standard or Enterprise.
// It has no score field: a score is requested by choosing a wire type.
type reCaptchaV3Task struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	WebsiteKey string `json:"websiteKey"`
	PageAction string `json:"pageAction,omitempty"`
}

// reCaptchaClassificationTask is YesCaptcha's reCAPTCHA v2 image recognition
// task. The server infers the grid from the image size.
type reCaptchaClassificationTask struct {
	Type     string `json:"type"`
	Image    string `json:"image"`
	Question string `json:"question"`
}

// reCaptchaV3Types are YesCaptcha's reCAPTCHA v3 wire types, keyed by the
// minimum score each one aims for. Zero is the type without a fixed score.
var reCaptchaV3Types = map[float64]string{
	0:   "RecaptchaV3TaskProxyless",
	0.7: "RecaptchaV3TaskProxylessM1S7",
	0.9: "RecaptchaV3TaskProxylessM1S9",
}

// mapReCaptchaV2 builds YesCaptcha's reCAPTCHA v2 task.
func mapReCaptchaV2(task *tasks.ReCaptchaV2Task) (reCaptchaV2Task, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("PageAction", task.PageAction != ""),
		solverapi.Setting("DataS", task.DataS != ""),
		solverapi.Setting("APIDomain", task.APIDomain != ""),
		solverapi.Setting("IsSession", task.IsSession),
		solverapi.Setting("Cookies", task.Cookies != ""),
		solverapi.Setting("UserAgent without Proxy", task.Proxy == nil && task.UserAgent != ""),
	); err != nil {
		return reCaptchaV2Task{}, err
	}

	taskType, err := variant("NoCaptchaTask").For(task.Proxy)
	if err != nil {
		return reCaptchaV2Task{}, err
	}

	return reCaptchaV2Task{
		Type:        taskType,
		WebsiteURL:  task.WebsiteURL,
		WebsiteKey:  task.WebsiteKey,
		IsInvisible: task.IsInvisible,
		UserAgent:   task.UserAgent,
		ProxyFields: solverapi.NewProxyFields(task.Proxy),
	}, nil
}

// mapReCaptchaV2Enterprise builds YesCaptcha's reCAPTCHA v2 Enterprise task.
func mapReCaptchaV2Enterprise(task *tasks.ReCaptchaV2EnterpriseTask) (reCaptchaV2Task, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("PageAction", task.PageAction != ""),
		solverapi.Setting("APIDomain", task.APIDomain != ""),
		solverapi.Setting("IsSession", task.IsSession),
		solverapi.Setting("Cookies", task.Cookies != ""),
		solverapi.Setting("UserAgent without Proxy", task.Proxy == nil && task.UserAgent != ""),
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
		ProxyFields:       solverapi.NewProxyFields(task.Proxy),
	}, nil
}

// mapReCaptchaV3 builds YesCaptcha's reCAPTCHA v3 task. MinScore selects the
// wire type, so only the scores in reCaptchaV3Types can be honoured.
func mapReCaptchaV3(task *tasks.ReCaptchaV3Task) (reCaptchaV3Task, error) {
	scoreType, hasScoreType := reCaptchaV3Types[task.MinScore]

	if err := solverapi.Unsupported(
		solverapi.Setting("MinScore without a matching task type", !hasScoreType),
		solverapi.Setting("HighScore", task.HighScore),
		solverapi.Setting("APIDomain", task.APIDomain != ""),
		solverapi.Setting("IsSession", task.IsSession),
	); err != nil {
		return reCaptchaV3Task{}, err
	}

	taskType, err := proxyless(scoreType).For(task.Proxy)
	if err != nil {
		return reCaptchaV3Task{}, err
	}

	return reCaptchaV3Task{
		Type:       taskType,
		WebsiteURL: task.WebsiteURL,
		WebsiteKey: task.WebsiteKey,
		PageAction: task.PageAction,
	}, nil
}

// mapReCaptchaV3Enterprise builds YesCaptcha's reCAPTCHA v3 Enterprise task.
// Its wire type is proxyless despite having no Proxyless suffix.
func mapReCaptchaV3Enterprise(task *tasks.ReCaptchaV3EnterpriseTask) (reCaptchaV3Task, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("MinScore", task.MinScore != 0),
		solverapi.Setting("HighScore", task.HighScore),
		solverapi.Setting("EnterprisePayload", len(task.EnterprisePayload) > 0),
		solverapi.Setting("APIDomain", task.APIDomain != ""),
		solverapi.Setting("IsSession", task.IsSession),
	); err != nil {
		return reCaptchaV3Task{}, err
	}

	taskType, err := proxyless("RecaptchaV3EnterpriseTask").For(task.Proxy)
	if err != nil {
		return reCaptchaV3Task{}, err
	}

	return reCaptchaV3Task{
		Type:       taskType,
		WebsiteURL: task.WebsiteURL,
		WebsiteKey: task.WebsiteKey,
		PageAction: task.PageAction,
	}, nil
}

// mapReCaptchaClassification builds YesCaptcha's reCAPTCHA v2 image
// recognition task, which is solved synchronously.
func mapReCaptchaClassification(task *tasks.ReCaptchaClassificationTask) (reCaptchaClassificationTask, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("WebsiteURL", task.WebsiteURL != ""),
		solverapi.Setting("WebsiteKey", task.WebsiteKey != ""),
	); err != nil {
		return reCaptchaClassificationTask{}, err
	}

	return reCaptchaClassificationTask{
		Type:     "ReCaptchaV2Classification",
		Image:    task.Image,
		Question: task.Question,
	}, nil
}
