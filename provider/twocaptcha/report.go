package twocaptcha

import (
	"context"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/internal/solverapi"
)

type reportRequest struct {
	solverapi.Auth

	TaskID any `json:"taskId"`
}

type reportResponse struct {
	solverapi.Envelope

	Status string `json:"status"`
}

// Report tells 2Captcha whether result's solution was accepted. Incorrect
// reports are reviewed and may be refunded.
func (c *Client) Report(ctx context.Context, result *unicap.Result, correct bool) error {
	taskID, err := c.api.EncodeTaskID(result.TaskID)
	if err != nil {
		return err
	}

	method := "reportIncorrect"
	if correct {
		method = "reportCorrect"
	}

	request := &reportRequest{
		TaskID: taskID,
	}

	var response reportResponse

	return c.api.Call(ctx, method, request, &response)
}
