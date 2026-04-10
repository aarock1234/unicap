package unicap

// Task represents a captcha solving task
type Task interface {
	// Type returns the captcha type identifier
	Type() TaskType

	// Validate ensures the task has the required inputs.
	Validate() error
}

// TaskType identifies the kind of captcha
type TaskType string

const (
	// TaskTypeReCaptchaV2 identifies a ReCaptcha v2 task.
	TaskTypeReCaptchaV2 TaskType = "recaptcha_v2"
	// TaskTypeReCaptchaV3 identifies a ReCaptcha v3 task.
	TaskTypeReCaptchaV3 TaskType = "recaptcha_v3"
	// TaskTypeReCaptchaV2Enterprise identifies a ReCaptcha v2 Enterprise task.
	TaskTypeReCaptchaV2Enterprise TaskType = "recaptcha_v2_enterprise"
	// TaskTypeReCaptchaV3Enterprise identifies a ReCaptcha v3 Enterprise task.
	TaskTypeReCaptchaV3Enterprise TaskType = "recaptcha_v3_enterprise"
	// TaskTypeHCaptcha identifies an hCaptcha task.
	TaskTypeHCaptcha TaskType = "hcaptcha"
	// TaskTypeFunCaptcha identifies a FunCaptcha task.
	TaskTypeFunCaptcha TaskType = "funcaptcha"
	// TaskTypeTurnstile identifies a Cloudflare Turnstile task.
	TaskTypeTurnstile TaskType = "turnstile"
	// TaskTypeCloudflareChallenge identifies a Cloudflare challenge task.
	TaskTypeCloudflareChallenge TaskType = "cloudflare_challenge"
	// TaskTypeDataDome identifies a DataDome slider task.
	TaskTypeDataDome TaskType = "datadome"
	// TaskTypeGeeTest identifies a GeeTest v3 task.
	TaskTypeGeeTest TaskType = "geetest"
	// TaskTypeGeeTestV4 identifies a GeeTest v4 task.
	TaskTypeGeeTestV4 TaskType = "geetest_v4"
	// TaskTypeImageToText identifies an image-to-text task.
	TaskTypeImageToText TaskType = "image_to_text"
)
