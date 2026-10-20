package capsolver

import (
	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

type turnstileTask struct {
	Type       string             `json:"type"`
	WebsiteURL string             `json:"websiteURL"`
	WebsiteKey string             `json:"websiteKey"`
	Metadata   *turnstileMetadata `json:"metadata,omitempty"`
}

type turnstileMetadata struct {
	Action string `json:"action,omitempty"`
	CData  string `json:"cdata,omitempty"`
}

type cloudflareChallengeTask struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	HTML       string `json:"html,omitempty"`
	UserAgent  string `json:"userAgent,omitempty"`
	Proxy      string `json:"proxy"`
}

type dataDomeTask struct {
	Type       string `json:"type"`
	CaptchaURL string `json:"captchaUrl"`
	UserAgent  string `json:"userAgent"`
	HTML       string `json:"html,omitempty"`
	Proxy      string `json:"proxy"`
}

type awsWAFTask struct {
	Type             string `json:"type"`
	WebsiteURL       string `json:"websiteURL"`
	AWSKey           string `json:"awsKey,omitempty"`
	AWSIV            string `json:"awsIv,omitempty"`
	AWSContext       string `json:"awsContext,omitempty"`
	AWSChallengeJS   string `json:"awsChallengeJS,omitempty"`
	AWSAPIJS         string `json:"awsApiJs,omitempty"`
	AWSProblemURL    string `json:"awsProblemUrl,omitempty"`
	AWSAPIKey        string `json:"awsApiKey,omitempty"`
	AWSExistingToken string `json:"awsExistingToken,omitempty"`
	Proxy            string `json:"proxy,omitempty"`
}

type botDeflectorTask struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	Domain     string `json:"domain"`
	FlowToken  string `json:"flowToken"`
}

func mapTurnstile(task *tasks.TurnstileTask) (turnstileTask, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("Proxy", task.Proxy != nil),
		solverapi.Setting("PageData", task.PageData != ""),
	); err != nil {
		return turnstileTask{}, err
	}

	wire := turnstileTask{
		Type:       "AntiTurnstileTaskProxyLess",
		WebsiteURL: task.WebsiteURL,
		WebsiteKey: task.WebsiteKey,
	}

	if task.Action != "" || task.CData != "" {
		wire.Metadata = &turnstileMetadata{
			Action: task.Action,
			CData:  task.CData,
		}
	}

	return wire, nil
}

func mapCloudflareChallenge(task *tasks.CloudflareChallengeTask) (cloudflareChallengeTask, error) {
	return cloudflareChallengeTask{
		Type:       "AntiCloudflareTask",
		WebsiteURL: task.WebsiteURL,
		HTML:       task.HTML,
		UserAgent:  task.UserAgent,
		Proxy:      solverapi.ProxyURL(task.Proxy),
	}, nil
}

func mapDataDome(task *tasks.DataDomeTask) (dataDomeTask, error) {
	return dataDomeTask{
		Type:       "DatadomeSliderTask",
		CaptchaURL: task.CaptchaURL,
		UserAgent:  task.UserAgent,
		HTML:       task.HTML,
		Proxy:      solverapi.ProxyURL(task.Proxy),
	}, nil
}

func mapAWSWAF(task *tasks.AWSWAFTask) (awsWAFTask, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("CaptchaScript", task.CaptchaScript != ""),
	); err != nil {
		return awsWAFTask{}, err
	}

	taskType, err := variant("AntiAwsWafTask").For(task.Proxy)
	if err != nil {
		return awsWAFTask{}, err
	}

	return awsWAFTask{
		Type:             taskType,
		WebsiteURL:       task.WebsiteURL,
		AWSKey:           task.Key,
		AWSIV:            task.IV,
		AWSContext:       task.Context,
		AWSChallengeJS:   task.ChallengeScript,
		AWSAPIJS:         task.JSAPIScript,
		AWSProblemURL:    task.ProblemURL,
		AWSAPIKey:        task.APIKey,
		AWSExistingToken: task.ExistingToken,
		Proxy:            solverapi.ProxyURL(task.Proxy),
	}, nil
}

func mapBotDeflector(task *tasks.BotDeflectorTask) (botDeflectorTask, error) {
	return botDeflectorTask{
		Type:       "AntiBotdeflectorTaskProxyLess",
		WebsiteURL: task.WebsiteURL,
		Domain:     task.Domain,
		FlowToken:  task.FlowToken,
	}, nil
}
