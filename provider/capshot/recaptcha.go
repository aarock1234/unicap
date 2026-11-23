package capshot

import (
	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

// reCaptchaV2Task is Capshot's reCAPTCHA v2 task.
type reCaptchaV2Task struct {
	Type                string `json:"type"`
	WebsiteURL          string `json:"websiteURL"`
	WebsiteKey          string `json:"websiteKey"`
	IsInvisible         bool   `json:"isInvisible,omitempty"`
	RecaptchaDataSValue string `json:"recaptchaDataSValue,omitempty"`
}

// reCaptchaV2EnterpriseTask is Capshot's reCAPTCHA v2 Enterprise task.
type reCaptchaV2EnterpriseTask struct {
	Type              string         `json:"type"`
	WebsiteURL        string         `json:"websiteURL"`
	WebsiteKey        string         `json:"websiteKey"`
	EnterprisePayload map[string]any `json:"enterprisePayload,omitempty"`
}

// reCaptchaV3Task is Capshot's reCAPTCHA v3 task, standard or Enterprise.
// Capshot requires the page action and the minimum score. The enterprise
// payload belongs to the Enterprise wire type only.
type reCaptchaV3Task struct {
	Type              string         `json:"type"`
	WebsiteURL        string         `json:"websiteURL"`
	WebsiteKey        string         `json:"websiteKey"`
	PageAction        string         `json:"pageAction"`
	MinScore          float64        `json:"minScore"`
	EnterprisePayload map[string]any `json:"enterprisePayload,omitempty"`
}

// mapReCaptchaV2 builds Capshot's reCAPTCHA v2 task.
func mapReCaptchaV2(task *tasks.ReCaptchaV2Task) (reCaptchaV2Task, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("PageAction", task.PageAction != ""),
		solverapi.Setting("APIDomain", task.APIDomain != ""),
		solverapi.Setting("IsSession", task.IsSession),
		solverapi.Setting("UserAgent", task.UserAgent != ""),
		solverapi.Setting("Cookies", task.Cookies != ""),
		solverapi.Setting("Proxy", task.Proxy != nil),
	); err != nil {
		return reCaptchaV2Task{}, err
	}

	return reCaptchaV2Task{
		Type:                "RecaptchaV2TaskProxyless",
		WebsiteURL:          task.WebsiteURL,
		WebsiteKey:          task.WebsiteKey,
		IsInvisible:         task.IsInvisible,
		RecaptchaDataSValue: task.DataS,
	}, nil
}

// mapReCaptchaV2Enterprise builds Capshot's reCAPTCHA v2 Enterprise task.
func mapReCaptchaV2Enterprise(task *tasks.ReCaptchaV2EnterpriseTask) (reCaptchaV2EnterpriseTask, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("IsInvisible", task.IsInvisible),
		solverapi.Setting("PageAction", task.PageAction != ""),
		solverapi.Setting("APIDomain", task.APIDomain != ""),
		solverapi.Setting("IsSession", task.IsSession),
		solverapi.Setting("UserAgent", task.UserAgent != ""),
		solverapi.Setting("Cookies", task.Cookies != ""),
		solverapi.Setting("Proxy", task.Proxy != nil),
	); err != nil {
		return reCaptchaV2EnterpriseTask{}, err
	}

	return reCaptchaV2EnterpriseTask{
		Type:              "RecaptchaV2EnterpriseTaskProxyless",
		WebsiteURL:        task.WebsiteURL,
		WebsiteKey:        task.WebsiteKey,
		EnterprisePayload: task.EnterprisePayload,
	}, nil
}

// mapReCaptchaV3 builds Capshot's reCAPTCHA v3 task. A zero MinScore cannot
// be told from an unset one, so it counts as missing.
func mapReCaptchaV3(task *tasks.ReCaptchaV3Task) (reCaptchaV3Task, error) {
	if err := solverapi.Check(
		solverapi.Required(
			solverapi.Setting("PageAction", task.PageAction != ""),
			solverapi.Setting("MinScore", task.MinScore != 0),
		),
		solverapi.Unsupported(
			solverapi.Setting("HighScore", task.HighScore),
			solverapi.Setting("APIDomain", task.APIDomain != ""),
			solverapi.Setting("IsSession", task.IsSession),
			solverapi.Setting("Proxy", task.Proxy != nil),
		),
	); err != nil {
		return reCaptchaV3Task{}, err
	}

	return reCaptchaV3Task{
		Type:       "RecaptchaV3TaskProxyless",
		WebsiteURL: task.WebsiteURL,
		WebsiteKey: task.WebsiteKey,
		PageAction: task.PageAction,
		MinScore:   task.MinScore,
	}, nil
}

// mapReCaptchaV3Enterprise builds Capshot's reCAPTCHA v3 Enterprise task. A
// zero MinScore cannot be told from an unset one, so it counts as missing.
func mapReCaptchaV3Enterprise(task *tasks.ReCaptchaV3EnterpriseTask) (reCaptchaV3Task, error) {
	if err := solverapi.Check(
		solverapi.Required(
			solverapi.Setting("PageAction", task.PageAction != ""),
			solverapi.Setting("MinScore", task.MinScore != 0),
		),
		solverapi.Unsupported(
			solverapi.Setting("HighScore", task.HighScore),
			solverapi.Setting("APIDomain", task.APIDomain != ""),
			solverapi.Setting("IsSession", task.IsSession),
			solverapi.Setting("Proxy", task.Proxy != nil),
		),
	); err != nil {
		return reCaptchaV3Task{}, err
	}

	return reCaptchaV3Task{
		Type:              "RecaptchaV3EnterpriseTaskProxyless",
		WebsiteURL:        task.WebsiteURL,
		WebsiteKey:        task.WebsiteKey,
		PageAction:        task.PageAction,
		MinScore:          task.MinScore,
		EnterprisePayload: task.EnterprisePayload,
	}, nil
}
