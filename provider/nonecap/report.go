package nonecap

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/internal/httpapi"
)

// Report sends NoneCap feedback on whether result's token worked. NoneCap
// accepts feedback for solved solves, the first report within 30 days of the
// solve, and keeps the last report for a solve.
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
			Outcome: outcome(correct),
		},
	}

	response, err := c.api.Do(ctx, request)
	if err != nil {
		return err
	}

	return c.requestError(request, response)
}

// outcome names the feedback outcome for a token that did or did not work.
func outcome(correct bool) string {
	if correct {
		return outcomeAccepted
	}

	return outcomeRejected
}
