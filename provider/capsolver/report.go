package capsolver

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/internal/solverapi"
)

type feedbackRequest struct {
	solverapi.Auth

	Solved bool           `json:"solved"`
	Task   any            `json:"task"`
	Result feedbackResult `json:"result"`
}

// feedbackResult mirrors the getTaskResult response the feedback refers to.
type feedbackResult struct {
	ErrorID  int             `json:"errorId"`
	Status   string          `json:"status"`
	Solution json.RawMessage `json:"solution"`
}

type feedbackResponse struct {
	solverapi.Envelope
}

// Report sends CapSolver feedback on whether result's solution worked.
// CapSolver needs the original task, so result must come from a
// unicap.Client call that set Result.Task.
func (c *Client) Report(ctx context.Context, result *unicap.Result, correct bool) error {
	if result.Task == nil {
		return fmt.Errorf("capsolver feedback needs the original task: %w", unicap.ErrInvalidTask)
	}

	payload, err := c.api.TaskPayload(result.Task)
	if err != nil {
		return fmt.Errorf("mapping task for feedback: %w", err)
	}

	request := &feedbackRequest{
		Solved: correct,
		Task:   payload,
		Result: feedbackResult{
			ErrorID:  0,
			Status:   "ready",
			Solution: result.Solution.Raw,
		},
	}

	var response feedbackResponse

	return c.api.Call(ctx, "feedbackTask", request, &response)
}
