package tasks

import (
	"github.com/aarock1234/unicap/v3"
)

// All returns a zero value of every typed task, in documentation order.
// RawTask is not included.
func All() []unicap.Task {
	return []unicap.Task{
		&ReCaptchaV2Task{},
		&ReCaptchaV2EnterpriseTask{},
		&ReCaptchaV3Task{},
		&ReCaptchaV3EnterpriseTask{},
		&ReCaptchaClassificationTask{},
		&HCaptchaTask{},
		&FunCaptchaTask{},
		&TurnstileTask{},
		&CloudflareChallengeTask{},
		&DataDomeTask{},
		&GeeTestTask{},
		&GeeTestV4Task{},
		&AWSWAFTask{},
		&AWSWAFClassificationTask{},
		&MTCaptchaTask{},
		&FriendlyCaptchaTask{},
		&LeminTask{},
		&CutCaptchaTask{},
		&ProsopoTask{},
		&AltchaTask{},
		&CapyTask{},
		&KeyCaptchaTask{},
		&CyberSiAraTask{},
		&AtbCaptchaTask{},
		&TencentTask{},
		&CaptchaFoxTask{},
		&YandexSmartCaptchaTask{},
		&BotDeflectorTask{},
		&VKCaptchaTask{},
		&YidunTask{},
		&BinanceTask{},
		&HuntTask{},
		&TSPDTask{},
		&BasiliskTask{},
		&ImpervaTask{},
		&AlibabaTask{},
		&AntiGateTask{},
		&ImageToTextTask{},
		&TextCaptchaTask{},
		&AudioTask{},
		&CoordinatesTask{},
		&GridTask{},
		&RotateTask{},
		&DrawAroundTask{},
		&BoundingBoxTask{},
		&DragAndDropTask{},
		&TemuTask{},
		&VKCaptchaImageTask{},
		&YandexImageTask{},
		&YandexPuzzleTask{},
		&VisionEngineTask{},
	}
}

// RawTask is asserted separately because All lists typed tasks only.
var _ unicap.Task = (*RawTask)(nil)
