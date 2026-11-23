package captchaai

import (
	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

// mapTurnstile builds the form of CaptchaAI's Turnstile task, which has no
// field for the widget's cdata or a challenge page's data.
func mapTurnstile(task *tasks.TurnstileTask) (form, error) {
	proxy, proxyErr := proxyFields(task.Proxy)

	if err := solverapi.Check(
		solverapi.Unsupported(
			solverapi.Setting("CData", task.CData != ""),
			solverapi.Setting("PageData", task.PageData != ""),
		),
		proxyErr,
	); err != nil {
		return nil, err
	}

	return newForm(methodTurnstile).
		with("sitekey", task.WebsiteKey).
		with("pageurl", task.WebsiteURL).
		optional("action", task.Action).
		extend(proxy), nil
}

// mapCloudflareChallenge builds the form of CaptchaAI's Cloudflare Challenge
// task, which is solved only through the caller's proxy. CaptchaAI documents
// no field for the challenge page's HTML or for a User-Agent: it reports the
// User-Agent its solver used instead.
func mapCloudflareChallenge(task *tasks.CloudflareChallengeTask) (form, error) {
	proxy, proxyErr := proxyFields(task.Proxy)

	if err := solverapi.Check(
		solverapi.Required(
			solverapi.Setting("Proxy", task.Proxy != nil),
		),
		solverapi.Unsupported(
			solverapi.Setting("HTML", task.HTML != ""),
			solverapi.Setting("UserAgent", task.UserAgent != ""),
		),
		proxyErr,
	); err != nil {
		return nil, err
	}

	return newForm(methodCloudflareChallenge).
		with("pageurl", task.WebsiteURL).
		extend(proxy), nil
}

// mapGeeTest builds the form of CaptchaAI's GeeTest v3 task. The User-Agent
// is rejected because CaptchaAI's documentation does not settle whether this
// task takes one.
func mapGeeTest(task *tasks.GeeTestTask) (form, error) {
	proxy, proxyErr := proxyFields(task.Proxy)

	if err := solverapi.Check(
		solverapi.Unsupported(
			solverapi.Setting("UserAgent", task.UserAgent != ""),
		),
		proxyErr,
	); err != nil {
		return nil, err
	}

	return newForm(methodGeeTest).
		with("gt", task.GT).
		with("challenge", task.Challenge).
		with("pageurl", task.WebsiteURL).
		optional("api_server", task.APIServerSubdomain).
		extend(proxy), nil
}

// mapFriendlyCaptcha builds the form of CaptchaAI's Friendly Captcha task,
// which documents no field for the widget's version or its scripts.
func mapFriendlyCaptcha(task *tasks.FriendlyCaptchaTask) (form, error) {
	proxy, proxyErr := proxyFields(task.Proxy)

	if err := solverapi.Check(
		solverapi.Unsupported(
			solverapi.Setting("Version", task.Version != ""),
			solverapi.Setting("ModuleScript", task.ModuleScript != ""),
			solverapi.Setting("NomoduleScript", task.NomoduleScript != ""),
		),
		proxyErr,
	); err != nil {
		return nil, err
	}

	return newForm(methodFriendlyCaptcha).
		with("sitekey", task.WebsiteKey).
		with("pageurl", task.WebsiteURL).
		extend(proxy), nil
}
