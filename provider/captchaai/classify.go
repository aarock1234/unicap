package captchaai

import (
	"encoding/json"

	"github.com/aarock1234/unicap/v3/internal/httpapi"
)

// classifySubmit sorts the reply to a submitted task and returns the ID of
// the task CaptchaAI created. A reply that carries a code, and one without a
// task ID of the documented form, are errors.
func (c *Client) classifySubmit(request httpapi.Request, response *httpapi.Response) (string, error) {
	answer, err := c.decodeReply(request, response)
	if err != nil {
		return "", err
	}

	if answer.Status != statusOK {
		return "", c.failure(request, response, answer)
	}

	vendorTaskID := answer.requestText()
	if !vendorTaskIDPattern.MatchString(vendorTaskID) {
		return "", c.unexplained(request, response)
	}

	return vendorTaskID, nil
}

// classifyResult sorts the reply to a fetch of a task's result. It returns
// the reply and true once the task is solved, and false while CaptchaAI
// reports the task as not ready. A reply that carries any other code is an
// error.
func (c *Client) classifyResult(request httpapi.Request, response *httpapi.Response) (reply, bool, error) {
	answer, err := c.decodeReply(request, response)
	if err != nil {
		return reply{}, false, err
	}

	if answer.solved() {
		return answer, true, nil
	}

	if answer.requestText() == codeNotReady {
		return reply{}, false, nil
	}

	return reply{}, false, c.failure(request, response, answer)
}

// decodeReply decodes the body of a reply, whatever its HTTP status: the
// body says whether the request succeeded. A body that is not a reply is
// unexplained.
func (c *Client) decodeReply(request httpapi.Request, response *httpapi.Response) (reply, error) {
	var answer reply
	if err := json.Unmarshal(response.Body, &answer); err != nil {
		return reply{}, c.unexplained(request, response)
	}

	return answer, nil
}

// failure returns the error a reply without a success status stands for: a
// final provider error named after the code in its request field and
// classified by errorRules. CaptchaAI describes a failure by its code alone,
// so the error has no message, and a code errorRules does not list carries
// no sentinel. A reply without a code is unexplained.
func (c *Client) failure(request httpapi.Request, response *httpapi.Response, answer reply) error {
	code := answer.requestText()
	if code == "" {
		return c.unexplained(request, response)
	}

	return c.api.Error(code, "", errorRules[code])
}

// unexplained reports a reply CaptchaAI did not explain with a code. It is a
// plain error, which the poller treats as transient, for a server error or a
// rate limit, and a final provider error named after the HTTP status
// otherwise.
func (c *Client) unexplained(request httpapi.Request, response *httpapi.Response) error {
	return c.api.StatusError(request.Method+" "+request.Path, response.StatusCode)
}
