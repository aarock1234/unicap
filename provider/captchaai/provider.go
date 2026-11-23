package captchaai

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/internal/httpapi"
)

// Name returns the provider identifier.
func (c *Client) Name() string {
	return c.api.Name()
}

// CreateTask submits task. CaptchaAI answers with a task ID only, so the
// result is never ready yet. Its task ID names the shape of the solution as
// well as the task, and reads "<shape>:<task id>".
func (c *Client) CreateTask(ctx context.Context, task unicap.Task) (*unicap.Result, error) {
	submitted, err := newSubmission(task)
	if err != nil {
		return nil, fmt.Errorf("%s: mapping task: %w", Name, err)
	}

	request, err := c.request(submitPath, submitted.fields)
	if err != nil {
		return nil, err
	}

	response, err := c.api.Do(ctx, request)
	if err != nil {
		return nil, err
	}

	vendorTaskID, err := c.classifySubmit(request, response)
	if err != nil {
		return nil, err
	}

	return newResult(submitted.ref(vendorTaskID), unicap.TaskStatusPending), nil
}

// GetTaskResult returns the current state of a submitted task. taskID is the
// task ID CreateTask returned, which reads "<shape>:<task id>".
func (c *Client) GetTaskResult(ctx context.Context, taskID string) (*unicap.Result, error) {
	ref, err := parseTaskRef(taskID)
	if err != nil {
		return nil, err
	}

	request, err := c.request(resultPath, ref.fields())
	if err != nil {
		return nil, err
	}

	response, err := c.api.Do(ctx, request)
	if err != nil {
		return nil, err
	}

	answer, solved, err := c.classifyResult(request, response)
	if err != nil {
		return nil, err
	}

	if !solved {
		return newResult(ref, unicap.TaskStatusProcessing), nil
	}

	return newReadyResult(ref, answer, response.Body), nil
}

// Supports reports whether CaptchaAI solves tasks of the given type. Raw
// tasks are always accepted.
func (c *Client) Supports(taskType unicap.TaskType) bool {
	if taskType == unicap.TaskTypeRaw {
		return true
	}

	_, ok := mappers[taskType]

	return ok
}

// RecommendedPolling returns the polling pace unicap recommends for
// CaptchaAI, which follows CaptchaAI's guidance.
func (c *Client) RecommendedPolling() unicap.Polling {
	return polling
}

// request builds the POST that sends fields to path as a multipart form. The
// API key goes first and the switch to JSON replies last, so both are part
// of every request and the key stays out of the URL.
func (c *Client) request(path string, fields form) (httpapi.Request, error) {
	complete := form{}.
		with(fieldKey, c.apiKey).
		extend(fields).
		with(fieldJSON, flagSet)

	// The encoding error is left out of the message so that no field value,
	// such as the API key or a proxy password, can reach it.
	body, err := encodeForm(complete)
	if err != nil {
		return httpapi.Request{}, c.api.Error(
			"INVALID_REQUEST",
			http.MethodPost+" "+path+": encoding request form",
			unicap.ErrInvalidTask,
		)
	}

	return httpapi.Request{
		Method:  http.MethodPost,
		Path:    path,
		Encoded: body,
	}, nil
}
