// Package anticaptcha provides the Anti-Captcha implementation of the unicap
// provider interface.
package anticaptcha

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/aarock1234/unicap"
	"github.com/aarock1234/unicap/internal/providerutil"
)

var _ unicap.Provider = (*client)(nil)

type client struct {
	apiKey string
	client *providerutil.HTTPClient
	errors *providerutil.ErrorMapper
}

// Option configures the provider.
type Option func(*client)

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *client) {
		c.client.HTTPClient = httpClient
	}
}

// WithBaseURL sets a custom base URL (for testing)
func WithBaseURL(url string) Option {
	return func(c *client) {
		c.client.BaseURL = url
	}
}

// WithLogger sets a custom logger
func WithLogger(logger *slog.Logger) Option {
	return func(c *client) {
		c.client.Logger = logger
	}
}

// New creates an Anti-Captcha provider.
func New(apiKey string, opts ...Option) (unicap.Provider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("api key: %w", unicap.ErrInvalidAPIKey)
	}

	c := &client{
		apiKey: apiKey,
		client: &providerutil.HTTPClient{
			HTTPClient: &http.Client{Timeout: 30 * time.Second},
			Logger:     slog.New(slog.NewTextHandler(io.Discard, nil)),
			BaseURL:    "https://api.anti-captcha.com",
		},
		errors: providerutil.StandardErrorMapper(
			"anticaptcha",
			[]string{"ERROR_KEY_DOES_NOT_EXIST", "ERROR_WRONG_USER_KEY"},
			[]string{"ERROR_ZERO_BALANCE", "ERROR_NO_SLOT_AVAILABLE"},
			[]string{"ERROR_TASK_ABSENT"},
			[]string{"ERROR_WRONG_TASK_DATA"},
		),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

func (c *client) CreateTask(ctx context.Context, task unicap.Task) (string, error) {
	anticaptchaTask, err := mapToAntiCaptchaTask(task)
	if err != nil {
		return "", fmt.Errorf("mapping task: %w", err)
	}

	req := createTaskRequest{
		ClientKey: c.apiKey,
		Task:      anticaptchaTask,
	}

	var resp createTaskResponse
	if err := c.client.DoJSON(ctx, "/createTask", req, &resp); err != nil {
		return "", err
	}

	if resp.ErrorID != 0 {
		return "", c.errors.MapError(resp.ErrorCode, resp.ErrorDescription)
	}

	c.client.Logger.InfoContext(ctx, "task created",
		slog.Int("task_id", resp.TaskID),
		slog.String("task_type", string(task.Type())),
	)

	return fmt.Sprintf("%d", resp.TaskID), nil
}

func (c *client) GetTaskResult(ctx context.Context, taskID string) (*unicap.TaskResult, error) {
	req := getTaskResultRequest{
		ClientKey: c.apiKey,
		TaskID:    taskID,
	}

	var resp getTaskResultResponse
	if err := c.client.DoJSON(ctx, "/getTaskResult", req, &resp); err != nil {
		return nil, err
	}

	if resp.ErrorID != 0 {
		return &unicap.TaskResult{
			Status: unicap.TaskStatusFailed,
			Error: &unicap.Error{
				Code:     resp.ErrorCode,
				Message:  resp.ErrorDescription,
				Provider: "anticaptcha",
			},
		}, nil
	}

	return &unicap.TaskResult{
		Status:   mapStatus(resp.Status),
		Solution: mapSolution(resp.Solution),
	}, nil
}

func (c *client) Name() string {
	return "anticaptcha"
}

func mapStatus(status string) unicap.TaskStatus {
	switch status {
	case "processing":
		return unicap.TaskStatusProcessing
	case "ready":
		return unicap.TaskStatusReady
	default:
		return unicap.TaskStatusPending
	}
}

func mapSolution(solution map[string]any) unicap.Solution {
	sol := unicap.Solution{
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
	TaskID           int    `json:"taskId,omitempty"`
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
