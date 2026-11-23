package easycaptcha

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/internal/httpapi"
)

// Report sends easycaptcha.to feedback on whether result's solution worked.
// easycaptcha.to accepts feedback for completed tasks it still retains,
// ignores an identical repeat, and rejects feedback that contradicts an
// earlier report as a conflict.
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
			Success: correct,
		},
	}

	response, err := c.api.Do(ctx, request)
	if err != nil {
		return err
	}

	return c.requestError(request, response, feedbackStatuses)
}
