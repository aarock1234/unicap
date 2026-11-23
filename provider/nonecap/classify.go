package nonecap

import (
	"cmp"
	"encoding/json"
	"net/http"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/internal/httpapi"
)

// unknownFailureCode stands in for the code of a solve that ended without a
// token and carries no error object.
const unknownFailureCode = "SOLVE_FAILED"

// classifySolve sorts the reply to a request NoneCap answers with a solve
// object. It returns the solve while it is in progress and once it is
// solved. A rejected request, a reply without a solve, and a solve that
// ended without a token are errors.
func (c *Client) classifySolve(request httpapi.Request, response *httpapi.Response) (solveObject, error) {
	if err := c.requestError(request, response); err != nil {
		return solveObject{}, err
	}

	var solve solveObject
	if err := json.Unmarshal(response.Body, &solve); err != nil {
		return solveObject{}, c.unexplained(request, response)
	}

	if solve.ID == "" {
		return solveObject{}, c.unexplained(request, response)
	}

	if isFailure(solve.Status) {
		return solveObject{}, c.solveFailure(solve.Error)
	}

	return solve, nil
}

// requestError returns the error a reply with an error status stands for,
// and nil for a success status. An error envelope with a code is a final
// provider error classified by requestErrorRules; a reply without one is
// unexplained.
func (c *Client) requestError(request httpapi.Request, response *httpapi.Response) error {
	if isSuccess(response.StatusCode) {
		return nil
	}

	var body errorEnvelope
	if err := json.Unmarshal(response.Body, &body); err != nil {
		return c.unexplained(request, response)
	}

	failure := body.Error
	if failure.Code == "" {
		return c.unexplained(request, response)
	}

	return c.api.Error(failure.Code, failure.Message, requestErrorRules[failure.Code])
}

// isSuccess reports whether an HTTP status carries a success body. NoneCap
// answers 200 for a solve that has ended and for recorded feedback, and 202
// for a solve that is in progress.
func isSuccess(statusCode int) bool {
	return statusCode == http.StatusOK || statusCode == http.StatusAccepted
}

// isFailure reports whether a solve status is one of the final statuses that
// carry an error instead of a token.
func isFailure(status string) bool {
	switch status {
	case statusFailed, statusCancelled, statusExpired:
		return true
	default:
		return false
	}
}

// solveFailure builds the final provider error of a solve that ended without
// a token.
func (c *Client) solveFailure(failure solveError) *unicap.Error {
	return c.api.Error(cmp.Or(failure.Code, unknownFailureCode), failure.Message, failureSentinel(failure))
}

// failureSentinel classifies a solve failure. A reason listed in
// solveFailureReasonRules decides first, then the code by solveFailureRules.
// A code NoneCap added since is classified by its retryable flag: a failure
// the same request can get past is unicap.ErrUnsolvable, and any other is
// unicap.ErrTaskFailed.
func failureSentinel(failure solveError) error {
	if sentinel, overridden := solveFailureReasonRules[failure.Reason]; overridden {
		return sentinel
	}

	if sentinel, documented := solveFailureRules[failure.Code]; documented {
		return sentinel
	}

	if failure.Retryable {
		return unicap.ErrUnsolvable
	}

	return unicap.ErrTaskFailed
}

// unexplained reports a reply NoneCap did not explain with an error code. It
// is a plain error, which the poller treats as transient, for a server error
// or a rate limit, and a final provider error named after the HTTP status
// otherwise.
func (c *Client) unexplained(request httpapi.Request, response *httpapi.Response) error {
	return c.api.StatusError(request.Method+" "+request.Path, response.StatusCode)
}

// taskStatus converts the status of a solve that did not fail.
func taskStatus(status string) unicap.TaskStatus {
	switch status {
	case statusSolved:
		return unicap.TaskStatusReady
	case statusPending:
		return unicap.TaskStatusPending
	default:
		return unicap.TaskStatusProcessing
	}
}
