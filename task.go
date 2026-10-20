package unicap

// Task is a captcha to solve. The tasks package provides an implementation
// for every captcha family unicap models.
type Task interface {
	// Type returns the captcha family identifier.
	Type() TaskType

	// Validate reports whether the task carries the inputs every provider
	// requires. Providers may reject further fields they cannot honour.
	Validate() error
}

// TaskType identifies a captcha family independently of any provider.
type TaskType string

// Token, cookie, and challenge captchas.
const (
	TaskTypeReCaptchaV2           TaskType = "recaptcha_v2"
	TaskTypeReCaptchaV2Enterprise TaskType = "recaptcha_v2_enterprise"
	TaskTypeReCaptchaV3           TaskType = "recaptcha_v3"
	TaskTypeReCaptchaV3Enterprise TaskType = "recaptcha_v3_enterprise"
	TaskTypeHCaptcha              TaskType = "hcaptcha"
	TaskTypeFunCaptcha            TaskType = "funcaptcha"
	TaskTypeTurnstile             TaskType = "turnstile"
	TaskTypeCloudflareChallenge   TaskType = "cloudflare_challenge"
	TaskTypeDataDome              TaskType = "datadome"
	TaskTypeGeeTest               TaskType = "geetest"
	TaskTypeGeeTestV4             TaskType = "geetest_v4"
	TaskTypeAWSWAF                TaskType = "aws_waf"
	TaskTypeMTCaptcha             TaskType = "mtcaptcha"
	TaskTypeFriendlyCaptcha       TaskType = "friendly_captcha"
	TaskTypeLemin                 TaskType = "lemin"
	TaskTypeCutCaptcha            TaskType = "cutcaptcha"
	TaskTypeProsopo               TaskType = "prosopo"
	TaskTypeAltcha                TaskType = "altcha"
	TaskTypeCapy                  TaskType = "capy"
	TaskTypeKeyCaptcha            TaskType = "keycaptcha"
	TaskTypeCyberSiAra            TaskType = "cybersiara"
	TaskTypeAtbCaptcha            TaskType = "atbcaptcha"
	TaskTypeTencent               TaskType = "tencent"
	TaskTypeCaptchaFox            TaskType = "captchafox"
	TaskTypeVKCaptcha             TaskType = "vk_captcha"
	TaskTypeYidun                 TaskType = "yidun"
	TaskTypeBinance               TaskType = "binance"
	TaskTypeHunt                  TaskType = "hunt"
	TaskTypeTSPD                  TaskType = "tspd"
	TaskTypeBasilisk              TaskType = "basilisk"
	TaskTypeImperva               TaskType = "imperva"
	TaskTypeAlibaba               TaskType = "alibaba"
	TaskTypeYandexSmartCaptcha    TaskType = "yandex_smartcaptcha"
	TaskTypeBotDeflector          TaskType = "botdeflector"
	TaskTypeAntiGate              TaskType = "antigate"
)

// Image, audio, and text recognition captchas.
const (
	TaskTypeImageToText             TaskType = "image_to_text"
	TaskTypeText                    TaskType = "text"
	TaskTypeAudio                   TaskType = "audio"
	TaskTypeCoordinates             TaskType = "coordinates"
	TaskTypeGrid                    TaskType = "grid"
	TaskTypeRotate                  TaskType = "rotate"
	TaskTypeDrawAround              TaskType = "draw_around"
	TaskTypeBoundingBox             TaskType = "bounding_box"
	TaskTypeDragAndDrop             TaskType = "drag_and_drop"
	TaskTypeReCaptchaClassification TaskType = "recaptcha_classification"
	TaskTypeAWSWAFClassification    TaskType = "aws_waf_classification"
	TaskTypeVisionEngine            TaskType = "vision_engine"
	TaskTypeTemu                    TaskType = "temu"
	TaskTypeVKCaptchaImage          TaskType = "vk_captcha_image"
	TaskTypeYandexImage             TaskType = "yandex_image"
	TaskTypeYandexPuzzle            TaskType = "yandex_puzzle"
)

// TaskTypeRaw identifies a provider-specific passthrough task.
const TaskTypeRaw TaskType = "raw"
