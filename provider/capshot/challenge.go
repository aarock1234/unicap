package capshot

import (
	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

// geeTestTask is Capshot's GeeTest task. Capshot detects the GeeTest version
// from the parameters sent, and documents only the v4 captcha ID.
type geeTestTask struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	CaptchaID  string `json:"captchaId"`
}

// funCaptchaTask is Capshot's FunCaptcha task.
type funCaptchaTask struct {
	Type                     string `json:"type"`
	WebsiteURL               string `json:"websiteURL"`
	WebsitePublicKey         string `json:"websitePublicKey"`
	FuncaptchaAPIJSSubdomain string `json:"funcaptchaApiJSSubdomain,omitempty"`
}

// friendlyCaptchaTask is Capshot's Friendly Captcha task.
type friendlyCaptchaTask struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	WebsiteKey string `json:"websiteKey"`
}

// awsWAFTask is Capshot's AWS WAF task. The proxy fields belong to the proxy
// wire type only.
type awsWAFTask struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	solverapi.ProxyFields
}

// awsWAFVariant names Capshot's AWS WAF wire types, the only pair its
// documentation prints in both the proxyless and the proxy form.
var awsWAFVariant = solverapi.Variant{
	Proxyless: "AntiAwsWafTaskProxyless",
	Proxy:     "AntiAwsWafTask",
}

// mapGeeTestV4 builds Capshot's GeeTest task for a v4 challenge.
func mapGeeTestV4(task *tasks.GeeTestV4Task) (geeTestTask, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("RiskType", task.RiskType != ""),
		solverapi.Setting("APIServerSubdomain", task.APIServerSubdomain != ""),
		solverapi.Setting("UserAgent", task.UserAgent != ""),
		solverapi.Setting("Proxy", task.Proxy != nil),
	); err != nil {
		return geeTestTask{}, err
	}

	return geeTestTask{
		Type:       "GeeTestTaskProxyless",
		WebsiteURL: task.WebsiteURL,
		CaptchaID:  task.CaptchaID,
	}, nil
}

// mapFunCaptcha builds Capshot's FunCaptcha task. Data is rejected because
// Capshot's documentation does not settle how the blob is sent.
func mapFunCaptcha(task *tasks.FunCaptchaTask) (funCaptchaTask, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("Data", task.Data != ""),
		solverapi.Setting("UserAgent", task.UserAgent != ""),
		solverapi.Setting("Proxy", task.Proxy != nil),
	); err != nil {
		return funCaptchaTask{}, err
	}

	return funCaptchaTask{
		Type:                     "FunCaptchaTaskProxyless",
		WebsiteURL:               task.WebsiteURL,
		WebsitePublicKey:         task.WebsitePublicKey,
		FuncaptchaAPIJSSubdomain: task.APIJSSubdomain,
	}, nil
}

// mapFriendlyCaptcha builds Capshot's Friendly Captcha task.
func mapFriendlyCaptcha(task *tasks.FriendlyCaptchaTask) (friendlyCaptchaTask, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("Version", task.Version != ""),
		solverapi.Setting("ModuleScript", task.ModuleScript != ""),
		solverapi.Setting("NomoduleScript", task.NomoduleScript != ""),
		solverapi.Setting("Proxy", task.Proxy != nil),
	); err != nil {
		return friendlyCaptchaTask{}, err
	}

	return friendlyCaptchaTask{
		Type:       "FriendlyCaptchaTaskProxyless",
		WebsiteURL: task.WebsiteURL,
		WebsiteKey: task.WebsiteKey,
	}, nil
}

// mapAWSWAF builds Capshot's AWS WAF task, which works from the protected
// page's URL alone.
func mapAWSWAF(task *tasks.AWSWAFTask) (awsWAFTask, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("Key", task.Key != ""),
		solverapi.Setting("IV", task.IV != ""),
		solverapi.Setting("Context", task.Context != ""),
		solverapi.Setting("ChallengeScript", task.ChallengeScript != ""),
		solverapi.Setting("CaptchaScript", task.CaptchaScript != ""),
		solverapi.Setting("JSAPIScript", task.JSAPIScript != ""),
		solverapi.Setting("ProblemURL", task.ProblemURL != ""),
		solverapi.Setting("APIKey", task.APIKey != ""),
		solverapi.Setting("ExistingToken", task.ExistingToken != ""),
	); err != nil {
		return awsWAFTask{}, err
	}

	taskType, err := awsWAFVariant.For(task.Proxy)
	if err != nil {
		return awsWAFTask{}, err
	}

	return awsWAFTask{
		Type:        taskType,
		WebsiteURL:  task.WebsiteURL,
		ProxyFields: solverapi.NewProxyFields(task.Proxy),
	}, nil
}
