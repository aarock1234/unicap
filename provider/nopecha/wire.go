package nopecha

import (
	"encoding/json"
	"regexp"
)

// versionPath is the root of the native v1 API, which every endpoint sits
// under.
const versionPath = "/v1/"

// The token endpoints of the native v1 API. A job is submitted to its
// endpoint with POST and its result is fetched from the same endpoint with
// GET.
const (
	endpointHCaptcha    = "token/hcaptcha"
	endpointReCaptchaV2 = "token/recaptcha2"
	endpointReCaptchaV3 = "token/recaptcha3"
	endpointTurnstile   = "token/turnstile"
)

// jobIDParameter is the query parameter that names the job whose result is
// fetched.
const jobIDParameter = "id"

// incompleteJobCode is the error code of a job that is not finished. With
// HTTP status 409 on a fetch it means "ask again later", not a failure.
const incompleteJobCode = 14

// endpointPattern is the form of an endpoint: segments of lowercase letters,
// digits and underscores, joined by single slashes. Every endpoint NopeCHA
// documents has it, the registered ones included. An endpoint becomes part
// of the request path, so nothing else is accepted from a raw task or a task
// ID.
var endpointPattern = regexp.MustCompile(`^[a-z0-9_]+(/[a-z0-9_]+)*$`)

// endpointPath returns the path of an endpoint.
func endpointPath(endpoint string) string {
	return versionPath + endpoint
}

// submitReply is the reply to a submitted job: Data is the job ID.
type submitReply struct {
	Data string `json:"data"`
}

// resultReply is the reply to a fetch of a finished job. Data is the
// solution, kept as received because its shape depends on the endpoint: a
// token endpoint answers with the token as a string.
type resultReply struct {
	Data json.RawMessage `json:"data"`
}

// errorBody is the body of an error reply. Type is null or says more about
// the error.
type errorBody struct {
	Message string `json:"message"`
	Code    int    `json:"error"`
	Type    string `json:"type"`
}
