package easycaptcha

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/internal/httpapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

// Name returns the provider identifier.
func (c *Client) Name() string {
	return c.api.Name()
}

// CreateTask submits task. easycaptcha.to asks clients never to repeat this
// request after an ambiguous timeout.
func (c *Client) CreateTask(ctx context.Context, task unicap.Task) (*unicap.Result, error) {
	body, err := newTaskRequest(task)
	if err != nil {
		return nil, fmt.Errorf("%s: mapping task: %w", Name, err)
	}

	request := httpapi.Request{
		Method: http.MethodPost,
		Path:   taskCollectionPath,
		Body:   body,
	}

	return c.taskResult(ctx, request, createStatuses)
}

// GetTaskResult returns the current state of a submitted task.
// easycaptcha.to keeps a result for five minutes after the task completes.
func (c *Client) GetTaskResult(ctx context.Context, taskID string) (*unicap.Result, error) {
	if err := c.checkTaskID(taskID); err != nil {
		return nil, err
	}

	request := httpapi.Request{
		Method: http.MethodGet,
		Path:   taskPath(taskID),
	}

	return c.taskResult(ctx, request, fetchStatuses)
}

// Supports reports whether easycaptcha.to solves tasks of the given type.
// Raw tasks are always accepted.
func (c *Client) Supports(taskType unicap.TaskType) bool {
	if taskType == unicap.TaskTypeRaw {
		return true
	}

	_, ok := mappers[taskType]

	return ok
}

// RecommendedPolling returns the polling pace unicap recommends for
// easycaptcha.to.
func (c *Client) RecommendedPolling() unicap.Polling {
	return polling
}

// newTaskRequest maps task to the create-task body. A raw task passes
// through as its Params with the "type" field set to its TaskType; every
// other task is mapped by its registered mapper.
func newTaskRequest(task unicap.Task) (any, error) {
	if raw, ok := task.(*tasks.RawTask); ok {
		return raw.Payload(), nil
	}

	mapper, ok := mappers[task.Type()]
	if !ok {
		return nil, fmt.Errorf("%s: %w", task.Type(), unicap.ErrUnsupportedTask)
	}

	return mapper(task)
}

// taskResult runs the three steps every task exchange shares: it sends a
// request that easycaptcha.to answers with a task object, classifies the
// reply against the endpoint's success statuses, and builds the result of a
// task that did not fail.
func (c *Client) taskResult(ctx context.Context, request httpapi.Request, success []int) (*unicap.Result, error) {
	response, err := c.api.Do(ctx, request)
	if err != nil {
		return nil, err
	}

	task, err := c.classifyTask(request, response, success)
	if err != nil {
		return nil, err
	}

	return newResult(task), nil
}

// checkTaskID rejects an empty task ID before any request is sent, because
// its path would name the task collection instead of one task. The error is
// a final provider error wrapping unicap.ErrTaskNotFound.
func (c *Client) checkTaskID(taskID string) error {
	if taskID != "" {
		return nil
	}

	return c.api.Error("INVALID_TASK_ID", "task id is empty", unicap.ErrTaskNotFound)
}
