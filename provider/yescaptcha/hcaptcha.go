package yescaptcha

import (
	"fmt"
	"maps"
	"slices"

	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

// hCaptchaTask is YesCaptcha's hCaptcha task. The proxy fields belong to the
// proxy wire type only.
type hCaptchaTask struct {
	Type        string `json:"type"`
	WebsiteURL  string `json:"websiteURL"`
	WebsiteKey  string `json:"websiteKey"`
	IsInvisible bool   `json:"isInvisible,omitempty"`
	RQData      string `json:"rqdata,omitempty"`
	solverapi.ProxyFields
}

// rqdataKey is the only tasks.HCaptchaTask.EnterprisePayload key YesCaptcha
// has a field for.
const rqdataKey = "rqdata"

// mapHCaptcha builds YesCaptcha's hCaptcha task. UserAgent is rejected
// because YesCaptcha accepts only the value its own /useragent endpoint
// returns.
func mapHCaptcha(task *tasks.HCaptchaTask) (hCaptchaTask, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("UserAgent", task.UserAgent != ""),
		solverapi.Setting("Cookies", task.Cookies != ""),
	); err != nil {
		return hCaptchaTask{}, err
	}

	rqdata, err := enterpriseRQData(task.EnterprisePayload)
	if err != nil {
		return hCaptchaTask{}, err
	}

	taskType, err := variant("HCaptchaTask").For(task.Proxy)
	if err != nil {
		return hCaptchaTask{}, err
	}

	return hCaptchaTask{
		Type:        taskType,
		WebsiteURL:  task.WebsiteURL,
		WebsiteKey:  task.WebsiteKey,
		IsInvisible: task.IsInvisible,
		RQData:      rqdata,
		ProxyFields: solverapi.NewProxyFields(task.Proxy),
	}, nil
}

// enterpriseRQData returns the rqdata string of an hCaptcha Enterprise
// payload. Every other key, and an rqdata that is not a string, is reported
// as unsupported, in key order.
func enterpriseRQData(payload map[string]any) (string, error) {
	var (
		rqdata      string
		unsupported []solverapi.Field
	)

	for _, key := range slices.Sorted(maps.Keys(payload)) {
		value, isString := payload[key].(string)
		if key == rqdataKey && isString {
			rqdata = value

			continue
		}

		unsupported = append(unsupported, solverapi.Setting(fmt.Sprintf("EnterprisePayload[%s]", key), true))
	}

	if err := solverapi.Unsupported(unsupported...); err != nil {
		return "", err
	}

	return rqdata, nil
}
