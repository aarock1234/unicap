package tasks

import "github.com/aarock1234/unicap"

// Compile-time assertions that every task type implements unicap.Task.
var (
	_ unicap.Task = (*ReCaptchaV2Task)(nil)
	_ unicap.Task = (*ReCaptchaV3Task)(nil)
	_ unicap.Task = (*ReCaptchaV2EnterpriseTask)(nil)
	_ unicap.Task = (*ReCaptchaV3EnterpriseTask)(nil)
	_ unicap.Task = (*HCaptchaTask)(nil)
	_ unicap.Task = (*FunCaptchaTask)(nil)
	_ unicap.Task = (*TurnstileTask)(nil)
	_ unicap.Task = (*CloudflareChallengeTask)(nil)
	_ unicap.Task = (*DataDomeTask)(nil)
	_ unicap.Task = (*GeeTestTask)(nil)
	_ unicap.Task = (*GeeTestV4Task)(nil)
	_ unicap.Task = (*ImageToTextTask)(nil)
	_ unicap.Task = (*AWSWAFTask)(nil)
	_ unicap.Task = (*MTCaptchaTask)(nil)
	_ unicap.Task = (*FriendlyCaptchaTask)(nil)
	_ unicap.Task = (*LeminTask)(nil)
	_ unicap.Task = (*CutCaptchaTask)(nil)
	_ unicap.Task = (*TextCaptchaTask)(nil)
	_ unicap.Task = (*ProsopoTask)(nil)
	_ unicap.Task = (*AltchaTask)(nil)
	_ unicap.Task = (*RawTask)(nil)
)
