package capmonster

import (
	"fmt"
	"maps"
	"slices"

	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

// reCaptchaV2Task is CapMonster Cloud's reCAPTCHA v2 task. One wire type
// covers both proxy modes.
type reCaptchaV2Task struct {
	Type                string `json:"type"`
	WebsiteURL          string `json:"websiteURL"`
	WebsiteKey          string `json:"websiteKey"`
	RecaptchaDataSValue string `json:"recaptchaDataSValue,omitempty"`
	UserAgent           string `json:"userAgent,omitempty"`
	Cookies             string `json:"cookies,omitempty"`
	IsInvisible         bool   `json:"isInvisible,omitempty"`
	solverapi.ProxyFields
}

// reCaptchaV2EnterpriseTask is CapMonster Cloud's reCAPTCHA v2 Enterprise
// task. One wire type covers both proxy modes.
type reCaptchaV2EnterpriseTask struct {
	Type              string             `json:"type"`
	WebsiteURL        string             `json:"websiteURL"`
	WebsiteKey        string             `json:"websiteKey"`
	PageAction        string             `json:"pageAction,omitempty"`
	EnterprisePayload *enterprisePayload `json:"enterprisePayload,omitempty"`
	APIDomain         string             `json:"apiDomain,omitempty"`
	UserAgent         string             `json:"userAgent,omitempty"`
	Cookies           string             `json:"cookies,omitempty"`
	solverapi.ProxyFields
}

// enterprisePayload is CapMonster Cloud's enterprisePayload object, which
// documents a single key.
type enterprisePayload struct {
	S string `json:"s"`
}

// reCaptchaV3Task is CapMonster Cloud's reCAPTCHA v3 task, standard or
// Enterprise. Both are solved through CapMonster Cloud's own proxies.
type reCaptchaV3Task struct {
	Type       string  `json:"type"`
	WebsiteURL string  `json:"websiteURL"`
	WebsiteKey string  `json:"websiteKey"`
	MinScore   float64 `json:"minScore,omitempty"`
	PageAction string  `json:"pageAction,omitempty"`
}

// enterpriseSKey is the only tasks.ReCaptchaV2EnterpriseTask.EnterprisePayload
// key CapMonster Cloud has a field for.
const enterpriseSKey = "s"

// mapReCaptchaV2 builds CapMonster Cloud's reCAPTCHA v2 task.
func mapReCaptchaV2(task *tasks.ReCaptchaV2Task) (reCaptchaV2Task, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("PageAction", task.PageAction != ""),
		solverapi.Setting("APIDomain", task.APIDomain != ""),
		solverapi.Setting("IsSession", task.IsSession),
	); err != nil {
		return reCaptchaV2Task{}, err
	}

	proxy, err := proxyFields(task.Proxy)
	if err != nil {
		return reCaptchaV2Task{}, err
	}

	return reCaptchaV2Task{
		Type:                "RecaptchaV2Task",
		WebsiteURL:          task.WebsiteURL,
		WebsiteKey:          task.WebsiteKey,
		RecaptchaDataSValue: task.DataS,
		UserAgent:           task.UserAgent,
		Cookies:             task.Cookies,
		IsInvisible:         task.IsInvisible,
		ProxyFields:         proxy,
	}, nil
}

// mapReCaptchaV2Enterprise builds CapMonster Cloud's reCAPTCHA v2 Enterprise
// task.
func mapReCaptchaV2Enterprise(task *tasks.ReCaptchaV2EnterpriseTask) (reCaptchaV2EnterpriseTask, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("IsInvisible", task.IsInvisible),
		solverapi.Setting("IsSession", task.IsSession),
	); err != nil {
		return reCaptchaV2EnterpriseTask{}, err
	}

	payload, err := newEnterprisePayload(task.EnterprisePayload)
	if err != nil {
		return reCaptchaV2EnterpriseTask{}, err
	}

	proxy, err := proxyFields(task.Proxy)
	if err != nil {
		return reCaptchaV2EnterpriseTask{}, err
	}

	return reCaptchaV2EnterpriseTask{
		Type:              "RecaptchaV2EnterpriseTask",
		WebsiteURL:        task.WebsiteURL,
		WebsiteKey:        task.WebsiteKey,
		PageAction:        task.PageAction,
		EnterprisePayload: payload,
		APIDomain:         task.APIDomain,
		UserAgent:         task.UserAgent,
		Cookies:           task.Cookies,
		ProxyFields:       proxy,
	}, nil
}

// mapReCaptchaV3 builds CapMonster Cloud's reCAPTCHA v3 task. CapMonster
// Cloud documents minimum scores from 0.1 to 0.9.
func mapReCaptchaV3(task *tasks.ReCaptchaV3Task) (reCaptchaV3Task, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("HighScore", task.HighScore),
		solverapi.Setting("APIDomain", task.APIDomain != ""),
		solverapi.Setting("IsSession", task.IsSession),
		solverapi.Setting("Proxy", task.Proxy != nil),
	); err != nil {
		return reCaptchaV3Task{}, err
	}

	return reCaptchaV3Task{
		Type:       "RecaptchaV3TaskProxyless",
		WebsiteURL: task.WebsiteURL,
		WebsiteKey: task.WebsiteKey,
		MinScore:   task.MinScore,
		PageAction: task.PageAction,
	}, nil
}

// mapReCaptchaV3Enterprise builds CapMonster Cloud's reCAPTCHA v3 Enterprise
// task. CapMonster Cloud documents minimum scores from 0.1 to 0.9.
func mapReCaptchaV3Enterprise(task *tasks.ReCaptchaV3EnterpriseTask) (reCaptchaV3Task, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("HighScore", task.HighScore),
		solverapi.Setting("EnterprisePayload", len(task.EnterprisePayload) > 0),
		solverapi.Setting("APIDomain", task.APIDomain != ""),
		solverapi.Setting("IsSession", task.IsSession),
		solverapi.Setting("Proxy", task.Proxy != nil),
	); err != nil {
		return reCaptchaV3Task{}, err
	}

	return reCaptchaV3Task{
		Type:       "RecaptchaV3EnterpriseTask",
		WebsiteURL: task.WebsiteURL,
		WebsiteKey: task.WebsiteKey,
		MinScore:   task.MinScore,
		PageAction: task.PageAction,
	}, nil
}

// newEnterprisePayload returns the enterprisePayload object for a reCAPTCHA
// Enterprise payload, or nil when the payload has no "s" value. Every other
// key, and an "s" that is not a string, is reported as unsupported, in key
// order.
func newEnterprisePayload(payload map[string]any) (*enterprisePayload, error) {
	var (
		wire        *enterprisePayload
		unsupported []solverapi.Field
	)

	for _, key := range slices.Sorted(maps.Keys(payload)) {
		value, isString := payload[key].(string)
		if key == enterpriseSKey && isString {
			wire = &enterprisePayload{
				S: value,
			}

			continue
		}

		unsupported = append(unsupported, solverapi.Setting(fmt.Sprintf("EnterprisePayload[%s]", key), true))
	}

	if err := solverapi.Unsupported(unsupported...); err != nil {
		return nil, err
	}

	return wire, nil
}
