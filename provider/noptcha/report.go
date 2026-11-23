package noptcha

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/internal/httpapi"
)

// Report sends Noptcha feedback on whether result's solution worked. Noptcha
// accepts feedback for solved tasks within 10 minutes of completion, counts
// only the first report for a task, and does not refund the task.
func (c *Client) Report(ctx context.Context, result *unicap.Result, correct bool) error {
	if result == nil {
		return fmt.Errorf("report needs a result: %w", unicap.ErrInvalidTask)
	}

	if err := c.checkTaskID(result.TaskID); err != nil {
		return err
	}

	request := httpapi.Request{
		Method: http.MethodPost,
		Path:   feedbackPath(result.TaskID),
		Body: feedbackRequest{
			Correct: correct,
		},
	}

	response, err := c.api.Do(ctx, request)
	if err != nil {
		return err
	}

	return c.requestError(request, response)
}
