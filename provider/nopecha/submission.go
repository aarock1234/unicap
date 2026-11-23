package nopecha

import (
	"fmt"
	"net/http"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/internal/httpapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

// submission is a task mapped for the endpoint that takes it, whatever the
// type of its body.
type submission interface {
	// request returns the request that submits the job.
	request() httpapi.Request

	// job returns the reference of the job the endpoint created as id.
	job(id string) jobRef
}

// endpointBody is the submission of a body of type B to one endpoint.
type endpointBody[B any] struct {
	endpoint string
	body     B
}

// request returns the request that submits the job.
func (s endpointBody[B]) request() httpapi.Request {
	return httpapi.Request{
		Method: http.MethodPost,
		Path:   endpointPath(s.endpoint),
		Body:   s.body,
	}
}

// job returns the reference of the job the endpoint created as id.
func (s endpointBody[B]) job(id string) jobRef {
	return jobRef{
		endpoint: s.endpoint,
		id:       id,
	}
}

// newSubmission maps task to its submission. A raw task names its own
// endpoint; every other task is mapped by its registered mapper.
func newSubmission(task unicap.Task) (submission, error) {
	if raw, ok := task.(*tasks.RawTask); ok {
		return newRawSubmission(raw)
	}

	mapper, ok := mappers[task.Type()]
	if !ok {
		return nil, fmt.Errorf("%s: %w", task.Type(), unicap.ErrUnsupportedTask)
	}

	mapped, err := mapper(task)
	if err != nil {
		return nil, err
	}

	mappedSubmission, ok := mapped.(submission)
	if !ok {
		return nil, fmt.Errorf("%s mapper returned unexpected type %T: %w", task.Type(), mapped, unicap.ErrUnsupportedTask)
	}

	return mappedSubmission, nil
}

// newRawSubmission passes a raw task through: its TaskType is the endpoint,
// such as "recognition/funcaptcha", and its Params are the request body. A
// TaskType that is not an endpoint is rejected, because it would address
// some other path.
func newRawSubmission(raw *tasks.RawTask) (submission, error) {
	if !endpointPattern.MatchString(raw.TaskType) {
		return nil, fmt.Errorf("raw task type must be an endpoint such as %q: %w", endpointHCaptcha, unicap.ErrInvalidTask)
	}

	return endpointBody[map[string]any]{
		endpoint: raw.TaskType,
		body:     raw.Params,
	}, nil
}
