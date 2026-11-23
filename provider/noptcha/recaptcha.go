package noptcha

import (
	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

// reCaptchaV2Params are the params of Noptcha's recaptcha.v2 and
// recaptcha.v2.enterprise tasks. The enterprise payload belongs to the
// Enterprise type only and is passed to the widget as is.
type reCaptchaV2Params struct {
	WebsiteURL        string         `json:"website_url"`
	WebsiteKey        string         `json:"website_key"`
	IsInvisible       bool           `json:"is_invisible,omitempty"`
	EnterprisePayload map[string]any `json:"enterprise_payload,omitempty"`
	APIDomain         string         `json:"api_domain,omitempty"`
	UserAgent         string         `json:"user_agent,omitempty"`
	Cookies           string         `json:"cookies,omitempty"`
	Proxy             *proxy         `json:"proxy,omitempty"`
}

// reCaptchaV3Params are the params of Noptcha's recaptcha.v3 and
// recaptcha.v3.enterprise tasks, which take no proxy. The enterprise payload
// belongs to the Enterprise type only and is passed on as is.
type reCaptchaV3Params struct {
	WebsiteURL        string         `json:"website_url"`
	WebsiteKey        string         `json:"website_key"`
	PageAction        string         `json:"page_action,omitempty"`
	MinScore          float64        `json:"min_score,omitempty"`
	EnterprisePayload map[string]any `json:"enterprise_payload,omitempty"`
	APIDomain         string         `json:"api_domain,omitempty"`
}

// reCaptchaGridParams are the params of Noptcha's image.grid.recaptcha task.
// The server infers the grid from the image size.
type reCaptchaGridParams struct {
	Image    string `json:"image"`
	Question string `json:"question"`
}

// The minimum scores Noptcha accepts for reCAPTCHA v3.
const (
	lowestMinScore  = 0.1
	highestMinScore = 0.9
)

// mapReCaptchaV2 builds the params of Noptcha's recaptcha.v2 task.
func mapReCaptchaV2(task *tasks.ReCaptchaV2Task) (reCaptchaV2Params, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("PageAction", task.PageAction != ""),
		solverapi.Setting("DataS", task.DataS != ""),
		solverapi.Setting("IsSession", task.IsSession),
	); err != nil {
		return reCaptchaV2Params{}, err
	}

	return reCaptchaV2Params{
		WebsiteURL:  task.WebsiteURL,
		WebsiteKey:  task.WebsiteKey,
		IsInvisible: task.IsInvisible,
		APIDomain:   task.APIDomain,
		UserAgent:   task.UserAgent,
		Cookies:     task.Cookies,
		Proxy:       newProxy(task.Proxy),
	}, nil
}

// mapReCaptchaV2Enterprise builds the params of Noptcha's
// recaptcha.v2.enterprise task.
func mapReCaptchaV2Enterprise(task *tasks.ReCaptchaV2EnterpriseTask) (reCaptchaV2Params, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("PageAction", task.PageAction != ""),
		solverapi.Setting("IsSession", task.IsSession),
	); err != nil {
		return reCaptchaV2Params{}, err
	}

	return reCaptchaV2Params{
		WebsiteURL:        task.WebsiteURL,
		WebsiteKey:        task.WebsiteKey,
		IsInvisible:       task.IsInvisible,
		EnterprisePayload: task.EnterprisePayload,
		APIDomain:         task.APIDomain,
		UserAgent:         task.UserAgent,
		Cookies:           task.Cookies,
		Proxy:             newProxy(task.Proxy),
	}, nil
}

// mapReCaptchaV3 builds the params of Noptcha's recaptcha.v3 task. A zero
// MinScore is left out, which asks for no particular score.
func mapReCaptchaV3(task *tasks.ReCaptchaV3Task) (reCaptchaV3Params, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("MinScore outside 0.1 to 0.9", !acceptsMinScore(task.MinScore)),
		solverapi.Setting("HighScore", task.HighScore),
		solverapi.Setting("IsSession", task.IsSession),
		solverapi.Setting("Proxy", task.Proxy != nil),
	); err != nil {
		return reCaptchaV3Params{}, err
	}

	return reCaptchaV3Params{
		WebsiteURL: task.WebsiteURL,
		WebsiteKey: task.WebsiteKey,
		PageAction: task.PageAction,
		MinScore:   task.MinScore,
		APIDomain:  task.APIDomain,
	}, nil
}

// mapReCaptchaV3Enterprise builds the params of Noptcha's
// recaptcha.v3.enterprise task. A zero MinScore is left out, which asks for
// no particular score.
func mapReCaptchaV3Enterprise(task *tasks.ReCaptchaV3EnterpriseTask) (reCaptchaV3Params, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("MinScore outside 0.1 to 0.9", !acceptsMinScore(task.MinScore)),
		solverapi.Setting("HighScore", task.HighScore),
		solverapi.Setting("IsSession", task.IsSession),
		solverapi.Setting("Proxy", task.Proxy != nil),
	); err != nil {
		return reCaptchaV3Params{}, err
	}

	return reCaptchaV3Params{
		WebsiteURL:        task.WebsiteURL,
		WebsiteKey:        task.WebsiteKey,
		PageAction:        task.PageAction,
		MinScore:          task.MinScore,
		EnterprisePayload: task.EnterprisePayload,
		APIDomain:         task.APIDomain,
	}, nil
}

// mapReCaptchaClassification builds the params of Noptcha's
// image.grid.recaptcha task.
func mapReCaptchaClassification(task *tasks.ReCaptchaClassificationTask) (reCaptchaGridParams, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("WebsiteURL", task.WebsiteURL != ""),
		solverapi.Setting("WebsiteKey", task.WebsiteKey != ""),
	); err != nil {
		return reCaptchaGridParams{}, err
	}

	return reCaptchaGridParams{
		Image:    task.Image,
		Question: task.Question,
	}, nil
}

// acceptsMinScore reports whether Noptcha can honour a minimum score: either
// none at all, or one within the range it accepts.
func acceptsMinScore(score float64) bool {
	return score == 0 || (score >= lowestMinScore && score <= highestMinScore)
}
