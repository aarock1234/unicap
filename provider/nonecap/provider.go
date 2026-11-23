package nonecap

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

// CreateTask submits task as a NoneCap solve.
func (c *Client) CreateTask(ctx context.Context, task unicap.Task) (*unicap.Result, error) {
	body, err := newSolveRequest(task)
	if err != nil {
		return nil, fmt.Errorf("%s: mapping task: %w", Name, err)
	}

	return c.solveResult(ctx, httpapi.Request{
		Method: http.MethodPost,
		Path:   solvesPath,
		Body:   body,
	})
}

// GetTaskResult returns the current state of a submitted solve.
func (c *Client) GetTaskResult(ctx context.Context, taskID string) (*unicap.Result, error) {
	if err := c.checkTaskID(taskID); err != nil {
		return nil, err
	}

	return c.solveResult(ctx, httpapi.Request{
		Method: http.MethodGet,
		Path:   solvePath(taskID),
	})
}

// Supports reports whether NoneCap solves tasks of the given type. Raw tasks
// are always accepted.
func (c *Client) Supports(taskType unicap.TaskType) bool {
	if taskType == unicap.TaskTypeRaw {
		return true
	}

	_, ok := mappers[taskType]

	return ok
}

// RecommendedPolling returns the polling pace NoneCap's Go guide suggests.
func (c *Client) RecommendedPolling() unicap.Polling {
	return polling
}

// newSolveRequest maps task to the create-solve body. A raw task passes
// through as its Params with the "type" field set to its TaskType; every
// other task is mapped by its registered mapper.
func newSolveRequest(task unicap.Task) (any, error) {
	if raw, ok := task.(*tasks.RawTask); ok {
		return raw.Payload(), nil
	}

	mapper, ok := mappers[task.Type()]
	if !ok {
		return nil, fmt.Errorf("%s: %w", task.Type(), unicap.ErrUnsupportedTask)
	}

	return mapper(task)
}

// solveResult runs the three steps every solve exchange shares: it sends a
// request that NoneCap answers with a solve object, classifies the reply, and
// builds the result of a solve that did not fail. The reply body is the
// solve object as received, which the result keeps as the raw solution.
func (c *Client) solveResult(ctx context.Context, request httpapi.Request) (*unicap.Result, error) {
	response, err := c.api.Do(ctx, request)
	if err != nil {
		return nil, err
	}

	solve, err := c.classifySolve(request, response)
	if err != nil {
		return nil, err
	}

	return newResult(solve, response.Body), nil
}

// checkTaskID rejects an empty task ID before any request is sent, because
// its path would name the solves collection instead of one solve. The error
// is a final provider error wrapping unicap.ErrTaskNotFound.
func (c *Client) checkTaskID(taskID string) error {
	if taskID != "" {
		return nil
	}

	return c.api.Error("INVALID_TASK_ID", "task id is empty", unicap.ErrTaskNotFound)
}
