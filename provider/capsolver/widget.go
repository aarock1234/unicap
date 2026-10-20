package capsolver

import (
	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

type geeTestTask struct {
	Type                      string `json:"type"`
	WebsiteURL                string `json:"websiteURL"`
	GT                        string `json:"gt"`
	Challenge                 string `json:"challenge"`
	GeetestAPIServerSubdomain string `json:"geetestApiServerSubdomain,omitempty"`
}

type geeTestV4Task struct {
	Type                      string `json:"type"`
	WebsiteURL                string `json:"websiteURL"`
	CaptchaID                 string `json:"captchaId"`
	RiskType                  string `json:"riskType,omitempty"`
	GeetestAPIServerSubdomain string `json:"geetestApiServerSubdomain,omitempty"`
}

type mtCaptchaTask struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	WebsiteKey string `json:"websiteKey"`
	Proxy      string `json:"proxy,omitempty"`
}

// geeTestTaskType is CapSolver's only GeeTest wire type; it covers v3 and v4
// and has no proxy variant.
const geeTestTaskType = "GeeTestTaskProxyLess"

func mapGeeTest(task *tasks.GeeTestTask) (geeTestTask, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("Proxy", task.Proxy != nil),
		solverapi.Setting("UserAgent", task.UserAgent != ""),
	); err != nil {
		return geeTestTask{}, err
	}

	return geeTestTask{
		Type:                      geeTestTaskType,
		WebsiteURL:                task.WebsiteURL,
		GT:                        task.GT,
		Challenge:                 task.Challenge,
		GeetestAPIServerSubdomain: task.APIServerSubdomain,
	}, nil
}

func mapGeeTestV4(task *tasks.GeeTestV4Task) (geeTestV4Task, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("Proxy", task.Proxy != nil),
		solverapi.Setting("UserAgent", task.UserAgent != ""),
	); err != nil {
		return geeTestV4Task{}, err
	}

	return geeTestV4Task{
		Type:                      geeTestTaskType,
		WebsiteURL:                task.WebsiteURL,
		CaptchaID:                 task.CaptchaID,
		RiskType:                  task.RiskType,
		GeetestAPIServerSubdomain: task.APIServerSubdomain,
	}, nil
}

func mapMTCaptcha(task *tasks.MTCaptchaTask) (mtCaptchaTask, error) {
	taskType, err := variant("MtCaptchaTask").For(task.Proxy)
	if err != nil {
		return mtCaptchaTask{}, err
	}

	return mtCaptchaTask{
		Type:       taskType,
		WebsiteURL: task.WebsiteURL,
		WebsiteKey: task.WebsiteKey,
		Proxy:      solverapi.ProxyURL(task.Proxy),
	}, nil
}
