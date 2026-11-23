package captchaai

import (
	"encoding/json"
	"regexp"
)

// Paths of the native API. A task is submitted to one and its result is
// fetched from the other.
const (
	submitPath = "/in.php"
	resultPath = "/res.php"
)

// The form fields of the protocol itself, which every request or every
// result fetch carries. The fields of a task are named by its mapper.
const (
	fieldKey    = "key"
	fieldMethod = "method"
	fieldJSON   = "json"
	fieldAction = "action"
	fieldID     = "id"
)

// flagSet is the value of a form field that switches something on, such as
// the JSON reply format.
const flagSet = "1"

// actionGet is the action that fetches the result of a task.
const actionGet = "get"

// The methods that select what CaptchaAI solves. reCAPTCHA v2, v3 and
// Enterprise share one method and differ by flags; an image to read and an
// image grid share another.
const (
	methodReCaptcha           = "userrecaptcha"
	methodTurnstile           = "turnstile"
	methodCloudflareChallenge = "cloudflare_challenge"
	methodGeeTest             = "geetest"
	methodFriendlyCaptcha     = "friendly_captcha"
	methodBase64              = "base64"
)

// statusOK is the status of a reply that carries what was asked for. Every
// other reply with a status carries a code instead.
const statusOK = 1

// codeNotReady is the code of a task that is still being solved, in
// CaptchaAI's own spelling. It arrives like an error code but means "ask
// again later".
const codeNotReady = "CAPCHA_NOT_READY"

// vendorTaskIDPattern is the form of the task ID CaptchaAI issues: decimal
// digits, kept as a string because an ID may start with a zero.
var vendorTaskIDPattern = regexp.MustCompile(`^[0-9]+$`)

// reply is the JSON body CaptchaAI answers with, covering every documented
// shape. A submission and most solved tasks answer with Status and Request;
// the Enterprise reCAPTCHA and Cloudflare Challenge solutions arrive in
// Result with UserAgent; a solved GeeTest task answers with Challenge,
// Validate and Seccode and no status at all.
type reply struct {
	// Status is 1 for a reply that carries what was asked for, and 0 for
	// one that carries a code.
	Status int `json:"status"`

	// Request is the task ID, the solution or a code. It is kept as
	// received because it is a string on most replies only.
	Request json.RawMessage `json:"request"`

	// Result is the solution of the task types that also report a
	// User-Agent.
	Result string `json:"result"`

	// UserAgent is the User-Agent the solver used, which the caller has to
	// reuse with the solution.
	UserAgent string `json:"user_agent"`

	// Challenge is the challenge of a solved GeeTest task.
	Challenge string `json:"challenge"`

	// Validate is the validation value of a solved GeeTest task.
	Validate string `json:"validate"`

	// Seccode is the security code of a solved GeeTest task.
	Seccode string `json:"seccode"`
}

// requestText returns Request when it is a string, and an empty string
// otherwise.
func (r reply) requestText() string {
	var text string

	_ = json.Unmarshal(r.Request, &text)

	return text
}

// solved reports whether the reply carries a solution: it has the success
// status, or the security code of a GeeTest task, whose replies carry no
// status.
func (r reply) solved() bool {
	return r.Status == statusOK || r.Seccode != ""
}
