package capmonster

import (
	"context"
	"fmt"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/internal/solverapi"
)

// reportRequest names the task an incorrect-answer report refers to.
type reportRequest struct {
	solverapi.Auth

	TaskID any `json:"taskId"`
}

// reportResponse is the reply to a report. Status is "success" when the
// report was accepted.
type reportResponse struct {
	solverapi.Envelope

	Status string `json:"status"`
}

// Report tells CapMonster Cloud that result's solution was rejected.
// CapMonster Cloud accepts incorrect reports only, for image-to-text,
// reCAPTCHA, GeeTest and Turnstile tasks; result must carry its task to pick
// the endpoint.
func (c *Client) Report(ctx context.Context, result *unicap.Result, correct bool) error {
	if result == nil {
		return fmt.Errorf("report needs a result: %w", unicap.ErrInvalidTask)
	}

	if result.Task == nil {
		return fmt.Errorf("capmonster reports need the original task: %w", unicap.ErrInvalidTask)
	}

	method, err := reportMethod(result.Task.Type(), correct)
	if err != nil {
		return err
	}

	taskID, err := c.api.EncodeTaskID(result.TaskID)
	if err != nil {
		return err
	}

	request := &reportRequest{
		TaskID: taskID,
	}

	var response reportResponse

	return c.api.Call(ctx, method, request, &response)
}

// reportMethod returns the report endpoint for a task type and verdict.
func reportMethod(taskType unicap.TaskType, correct bool) (string, error) {
	if correct {
		return "", fmt.Errorf("correct reports: %w", unicap.ErrUnsupportedOperation)
	}

	switch taskType {
	case unicap.TaskTypeImageToText:
		return "reportIncorrectImageCaptcha", nil

	case unicap.TaskTypeReCaptchaV2,
		unicap.TaskTypeReCaptchaV2Enterprise,
		unicap.TaskTypeReCaptchaV3,
		unicap.TaskTypeReCaptchaV3Enterprise,
		unicap.TaskTypeGeeTest,
		unicap.TaskTypeGeeTestV4,
		unicap.TaskTypeTurnstile:
		return "reportIncorrectTokenCaptcha", nil

	default:
		return "", fmt.Errorf("reports for %s tasks: %w", taskType, unicap.ErrUnsupportedOperation)
	}
}
