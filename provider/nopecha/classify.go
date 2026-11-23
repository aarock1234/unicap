package nopecha

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aarock1234/unicap/v3/internal/httpapi"
)

// jsonNull is the JSON value of a reply field that carries nothing.
var jsonNull = []byte("null")

// classifySubmit sorts the reply to a submitted job and returns the ID of
// the job NopeCHA created. A rejected request and a reply without a job ID
// are errors.
func (c *Client) classifySubmit(request httpapi.Request, response *httpapi.Response) (string, error) {
	if response.StatusCode != http.StatusOK {
		return "", c.requestError(request, response)
	}

	var reply submitReply
	if err := json.Unmarshal(response.Body, &reply); err != nil {
		return "", c.unexplained(request, response)
	}

	if reply.Data == "" {
		return "", c.unexplained(request, response)
	}

	return reply.Data, nil
}

// classifyFetch sorts the reply to a fetch of a job's result. It returns the
// solution data and true once the job is finished, and false while NopeCHA
// reports the job as incomplete. A rejected request and a reply without data
// are errors.
func (c *Client) classifyFetch(request httpapi.Request, response *httpapi.Response) (json.RawMessage, bool, error) {
	if isIncomplete(response) {
		return nil, false, nil
	}

	if response.StatusCode != http.StatusOK {
		return nil, false, c.requestError(request, response)
	}

	var reply resultReply
	if err := json.Unmarshal(response.Body, &reply); err != nil {
		return nil, false, c.unexplained(request, response)
	}

	if len(reply.Data) == 0 || bytes.Equal(reply.Data, jsonNull) {
		return nil, false, c.unexplained(request, response)
	}

	return reply.Data, true, nil
}

// isIncomplete reports whether a reply says the job is not finished: HTTP
// status 409 with the incomplete-job code in an error body.
func isIncomplete(response *httpapi.Response) bool {
	if response.StatusCode != http.StatusConflict {
		return false
	}

	var body errorBody
	if err := json.Unmarshal(response.Body, &body); err != nil {
		return false
	}

	return body.Code == incompleteJobCode
}

// requestError returns the error a reply without a success status stands
// for. An error body with a documented code is a final provider error
// classified by errorRules; anything else is unexplained.
func (c *Client) requestError(request httpapi.Request, response *httpapi.Response) error {
	var body errorBody
	if err := json.Unmarshal(response.Body, &body); err != nil {
		return c.unexplained(request, response)
	}

	sentinel, documented := errorRules[body.Code]
	if !documented {
		return c.unexplained(request, response)
	}

	return c.api.Error(strconv.Itoa(body.Code), errorMessage(body), sentinel)
}

// errorMessage returns NopeCHA's description of an error: its message,
// followed by its type when the error has one.
func errorMessage(body errorBody) string {
	if body.Type == "" {
		return body.Message
	}

	return body.Message + ": " + body.Type
}

// unexplained reports a reply NopeCHA did not explain with a documented
// error code. It is a plain error, which the poller treats as transient, for
// a server error or a rate limit, and a final provider error named after the
// HTTP status otherwise.
func (c *Client) unexplained(request httpapi.Request, response *httpapi.Response) error {
	return c.api.StatusError(request.Method+" "+request.Path, response.StatusCode)
}
