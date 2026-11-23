package easycaptcha

import (
	"cmp"
	"encoding/json"
	"slices"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/internal/httpapi"
)

// unknownFailureCode stands in for the code of a failed task that carries no
// error object.
const unknownFailureCode = "TASK_FAILED"

// classifyTask sorts the reply to a request easycaptcha.to answers with a
// task object; success lists the endpoint's success statuses. It returns the
// task while it is queued or running and once it is completed. A rejected
// request, a reply without a task, and a failed task are errors.
func (c *Client) classifyTask(request httpapi.Request, response *httpapi.Response, success []int) (taskObject, error) {
	task, err := classifyReply[taskObject](c, request, response, success)
	if err != nil {
		return taskObject{}, err
	}

	if task.ID == "" {
		return taskObject{}, c.unexplained(request, response)
	}

	if task.Status == statusFailed {
		return taskObject{}, c.taskFailure(task.Error)
	}

	return task, nil
}

// classifyReply sorts a reply into the success body T or the error the reply
// stands for. A success status whose body is not a T is unexplained.
func classifyReply[T any](c *Client, request httpapi.Request, response *httpapi.Response, success []int) (T, error) {
	var body T

	if err := c.requestError(request, response, success); err != nil {
		return body, err
	}

	if err := json.Unmarshal(response.Body, &body); err != nil {
		return body, c.unexplained(request, response)
	}

	return body, nil
}

// requestError returns the error a reply stands for, and nil when its status
// is one of the endpoint's success statuses. A problem body with a
// documented code is a final provider error classified by requestErrorRules;
// anything else is unexplained, as easycaptcha.to asks clients to treat an
// unknown problem by its HTTP status.
func (c *Client) requestError(request httpapi.Request, response *httpapi.Response, success []int) error {
	if slices.Contains(success, response.StatusCode) {
		return nil
	}

	var body problem
	if err := json.Unmarshal(response.Body, &body); err != nil {
		return c.unexplained(request, response)
	}

	sentinel, documented := requestErrorRules[body.Code]
	if !documented {
		return c.unexplained(request, response)
	}

	return c.api.Error(body.Code, cmp.Or(body.Detail, body.Title), sentinel)
}

// taskFailure builds the final provider error of a failed task, classified
// by taskFailureRules. Any other code, and a failure without one, maps to
// unicap.ErrTaskFailed.
func (c *Client) taskFailure(failure problem) *unicap.Error {
	sentinel, documented := taskFailureRules[failure.Code]
	if !documented {
		sentinel = unicap.ErrTaskFailed
	}

	return c.api.Error(
		cmp.Or(failure.Code, unknownFailureCode),
		cmp.Or(failure.Detail, failure.Title),
		sentinel,
	)
}

// unexplained reports a reply easycaptcha.to did not explain with a
// documented problem code. It is a plain error, which the poller treats as
// transient, for a server error or a rate limit, and a final provider error
// named after the HTTP status otherwise.
func (c *Client) unexplained(request httpapi.Request, response *httpapi.Response) error {
	return c.api.StatusError(request.Method+" "+request.Path, response.StatusCode)
}

// taskStatus converts the status of a task that did not fail.
func taskStatus(status string) unicap.TaskStatus {
	switch status {
	case statusCompleted:
		return unicap.TaskStatusReady
	case statusQueued:
		return unicap.TaskStatusPending
	default:
		return unicap.TaskStatusProcessing
	}
}
