package twocaptcha

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/aarock1234/unicap/pkg/upicap"
)

// twocaptchaClient implements the Provider interface
type twocaptchaClient struct {
	apiKey string
	client *upicap.BaseHTTPClient
	errors *upicap.ErrorMapper
}

// ProviderOption configures a provider
type ProviderOption func(*twocaptchaClient)

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client *http.Client) ProviderOption {
	return func(c *twocaptchaClient) {
		c.client.HTTPClient = client
	}
}

// WithBaseURL sets a custom base URL (for testing)
func WithBaseURL(url string) ProviderOption {
	return func(c *twocaptchaClient) {
		c.client.BaseURL = url
	}
}

// WithLogger sets a custom logger
func WithLogger(logger *slog.Logger) ProviderOption {
	return func(c *twocaptchaClient) {
		c.client.Logger = logger
	}
}

// NewTwoCaptchaProvider creates a 2Captcha provider
func NewTwoCaptchaProvider(apiKey string, opts ...ProviderOption) (upicap.Provider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("api key: %w", upicap.ErrInvalidAPIKey)
	}

	c := &twocaptchaClient{
		apiKey: apiKey,
		client: &upicap.BaseHTTPClient{
			HTTPClient: &http.Client{Timeout: 30 * time.Second},
			Logger:     slog.New(slog.NewTextHandler(io.Discard, nil)),
			BaseURL:    "https://api.2captcha.com",
		},
		errors: upicap.StandardErrorMapper(
			"2captcha",
			[]string{"ERROR_KEY_DOES_NOT_EXIST", "ERROR_WRONG_USER_KEY"},
			[]string{"ERROR_ZERO_BALANCE"},
			[]string{"ERROR_TASK_ABSENT"},
			[]string{"ERROR_WRONG_TASK_DATA"},
		),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

func (c *twocaptchaClient) CreateTask(ctx context.Context, task upicap.Task) (string, error) {
	twocaptchaTask, err := mapToTwoCaptchaTask(task)
	if err != nil {
		return "", fmt.Errorf("mapping task: %w", err)
	}

	req := createTaskRequest{
		ClientKey: c.apiKey,
		Task:      twocaptchaTask,
	}

	var resp createTaskResponse
	if err := c.client.DoJSON(ctx, "/createTask", req, &resp); err != nil {
		return "", err
	}

	if resp.ErrorID != 0 {
		return "", c.errors.MapError(resp.ErrorCode, resp.ErrorDescription)
	}

	c.client.Logger.InfoContext(ctx, "task created",
		slog.String("task_id", resp.TaskID),
		slog.String("task_type", string(task.Type())),
	)

	return resp.TaskID, nil
}

func (c *twocaptchaClient) GetTaskResult(ctx context.Context, taskID string) (*upicap.TaskResult, error) {
	req := getTaskResultRequest{
		ClientKey: c.apiKey,
		TaskID:    taskID,
	}

	var resp getTaskResultResponse
	if err := c.client.DoJSON(ctx, "/getTaskResult", req, &resp); err != nil {
		return nil, err
	}

	if resp.ErrorID != 0 {
		return &upicap.TaskResult{
			Status: upicap.TaskStatusFailed,
			Error: &upicap.Error{
				Code:     resp.ErrorCode,
				Message:  resp.ErrorDescription,
				Provider: "2captcha",
			},
		}, nil
	}

	return &upicap.TaskResult{
		Status:   mapStatus(resp.Status),
		Solution: mapSolution(resp.Solution),
	}, nil
}

func (c *twocaptchaClient) Name() string {
	return "2captcha"
}

func mapStatus(status string) upicap.TaskStatus {
	switch status {
	case "processing":
		return upicap.TaskStatusProcessing
	case "ready":
		return upicap.TaskStatusReady
	default:
		return upicap.TaskStatusPending
	}
}

func mapSolution(solution map[string]any) upicap.Solution {
	sol := upicap.Solution{
		Extra: solution,
	}

	if token, ok := solution["gRecaptchaResponse"].(string); ok {
		sol.Token = token
	} else if token, ok := solution["token"].(string); ok {
		sol.Token = token
	}

	if text, ok := solution["text"].(string); ok {
		sol.Text = text
	}

	return sol
}

// Request/Response types
type createTaskRequest struct {
	ClientKey string `json:"clientKey"`
	Task      any    `json:"task"`
}

type createTaskResponse struct {
	ErrorID          int    `json:"errorId"`
	ErrorCode        string `json:"errorCode,omitempty"`
	ErrorDescription string `json:"errorDescription,omitempty"`
	TaskID           string `json:"taskId,omitempty"`
}

type getTaskResultRequest struct {
	ClientKey string `json:"clientKey"`
	TaskID    string `json:"taskId"`
}

type getTaskResultResponse struct {
	ErrorID          int            `json:"errorId"`
	ErrorCode        string         `json:"errorCode,omitempty"`
	ErrorDescription string         `json:"errorDescription,omitempty"`
	Status           string         `json:"status,omitempty"`
	Solution         map[string]any `json:"solution,omitempty"`
}
