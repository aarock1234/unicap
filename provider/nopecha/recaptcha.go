package nopecha

import (
	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

// enterpriseSKey is the only EnterprisePayload key of the reCAPTCHA
// Enterprise tasks that NopeCHA has a field for.
const enterpriseSKey = "s"

// reCaptchaV2Request is the body NopeCHA's reCAPTCHA v2 token endpoint
// takes, for the standard and the Enterprise widget alike.
type reCaptchaV2Request struct {
	Sitekey    string          `json:"sitekey"`
	URL        string          `json:"url"`
	Enterprise bool            `json:"enterprise,omitempty"`
	Data       reCaptchaV2Data `json:"data,omitzero"`
	Proxy      *proxy          `json:"proxy,omitempty"`
	UserAgent  string          `json:"useragent,omitempty"`
}

// reCaptchaV2Data holds the widget parameters of a reCAPTCHA v2 job.
type reCaptchaV2Data struct {
	S string `json:"s,omitempty"`
}

// reCaptchaV3Request is the body NopeCHA's reCAPTCHA v3 token endpoint
// takes, for standard and Enterprise reCAPTCHA alike.
type reCaptchaV3Request struct {
	Sitekey    string          `json:"sitekey"`
	URL        string          `json:"url"`
	Enterprise bool            `json:"enterprise,omitempty"`
	Data       reCaptchaV3Data `json:"data,omitzero"`
	Proxy      *proxy          `json:"proxy,omitempty"`
}

// reCaptchaV3Data holds the parameters of a reCAPTCHA v3 job.
type reCaptchaV3Data struct {
	Action string `json:"action,omitempty"`
	S      string `json:"s,omitempty"`
}

// mapReCaptchaV2 builds the body of NopeCHA's reCAPTCHA v2 token endpoint.
// NopeCHA has no field that marks an invisible widget, and its cookie
// objects need attributes that the task's cookie string does not carry.
func mapReCaptchaV2(task *tasks.ReCaptchaV2Task) (reCaptchaV2Request, error) {
	if err := solverapi.Check(
		solverapi.Unsupported(
			solverapi.Setting("IsInvisible", task.IsInvisible),
			solverapi.Setting("PageAction", task.PageAction != ""),
			solverapi.Setting("APIDomain", task.APIDomain != ""),
			solverapi.Setting("IsSession", task.IsSession),
			solverapi.Setting("Cookies", task.Cookies != ""),
		),
		checkProxy(task.Proxy),
	); err != nil {
		return reCaptchaV2Request{}, err
	}

	return reCaptchaV2Request{
		Sitekey: task.WebsiteKey,
		URL:     task.WebsiteURL,
		Data: reCaptchaV2Data{
			S: task.DataS,
		},
		Proxy:     newProxy(task.Proxy),
		UserAgent: task.UserAgent,
	}, nil
}

// mapReCaptchaV2Enterprise builds the body of NopeCHA's reCAPTCHA v2 token
// endpoint for an Enterprise widget. It rejects what mapReCaptchaV2 rejects.
func mapReCaptchaV2Enterprise(task *tasks.ReCaptchaV2EnterpriseTask) (reCaptchaV2Request, error) {
	if err := solverapi.Check(
		solverapi.Unsupported(
			solverapi.Setting("IsInvisible", task.IsInvisible),
			solverapi.Setting("PageAction", task.PageAction != ""),
			solverapi.Setting("APIDomain", task.APIDomain != ""),
			solverapi.Setting("IsSession", task.IsSession),
			solverapi.Setting("Cookies", task.Cookies != ""),
		),
		checkProxy(task.Proxy),
	); err != nil {
		return reCaptchaV2Request{}, err
	}

	dataS, err := enterpriseValue(task.EnterprisePayload, enterpriseSKey)
	if err != nil {
		return reCaptchaV2Request{}, err
	}

	return reCaptchaV2Request{
		Sitekey:    task.WebsiteKey,
		URL:        task.WebsiteURL,
		Enterprise: true,
		Data: reCaptchaV2Data{
			S: dataS,
		},
		Proxy:     newProxy(task.Proxy),
		UserAgent: task.UserAgent,
	}, nil
}

// mapReCaptchaV3 builds the body of NopeCHA's reCAPTCHA v3 token endpoint.
// NopeCHA has no field for a minimum score.
func mapReCaptchaV3(task *tasks.ReCaptchaV3Task) (reCaptchaV3Request, error) {
	if err := solverapi.Check(
		solverapi.Unsupported(
			solverapi.Setting("MinScore", task.MinScore != 0),
			solverapi.Setting("HighScore", task.HighScore),
			solverapi.Setting("APIDomain", task.APIDomain != ""),
			solverapi.Setting("IsSession", task.IsSession),
		),
		checkProxy(task.Proxy),
	); err != nil {
		return reCaptchaV3Request{}, err
	}

	return reCaptchaV3Request{
		Sitekey: task.WebsiteKey,
		URL:     task.WebsiteURL,
		Data: reCaptchaV3Data{
			Action: task.PageAction,
		},
		Proxy: newProxy(task.Proxy),
	}, nil
}

// mapReCaptchaV3Enterprise builds the body of NopeCHA's reCAPTCHA v3 token
// endpoint for reCAPTCHA Enterprise. It rejects what mapReCaptchaV3 rejects.
func mapReCaptchaV3Enterprise(task *tasks.ReCaptchaV3EnterpriseTask) (reCaptchaV3Request, error) {
	if err := solverapi.Check(
		solverapi.Unsupported(
			solverapi.Setting("MinScore", task.MinScore != 0),
			solverapi.Setting("HighScore", task.HighScore),
			solverapi.Setting("APIDomain", task.APIDomain != ""),
			solverapi.Setting("IsSession", task.IsSession),
		),
		checkProxy(task.Proxy),
	); err != nil {
		return reCaptchaV3Request{}, err
	}

	dataS, err := enterpriseValue(task.EnterprisePayload, enterpriseSKey)
	if err != nil {
		return reCaptchaV3Request{}, err
	}

	return reCaptchaV3Request{
		Sitekey:    task.WebsiteKey,
		URL:        task.WebsiteURL,
		Enterprise: true,
		Data: reCaptchaV3Data{
			Action: task.PageAction,
			S:      dataS,
		},
		Proxy: newProxy(task.Proxy),
	}, nil
}
