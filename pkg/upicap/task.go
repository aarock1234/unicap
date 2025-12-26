package upicap

// Task represents a captcha solving task
type Task interface {
	// Type returns the captcha type identifier
	Type() TaskType

	// MarshalJSON converts the task to provider-specific JSON
	MarshalJSON() ([]byte, error)
}

// TaskType identifies the kind of captcha
type TaskType string

const (
	TaskTypeReCaptchaV2           TaskType = "recaptcha_v2"
	TaskTypeReCaptchaV3           TaskType = "recaptcha_v3"
	TaskTypeReCaptchaV2Enterprise TaskType = "recaptcha_v2_enterprise"
	TaskTypeReCaptchaV3Enterprise TaskType = "recaptcha_v3_enterprise"
	TaskTypeHCaptcha              TaskType = "hcaptcha"
	TaskTypeFunCaptcha            TaskType = "funcaptcha"
	TaskTypeTurnstile             TaskType = "turnstile"
	TaskTypeCloudflareChallenge   TaskType = "cloudflare_challenge"
	TaskTypeDataDome              TaskType = "datadome"
	TaskTypeGeeTest               TaskType = "geetest"
	TaskTypeGeeTestV4             TaskType = "geetest_v4"
	TaskTypeImageToText           TaskType = "image_to_text"
)
