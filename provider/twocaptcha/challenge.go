package twocaptcha

import (
	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

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

type dataDomeTask struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	CaptchaURL string `json:"captchaUrl"`
	UserAgent  string `json:"userAgent"`
	solverapi.ProxyFields
}

type geeTestTask struct {
	Type                      string         `json:"type"`
	WebsiteURL                string         `json:"websiteURL"`
	Version                   int            `json:"version"`
	GT                        string         `json:"gt,omitempty"`
	Challenge                 string         `json:"challenge,omitempty"`
	InitParameters            *geeTestParams `json:"initParameters,omitempty"`
	RiskType                  string         `json:"riskType,omitempty"`
	GeetestAPIServerSubdomain string         `json:"geetestApiServerSubdomain,omitempty"`
	UserAgent                 string         `json:"userAgent,omitempty"`
	solverapi.ProxyFields
}

type geeTestParams struct {
	CaptchaID string `json:"captcha_id"`
}

type amazonTask struct {
	Type            string `json:"type"`
	WebsiteURL      string `json:"websiteURL"`
	WebsiteKey      string `json:"websiteKey,omitempty"`
	IV              string `json:"iv,omitempty"`
	Context         string `json:"context,omitempty"`
	ChallengeScript string `json:"challengeScript,omitempty"`
	CaptchaScript   string `json:"captchaScript,omitempty"`
	JSAPIScript     string `json:"jsapiScript,omitempty"`
	solverapi.ProxyFields
}

func mapFunCaptcha(task *tasks.FunCaptchaTask) (funCaptchaTask, error) {
	taskType, err := variant("FunCaptchaTask").For(task.Proxy)
	if err != nil {
		return funCaptchaTask{}, err
	}

	return funCaptchaTask{
		Type:                     taskType,
		WebsiteURL:               task.WebsiteURL,
		WebsitePublicKey:         task.WebsitePublicKey,
		FuncaptchaAPIJSSubdomain: task.APIJSSubdomain,
		Data:                     task.Data,
		UserAgent:                task.UserAgent,
		ProxyFields:              solverapi.NewProxyFields(task.Proxy),
	}, nil
}

func mapTurnstile(task *tasks.TurnstileTask) (turnstileTask, error) {
	taskType, err := variant("TurnstileTask").For(task.Proxy)
	if err != nil {
		return turnstileTask{}, err
	}

	return turnstileTask{
		Type:        taskType,
		WebsiteURL:  task.WebsiteURL,
		WebsiteKey:  task.WebsiteKey,
		Action:      task.Action,
		Data:        task.CData,
		PageData:    task.PageData,
		ProxyFields: solverapi.NewProxyFields(task.Proxy),
	}, nil
}

func mapDataDome(task *tasks.DataDomeTask) (dataDomeTask, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("HTML", task.HTML != ""),
	); err != nil {
		return dataDomeTask{}, err
	}

	taskType, err := proxyOnly("DataDomeSliderTask").For(task.Proxy)
	if err != nil {
		return dataDomeTask{}, err
	}

	return dataDomeTask{
		Type:        taskType,
		WebsiteURL:  task.WebsiteURL,
		CaptchaURL:  task.CaptchaURL,
		UserAgent:   task.UserAgent,
		ProxyFields: solverapi.NewProxyFields(task.Proxy),
	}, nil
}

func mapGeeTest(task *tasks.GeeTestTask) (geeTestTask, error) {
	taskType, err := variant("GeeTestTask").For(task.Proxy)
	if err != nil {
		return geeTestTask{}, err
	}

	return geeTestTask{
		Type:                      taskType,
		WebsiteURL:                task.WebsiteURL,
		Version:                   3,
		GT:                        task.GT,
		Challenge:                 task.Challenge,
		GeetestAPIServerSubdomain: task.APIServerSubdomain,
		UserAgent:                 task.UserAgent,
		ProxyFields:               solverapi.NewProxyFields(task.Proxy),
	}, nil
}

func mapGeeTestV4(task *tasks.GeeTestV4Task) (geeTestTask, error) {
	taskType, err := variant("GeeTestTask").For(task.Proxy)
	if err != nil {
		return geeTestTask{}, err
	}

	return geeTestTask{
		Type:       taskType,
		WebsiteURL: task.WebsiteURL,
		Version:    4,
		InitParameters: &geeTestParams{
			CaptchaID: task.CaptchaID,
		},
		RiskType:                  task.RiskType,
		GeetestAPIServerSubdomain: task.APIServerSubdomain,
		UserAgent:                 task.UserAgent,
		ProxyFields:               solverapi.NewProxyFields(task.Proxy),
	}, nil
}

func mapAWSWAF(task *tasks.AWSWAFTask) (amazonTask, error) {
	hasGokuProps := task.Key != "" && task.IV != "" && task.Context != ""

	if err := solverapi.Check(
		solverapi.Required(
			solverapi.Setting("Key, IV, and Context, or JSAPIScript", hasGokuProps || task.JSAPIScript != ""),
		),
		solverapi.Unsupported(
			solverapi.Setting("ProblemURL", task.ProblemURL != ""),
			solverapi.Setting("APIKey", task.APIKey != ""),
			solverapi.Setting("ExistingToken", task.ExistingToken != ""),
		),
	); err != nil {
		return amazonTask{}, err
	}

	taskType, err := variant("AmazonTask").For(task.Proxy)
	if err != nil {
		return amazonTask{}, err
	}

	return amazonTask{
		Type:            taskType,
		WebsiteURL:      task.WebsiteURL,
		WebsiteKey:      task.Key,
		IV:              task.IV,
		Context:         task.Context,
		ChallengeScript: task.ChallengeScript,
		CaptchaScript:   task.CaptchaScript,
		JSAPIScript:     task.JSAPIScript,
		ProxyFields:     solverapi.NewProxyFields(task.Proxy),
	}, nil
}
