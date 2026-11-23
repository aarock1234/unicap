package nopecha

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/internal/httpapi"
)

// jobRefSeparator parts the endpoint from the job ID in a unicap task ID.
const jobRefSeparator = ":"

// jobRef names a NopeCHA job. A job's result is fetched from the endpoint
// the job was submitted to, so the job ID alone does not identify it. As a
// string, which is the unicap task ID, it reads "<endpoint>:<job id>".
type jobRef struct {
	endpoint string
	id       string
}

// String returns the unicap task ID of the job.
func (j jobRef) String() string {
	return j.endpoint + jobRefSeparator + j.id
}

// request returns the request that fetches the job's result.
func (j jobRef) request() httpapi.Request {
	return httpapi.Request{
		Method: http.MethodGet,
		Path:   endpointPath(j.endpoint),
		Query: url.Values{
			jobIDParameter: []string{
				j.id,
			},
		},
	}
}

// parseJobRef reads a unicap task ID. The ID is cut at the first separator,
// and both parts must be present. The endpoint must have the form of one,
// because it becomes part of the request path. A malformed ID is a final
// provider error wrapping unicap.ErrTaskNotFound.
func parseJobRef(taskID string) (jobRef, error) {
	endpoint, id, found := strings.Cut(taskID, jobRefSeparator)
	if !found || id == "" || !endpointPattern.MatchString(endpoint) {
		return jobRef{}, &unicap.Error{
			Provider: Name,
			Code:     "INVALID_TASK_ID",
			Message:  `task id must read "<endpoint>:<job id>"`,
			Err:      unicap.ErrTaskNotFound,
		}
	}

	return jobRef{
		endpoint: endpoint,
		id:       id,
	}, nil
}
