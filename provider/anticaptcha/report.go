package anticaptcha

import (
	"context"
	"fmt"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/internal/solverapi"
)

type taskRequest struct {
	solverapi.Auth

	TaskID any `json:"taskId"`
}

type pushVariableRequest struct {
	solverapi.Auth

	TaskID any    `json:"taskId"`
	Name   string `json:"name"`
	Value  any    `json:"value"`
}

type statusResponse struct {
	solverapi.Envelope

	Status string `json:"status"`
}

// Report tells Anti-Captcha whether result's solution was accepted.
// Anti-Captcha accepts reports for reCAPTCHA tasks, and incorrect reports
// for image-to-text tasks; result must carry its task to pick the endpoint.
func (c *Client) Report(ctx context.Context, result *unicap.Result, correct bool) error {
	if result.Task == nil {
		return fmt.Errorf("anti-captcha reports need the original task: %w", unicap.ErrInvalidTask)
	}

	method, err := reportMethod(result.Task.Type(), correct)
	if err != nil {
		return err
	}

	taskID, err := c.api.EncodeTaskID(result.TaskID)
	if err != nil {
		return err
	}

	request := &taskRequest{
		TaskID: taskID,
	}

	var response statusResponse

	return c.api.Call(ctx, method, request, &response)
}

// PushVariable supplies the value of an AntiGate template variable that was
// submitted as "_WAIT_FOR_IT_".
func (c *Client) PushVariable(ctx context.Context, taskID, name string, value any) error {
	encodedID, err := c.api.EncodeTaskID(taskID)
	if err != nil {
		return err
	}

	request := &pushVariableRequest{
		TaskID: encodedID,
		Name:   name,
		Value:  value,
	}

	var response statusResponse

	return c.api.Call(ctx, "pushAntiGateVariable", request, &response)
}

// reportMethod returns the report endpoint for a task type and verdict.
func reportMethod(taskType unicap.TaskType, correct bool) (string, error) {
	switch taskType {
	case unicap.TaskTypeReCaptchaV2,
		unicap.TaskTypeReCaptchaV2Enterprise,
		unicap.TaskTypeReCaptchaV3,
		unicap.TaskTypeReCaptchaV3Enterprise:
		if correct {
			return "reportCorrectRecaptcha", nil
		}

		return "reportIncorrectRecaptcha", nil

	case unicap.TaskTypeImageToText:
		if correct {
			return "", fmt.Errorf("correct reports for image tasks: %w", unicap.ErrUnsupportedOperation)
		}

		return "reportIncorrectImageCaptcha", nil

	default:
		return "", fmt.Errorf("reports for %s tasks: %w", taskType, unicap.ErrUnsupportedOperation)
	}
}
