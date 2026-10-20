package capsolver

import (
	"fmt"
	"net/http"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

type reCaptchaV2Task struct {
	Type                string         `json:"type"`
	WebsiteURL          string         `json:"websiteURL"`
	WebsiteKey          string         `json:"websiteKey"`
	IsInvisible         bool           `json:"isInvisible,omitempty"`
	PageAction          string         `json:"pageAction,omitempty"`
	RecaptchaDataSValue string         `json:"recaptchaDataSValue,omitempty"`
	EnterprisePayload   map[string]any `json:"enterprisePayload,omitempty"`
	IsSession           bool           `json:"isSession,omitempty"`
	APIDomain           string         `json:"apiDomain,omitempty"`
	Cookies             []cookie       `json:"cookies,omitempty"`
	Proxy               string         `json:"proxy,omitempty"`
}

type reCaptchaV3Task struct {
	Type              string         `json:"type"`
	WebsiteURL        string         `json:"websiteURL"`
	WebsiteKey        string         `json:"websiteKey"`
	PageAction        string         `json:"pageAction,omitempty"`
	EnterprisePayload map[string]any `json:"enterprisePayload,omitempty"`
	IsSession         bool           `json:"isSession,omitempty"`
	APIDomain         string         `json:"apiDomain,omitempty"`
	Proxy             string         `json:"proxy,omitempty"`
}

type reCaptchaClassificationTask struct {
	Type       string `json:"type"`
	Image      string `json:"image"`
	Question   string `json:"question"`
	WebsiteURL string `json:"websiteURL,omitempty"`
	WebsiteKey string `json:"websiteKey,omitempty"`
}

// cookie is CapSolver's cookie object.
type cookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func mapReCaptchaV2(task *tasks.ReCaptchaV2Task) (reCaptchaV2Task, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("UserAgent", task.UserAgent != ""),
	); err != nil {
		return reCaptchaV2Task{}, err
	}

	taskType, err := variant("ReCaptchaV2Task").For(task.Proxy)
	if err != nil {
		return reCaptchaV2Task{}, err
	}

	cookies, err := parseCookies(task.Cookies)
	if err != nil {
		return reCaptchaV2Task{}, err
	}

	return reCaptchaV2Task{
		Type:                taskType,
		WebsiteURL:          task.WebsiteURL,
		WebsiteKey:          task.WebsiteKey,
		IsInvisible:         task.IsInvisible,
		PageAction:          task.PageAction,
		RecaptchaDataSValue: task.DataS,
		IsSession:           task.IsSession,
		APIDomain:           task.APIDomain,
		Cookies:             cookies,
		Proxy:               solverapi.ProxyURL(task.Proxy),
	}, nil
}

func mapReCaptchaV2Enterprise(task *tasks.ReCaptchaV2EnterpriseTask) (reCaptchaV2Task, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("UserAgent", task.UserAgent != ""),
	); err != nil {
		return reCaptchaV2Task{}, err
	}

	taskType, err := variant("ReCaptchaV2EnterpriseTask").For(task.Proxy)
	if err != nil {
		return reCaptchaV2Task{}, err
	}

	cookies, err := parseCookies(task.Cookies)
	if err != nil {
		return reCaptchaV2Task{}, err
	}

	return reCaptchaV2Task{
		Type:              taskType,
		WebsiteURL:        task.WebsiteURL,
		WebsiteKey:        task.WebsiteKey,
		IsInvisible:       task.IsInvisible,
		PageAction:        task.PageAction,
		EnterprisePayload: task.EnterprisePayload,
		IsSession:         task.IsSession,
		APIDomain:         task.APIDomain,
		Cookies:           cookies,
		Proxy:             solverapi.ProxyURL(task.Proxy),
	}, nil
}

func mapReCaptchaV3(task *tasks.ReCaptchaV3Task) (reCaptchaV3Task, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("MinScore", task.MinScore != 0),
	); err != nil {
		return reCaptchaV3Task{}, err
	}

	base := "ReCaptchaV3Task"
	if task.HighScore {
		base = "ReCaptchaV3M1Task"
	}

	taskType, err := variant(base).For(task.Proxy)
	if err != nil {
		return reCaptchaV3Task{}, err
	}

	return reCaptchaV3Task{
		Type:       taskType,
		WebsiteURL: task.WebsiteURL,
		WebsiteKey: task.WebsiteKey,
		PageAction: task.PageAction,
		IsSession:  task.IsSession,
		APIDomain:  task.APIDomain,
		Proxy:      solverapi.ProxyURL(task.Proxy),
	}, nil
}

func mapReCaptchaV3Enterprise(task *tasks.ReCaptchaV3EnterpriseTask) (reCaptchaV3Task, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("MinScore", task.MinScore != 0),
	); err != nil {
		return reCaptchaV3Task{}, err
	}

	base := "ReCaptchaV3EnterpriseTask"
	if task.HighScore {
		base = "ReCaptchaV3EnterpriseM1Task"
	}

	taskType, err := variant(base).For(task.Proxy)
	if err != nil {
		return reCaptchaV3Task{}, err
	}

	return reCaptchaV3Task{
		Type:              taskType,
		WebsiteURL:        task.WebsiteURL,
		WebsiteKey:        task.WebsiteKey,
		PageAction:        task.PageAction,
		EnterprisePayload: task.EnterprisePayload,
		IsSession:         task.IsSession,
		APIDomain:         task.APIDomain,
		Proxy:             solverapi.ProxyURL(task.Proxy),
	}, nil
}

func mapReCaptchaClassification(task *tasks.ReCaptchaClassificationTask) (reCaptchaClassificationTask, error) {
	return reCaptchaClassificationTask{
		Type:       "ReCaptchaV2Classification",
		Image:      task.Image,
		Question:   task.Question,
		WebsiteURL: task.WebsiteURL,
		WebsiteKey: task.WebsiteKey,
	}, nil
}

// parseCookies converts a "name=value; name2=value2" string into CapSolver
// cookie objects.
func parseCookies(text string) ([]cookie, error) {
	if text == "" {
		return nil, nil
	}

	parsed, err := http.ParseCookie(text)
	if err != nil {
		return nil, fmt.Errorf("parsing cookies: %w", unicap.ErrInvalidTask)
	}

	cookies := make([]cookie, 0, len(parsed))
	for _, parsedCookie := range parsed {
		cookies = append(cookies, cookie{
			Name:  parsedCookie.Name,
			Value: parsedCookie.Value,
		})
	}

	return cookies, nil
}
