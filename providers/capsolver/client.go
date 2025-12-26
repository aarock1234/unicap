package capsolver

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/aarock1234/unicap/unicap"
)

// capsolverClient implements the Provider interface
type capsolverClient struct {
	apiKey string
	client *unicap.BaseHTTPClient
	errors *unicap.ErrorMapper
}

// ProviderOption configures a provider
type ProviderOption func(*capsolverClient)

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client *http.Client) ProviderOption {
	return func(c *capsolverClient) {
		c.client.HTTPClient = client
	}
}

// WithBaseURL sets a custom base URL (for testing)
func WithBaseURL(url string) ProviderOption {
	return func(c *capsolverClient) {
		c.client.BaseURL = url
	}
}

// WithLogger sets a custom logger
func WithLogger(logger *slog.Logger) ProviderOption {
	return func(c *capsolverClient) {
		c.client.Logger = logger
	}
}

// NewCapSolverProvider creates a CapSolver provider
func NewCapSolverProvider(apiKey string, opts ...ProviderOption) (unicap.Provider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("api key: %w", unicap.ErrInvalidAPIKey)
	}

	c := &capsolverClient{
		apiKey: apiKey,
		client: &unicap.BaseHTTPClient{
			HTTPClient: &http.Client{Timeout: 30 * time.Second},
			Logger:     slog.New(slog.NewTextHandler(io.Discard, nil)),
			BaseURL:    "https://api.capsolver.com",
		},
		errors: unicap.StandardErrorMapper(
			"capsolver",
			[]string{"ERROR_KEY_INVALID", "ERROR_KEY_DOES_NOT_EXIST"},
			[]string{"ERROR_ZERO_BALANCE", "ERROR_NO_SLOT_AVAILABLE"},
			[]string{"ERROR_TASK_NOT_FOUND"},
			[]string{"ERROR_INVALID_TASK_DATA"},
		),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

func (c *capsolverClient) CreateTask(ctx context.Context, task unicap.Task) (string, error) {
	capsolverTask, err := mapToCapSolverTask(task)
	if err != nil {
		return "", fmt.Errorf("mapping task: %w", err)
	}

	req := createTaskRequest{
		ClientKey: c.apiKey,
		Task:      capsolverTask,
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

func (c *capsolverClient) GetTaskResult(ctx context.Context, taskID string) (*unicap.TaskResult, error) {
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
				Provider: "capsolver",
			},
		}, nil
	}

	return &unicap.TaskResult{
		Status:   mapStatus(resp.Status),
		Solution: mapSolution(resp.Solution),
	}, nil
}

func (c *capsolverClient) Name() string {
	return "capsolver"
}

func mapStatus(status string) unicap.TaskStatus {
	switch status {
	case "processing":
		return unicap.TaskStatusProcessing
	case "ready":
		return unicap.TaskStatusReady
	case "failed":
		return unicap.TaskStatusFailed
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
