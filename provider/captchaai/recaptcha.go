package captchaai

import (
	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

// versionV3 is the version flag of a reCAPTCHA v3 task. A task without the
// flag is reCAPTCHA v2.
const versionV3 = "v3"

// mapReCaptchaV2 builds the form of CaptchaAI's reCAPTCHA v2 task. CaptchaAI
// names the widget domains without the "www." unicap documents, has no
// action field for this task, and takes cookies as objects with a domain,
// which the task's cookie string does not carry.
func mapReCaptchaV2(task *tasks.ReCaptchaV2Task) (form, error) {
	proxy, proxyErr := proxyFields(task.Proxy)

	if err := solverapi.Check(
		solverapi.Unsupported(
			solverapi.Setting("PageAction", task.PageAction != ""),
			solverapi.Setting("APIDomain", task.APIDomain != ""),
			solverapi.Setting("IsSession", task.IsSession),
			solverapi.Setting("Cookies", task.Cookies != ""),
		),
		proxyErr,
	); err != nil {
		return nil, err
	}

	return newForm(methodReCaptcha).
		with("googlekey", task.WebsiteKey).
		with("pageurl", task.WebsiteURL).
		flag("invisible", task.IsInvisible).
		optional("data-s", task.DataS).
		optional("userAgent", task.UserAgent).
		extend(proxy), nil
}

// mapReCaptchaV2Enterprise builds the form of CaptchaAI's reCAPTCHA v2
// Enterprise task, which has no field for an enterprise payload or a widget
// domain. Cookies are rejected as they are for reCAPTCHA v2.
func mapReCaptchaV2Enterprise(task *tasks.ReCaptchaV2EnterpriseTask) (form, error) {
	proxy, proxyErr := proxyFields(task.Proxy)

	if err := solverapi.Check(
		solverapi.Unsupported(
			solverapi.Setting("EnterprisePayload", len(task.EnterprisePayload) > 0),
			solverapi.Setting("APIDomain", task.APIDomain != ""),
			solverapi.Setting("IsSession", task.IsSession),
			solverapi.Setting("Cookies", task.Cookies != ""),
		),
		proxyErr,
	); err != nil {
		return nil, err
	}

	return newForm(methodReCaptcha).
		with("enterprise", flagSet).
		with("googlekey", task.WebsiteKey).
		with("pageurl", task.WebsiteURL).
		optional("action", task.PageAction).
		flag("invisible", task.IsInvisible).
		optional("userAgent", task.UserAgent).
		extend(proxy), nil
}

// mapReCaptchaV3 builds the form of CaptchaAI's reCAPTCHA v3 task. CaptchaAI
// documents no field for a minimum score, and names the script domains
// without the "www." unicap documents.
func mapReCaptchaV3(task *tasks.ReCaptchaV3Task) (form, error) {
	proxy, proxyErr := proxyFields(task.Proxy)

	if err := solverapi.Check(
		solverapi.Unsupported(
			solverapi.Setting("MinScore", task.MinScore != 0),
			solverapi.Setting("HighScore", task.HighScore),
			solverapi.Setting("APIDomain", task.APIDomain != ""),
			solverapi.Setting("IsSession", task.IsSession),
		),
		proxyErr,
	); err != nil {
		return nil, err
	}

	return newForm(methodReCaptcha).
		with("version", versionV3).
		with("googlekey", task.WebsiteKey).
		with("pageurl", task.WebsiteURL).
		optional("action", task.PageAction).
		extend(proxy), nil
}

// mapReCaptchaV3Enterprise builds the form of CaptchaAI's reCAPTCHA v3
// Enterprise task. It rejects what mapReCaptchaV3 rejects, and the
// enterprise payload, which CaptchaAI has no field for.
func mapReCaptchaV3Enterprise(task *tasks.ReCaptchaV3EnterpriseTask) (form, error) {
	proxy, proxyErr := proxyFields(task.Proxy)

	if err := solverapi.Check(
		solverapi.Unsupported(
			solverapi.Setting("MinScore", task.MinScore != 0),
			solverapi.Setting("HighScore", task.HighScore),
			solverapi.Setting("EnterprisePayload", len(task.EnterprisePayload) > 0),
			solverapi.Setting("APIDomain", task.APIDomain != ""),
			solverapi.Setting("IsSession", task.IsSession),
		),
		proxyErr,
	); err != nil {
		return nil, err
	}

	return newForm(methodReCaptcha).
		with("version", versionV3).
		with("enterprise", flagSet).
		with("googlekey", task.WebsiteKey).
		with("pageurl", task.WebsiteURL).
		optional("action", task.PageAction).
		extend(proxy), nil
}
