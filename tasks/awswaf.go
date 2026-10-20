package tasks

import (
	"github.com/aarock1234/unicap/v3"
)

// AWSWAFTask solves an Amazon AWS WAF captcha. Providers return either the
// aws-waf-token value in unicap.Solution.Token or a voucher pair readable
// with unicap.DecodeSolution[AWSWAFSolution].
//
// Providers differ in which inputs they need: CapSolver can work from
// WebsiteURL alone, while 2Captcha and Anti-Captcha need Key, IV, and
// Context from the page's window.gokuProps.
type AWSWAFTask struct {
	// WebsiteURL is the protected page.
	WebsiteURL string

	// Key is the gokuProps "key" value.
	Key string

	// IV is the gokuProps "iv" value.
	IV string

	// Context is the gokuProps "context" value.
	Context string

	// ChallengeScript is the URL of challenge.js.
	ChallengeScript string

	// CaptchaScript is the URL of captcha.js.
	CaptchaScript string

	// JSAPIScript is the URL of jsapi.js, used by sites that integrate the
	// captcha through the JavaScript API.
	JSAPIScript string

	// ProblemURL is the URL of the captcha "problem" endpoint.
	ProblemURL string

	// APIKey is the integration API key some sites expose.
	APIKey string

	// ExistingToken is a previous aws-waf-token to refresh.
	ExistingToken string

	// Proxy is the optional proxy the worker loads the page through.
	Proxy *unicap.Proxy
}

// Type returns unicap.TaskTypeAWSWAF.
func (t *AWSWAFTask) Type() unicap.TaskType {
	return unicap.TaskTypeAWSWAF
}

// Validate reports a missing URL and malformed proxies.
func (t *AWSWAFTask) Validate() error {
	return validate(
		required("WebsiteURL", t.WebsiteURL),
		optionalProxy(t.Proxy),
	)
}

// AWSWAFSolution is the decoded payload of an AWSWAFTask on providers that
// return a voucher instead of a token.
type AWSWAFSolution struct {
	CaptchaVoucher string `json:"captcha_voucher"`
	ExistingToken  string `json:"existing_token"`
}

// AWSWAFClassificationTask answers an AWS WAF image challenge.
type AWSWAFClassificationTask struct {
	// Images are the base64-encoded challenge images: nine for grid
	// questions, otherwise one.
	Images []string

	// Question is the challenge identifier, such as "aws:grid:bed" or
	// "aws:toycarcity:carcity".
	Question string

	// WebsiteURL is the optional protected page.
	WebsiteURL string
}

// Type returns unicap.TaskTypeAWSWAFClassification.
func (t *AWSWAFClassificationTask) Type() unicap.TaskType {
	return unicap.TaskTypeAWSWAFClassification
}

// Validate reports missing required fields.
func (t *AWSWAFClassificationTask) Validate() error {
	return validate(
		requiredItems("Images", t.Images),
		required("Question", t.Question),
	)
}

// AWSWAFClassificationSolution is the decoded payload of an
// AWSWAFClassificationTask. Grid matches are also available in
// unicap.Solution.Indexes.
type AWSWAFClassificationSolution struct {
	// Objects are the zero-based indexes of matching grid images.
	Objects []int `json:"objects"`

	// Box is the [x, y] target for "carcity" questions.
	Box []float64 `json:"box"`

	// Distance is the slide distance for "bifurcatedzoo" questions.
	Distance float64 `json:"distance"`
}
