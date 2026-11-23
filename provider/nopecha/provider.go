package nopecha

import (
	"context"
	"fmt"

	"github.com/aarock1234/unicap/v3"
)

// Name returns the provider identifier.
func (c *Client) Name() string {
	return c.api.Name()
}

// CreateTask submits task as a NopeCHA job. NopeCHA answers with a job ID
// only, so the result is never ready yet. Its task ID names the endpoint as
// well as the job, and reads "<endpoint>:<job id>".
func (c *Client) CreateTask(ctx context.Context, task unicap.Task) (*unicap.Result, error) {
	submission, err := newSubmission(task)
	if err != nil {
		return nil, fmt.Errorf("%s: mapping task: %w", Name, err)
	}

	request := submission.request()

	response, err := c.api.Do(ctx, request)
	if err != nil {
		return nil, err
	}

	jobID, err := c.classifySubmit(request, response)
	if err != nil {
		return nil, err
	}

	return newResult(submission.job(jobID), unicap.TaskStatusPending), nil
}

// GetTaskResult returns the current state of a submitted job. taskID is the
// task ID CreateTask returned, which reads "<endpoint>:<job id>". NopeCHA
// answers only to the API key that submitted the job.
func (c *Client) GetTaskResult(ctx context.Context, taskID string) (*unicap.Result, error) {
	job, err := parseJobRef(taskID)
	if err != nil {
		return nil, err
	}

	request := job.request()

	response, err := c.api.Do(ctx, request)
	if err != nil {
		return nil, err
	}

	data, finished, err := c.classifyFetch(request, response)
	if err != nil {
		return nil, err
	}

	if !finished {
		return newResult(job, unicap.TaskStatusProcessing), nil
	}

	return newReadyResult(job, data), nil
}

// Supports reports whether NopeCHA solves tasks of the given type. Raw tasks
// are always accepted.
func (c *Client) Supports(taskType unicap.TaskType) bool {
	if taskType == unicap.TaskTypeRaw {
		return true
	}

	_, ok := mappers[taskType]

	return ok
}

// RecommendedPolling returns the polling pace unicap recommends for NopeCHA,
// which follows NopeCHA's advice.
func (c *Client) RecommendedPolling() unicap.Polling {
	return polling
}
