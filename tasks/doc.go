// Package tasks provides the captcha task types unicap models, one per
// captcha family, together with typed solution payloads for families whose
// answers do not fit unicap.Solution's common fields.
//
// Every task validates the inputs all providers require. Providers reject
// options they cannot honour with unicap.ErrUnsupportedOption rather than
// silently dropping them.
package tasks

import (
	"github.com/aarock1234/unicap/v3"
)

// Compile-time assertions that every task type implements unicap.Task.
var (
	_ unicap.Task = (*ReCaptchaV2Task)(nil)
	_ unicap.Task = (*ReCaptchaV2EnterpriseTask)(nil)
	_ unicap.Task = (*ReCaptchaV3Task)(nil)
	_ unicap.Task = (*ReCaptchaV3EnterpriseTask)(nil)
	_ unicap.Task = (*ReCaptchaClassificationTask)(nil)
	_ unicap.Task = (*HCaptchaTask)(nil)
	_ unicap.Task = (*FunCaptchaTask)(nil)
	_ unicap.Task = (*TurnstileTask)(nil)
	_ unicap.Task = (*CloudflareChallengeTask)(nil)
	_ unicap.Task = (*DataDomeTask)(nil)
	_ unicap.Task = (*GeeTestTask)(nil)
	_ unicap.Task = (*GeeTestV4Task)(nil)
	_ unicap.Task = (*AWSWAFTask)(nil)
	_ unicap.Task = (*AWSWAFClassificationTask)(nil)
	_ unicap.Task = (*MTCaptchaTask)(nil)
	_ unicap.Task = (*FriendlyCaptchaTask)(nil)
	_ unicap.Task = (*LeminTask)(nil)
	_ unicap.Task = (*CutCaptchaTask)(nil)
	_ unicap.Task = (*ProsopoTask)(nil)
	_ unicap.Task = (*AltchaTask)(nil)
	_ unicap.Task = (*CapyTask)(nil)
	_ unicap.Task = (*KeyCaptchaTask)(nil)
	_ unicap.Task = (*CyberSiAraTask)(nil)
	_ unicap.Task = (*AtbCaptchaTask)(nil)
	_ unicap.Task = (*TencentTask)(nil)
	_ unicap.Task = (*CaptchaFoxTask)(nil)
	_ unicap.Task = (*YandexSmartCaptchaTask)(nil)
	_ unicap.Task = (*BotDeflectorTask)(nil)
	_ unicap.Task = (*VKCaptchaTask)(nil)
	_ unicap.Task = (*YidunTask)(nil)
	_ unicap.Task = (*BinanceTask)(nil)
	_ unicap.Task = (*HuntTask)(nil)
	_ unicap.Task = (*TSPDTask)(nil)
	_ unicap.Task = (*BasiliskTask)(nil)
	_ unicap.Task = (*ImpervaTask)(nil)
	_ unicap.Task = (*AlibabaTask)(nil)
	_ unicap.Task = (*AntiGateTask)(nil)
	_ unicap.Task = (*ImageToTextTask)(nil)
	_ unicap.Task = (*TextCaptchaTask)(nil)
	_ unicap.Task = (*AudioTask)(nil)
	_ unicap.Task = (*CoordinatesTask)(nil)
	_ unicap.Task = (*GridTask)(nil)
	_ unicap.Task = (*RotateTask)(nil)
	_ unicap.Task = (*DrawAroundTask)(nil)
	_ unicap.Task = (*BoundingBoxTask)(nil)
	_ unicap.Task = (*DragAndDropTask)(nil)
	_ unicap.Task = (*TemuTask)(nil)
	_ unicap.Task = (*VKCaptchaImageTask)(nil)
	_ unicap.Task = (*YandexImageTask)(nil)
	_ unicap.Task = (*YandexPuzzleTask)(nil)
	_ unicap.Task = (*VisionEngineTask)(nil)
	_ unicap.Task = (*RawTask)(nil)
)
