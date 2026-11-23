package noptcha

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

// CreateTask submits task. Task types Noptcha solves synchronously come back
// already ready.
func (c *Client) CreateTask(ctx context.Context, task unicap.Task) (*unicap.Result, error) {
	body, err := newTaskRequest(task)
	if err != nil {
		return nil, fmt.Errorf("%s: mapping task: %w", Name, err)
	}

	return c.taskResult(ctx, httpapi.Request{
		Method: http.MethodPost,
		Path:   tasksPath,
		Body:   body,
	})
}

// GetTaskResult returns the current state of a submitted task. Noptcha keeps
// a result for 10 minutes after the task completes.
func (c *Client) GetTaskResult(ctx context.Context, taskID string) (*unicap.Result, error) {
	if err := c.checkTaskID(taskID); err != nil {
		return nil, err
	}

	return c.taskResult(ctx, httpapi.Request{
		Method: http.MethodGet,
		Path:   taskPath(taskID),
	})
}

// Supports reports whether Noptcha solves tasks of the given type. Raw tasks
// are always accepted.
func (c *Client) Supports(taskType unicap.TaskType) bool {
	if taskType == unicap.TaskTypeRaw {
		return true
	}

	_, ok := mappers[taskType]

	return ok
}

// RecommendedPolling returns the polling pace unicap recommends for Noptcha.
func (c *Client) RecommendedPolling() unicap.Polling {
	return polling
}

// newTaskRequest maps task to the create-task body. A raw task passes
// through with its TaskType as the native type and its Params as the params
// object; every other task is mapped by its registered mapper.
func newTaskRequest(task unicap.Task) (any, error) {
	if raw, ok := task.(*tasks.RawTask); ok {
		return taskRequest[map[string]any]{
			Type:   raw.TaskType,
			Params: raw.Params,
		}, nil
	}

	mapper, ok := mappers[task.Type()]
	if !ok {
		return nil, fmt.Errorf("%s: %w", task.Type(), unicap.ErrUnsupportedTask)
	}

	return mapper(task)
}

// taskResult runs the three steps every task exchange shares: it sends a
// request that Noptcha answers with a task object, classifies the reply, and
// builds the result of a task that did not fail.
func (c *Client) taskResult(ctx context.Context, request httpapi.Request) (*unicap.Result, error) {
	response, err := c.api.Do(ctx, request)
	if err != nil {
		return nil, err
	}

	task, err := c.classifyTask(request, response)
	if err != nil {
		return nil, err
	}

	return newResult(task), nil
}

// checkTaskID rejects a task ID Noptcha could not have issued before any
// request is sent. The error is a final provider error wrapping
// unicap.ErrTaskNotFound.
func (c *Client) checkTaskID(taskID string) error {
	if taskIDPattern.MatchString(taskID) {
		return nil
	}

	return c.api.Error("INVALID_TASK_ID", "task id must be 1 to 16 decimal digits", unicap.ErrTaskNotFound)
}
