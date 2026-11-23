package noptcha

import (
	"cmp"
	"encoding/json"
	"net/http"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/internal/httpapi"
)

// unknownFailureCode stands in for the code of a failed task that carries no
// error object.
const unknownFailureCode = "TASK_FAILED"

// classifyTask sorts the reply to a request Noptcha answers with a task
// object. It returns the task while it is processing and once it is solved.
// A rejected request, a reply without a task, and a failed task are errors.
func (c *Client) classifyTask(request httpapi.Request, response *httpapi.Response) (taskObject, error) {
	task, err := classifyReply[taskObject](c, request, response)
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
func classifyReply[T any](c *Client, request httpapi.Request, response *httpapi.Response) (T, error) {
	var body T

	if err := c.requestError(request, response); err != nil {
		return body, err
	}

	if err := json.Unmarshal(response.Body, &body); err != nil {
		return body, c.unexplained(request, response)
	}

	return body, nil
}

// requestError returns the error a reply with an error status stands for,
// and nil for a success status. A problem body with a documented code is a
// provider error classified by requestErrorRules; anything else is
// unexplained, as Noptcha asks clients to treat an unknown code by its HTTP
// status.
func (c *Client) requestError(request httpapi.Request, response *httpapi.Response) error {
	if isSuccess(response.StatusCode) {
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

// isSuccess reports whether an HTTP status carries a success body. Noptcha
// answers 200 for a final result and 202 for a task that is still
// processing.
func isSuccess(statusCode int) bool {
	return statusCode == http.StatusOK || statusCode == http.StatusAccepted
}

// taskFailure builds the final provider error of a failed task, classified
// by taskFailureRules. A code Noptcha added since maps to
// unicap.ErrTaskFailed.
func (c *Client) taskFailure(failure taskError) *unicap.Error {
	sentinel, documented := taskFailureRules[failure.Code]
	if !documented {
		sentinel = unicap.ErrTaskFailed
	}

	return c.api.Error(cmp.Or(failure.Code, unknownFailureCode), failure.Message, sentinel)
}

// unexplained reports a reply Noptcha did not explain with a documented
// error code. It is a plain error, which the poller treats as transient, for
// a server error or a rate limit, and a final provider error named after the
// HTTP status otherwise.
func (c *Client) unexplained(request httpapi.Request, response *httpapi.Response) error {
	return c.api.StatusError(request.Method+" "+request.Path, response.StatusCode)
}

// taskStatus converts the status of a task that did not fail.
func taskStatus(status string) unicap.TaskStatus {
	if status == statusSolved {
		return unicap.TaskStatusReady
	}

	return unicap.TaskStatusProcessing
}
