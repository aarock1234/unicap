// Package solverapi implements the shared Anti-Captcha-style HTTP protocol
// (createTask / getTaskResult) spoken by the built-in providers. Providers
// supply only a base URL, an error mapper, and a task mapper; this package
// owns the request/response flow, status and solution decoding, and the
// raw-task passthrough.
package solverapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/aarock1234/unicap"
	"github.com/aarock1234/unicap/tasks"
)

var _ unicap.Provider = (*Client)(nil)

// TaskMapper converts a universal task into a provider-specific task payload.
type TaskMapper func(unicap.Task) (any, error)

// Client is a provider client that speaks the createTask / getTaskResult
// protocol.
type Client struct {
	http    *http.Client
	logger  *slog.Logger
	baseURL string
	apiKey  string
	name    string
	errors  *ErrorMapper
	mapTask TaskMapper
}

// Option configures a Client.
type Option func(*Client)

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(h *http.Client) Option {
	return func(c *Client) {
		if h != nil {
			c.http = h
		}
	}
}

// WithLogger sets a custom logger.
func WithLogger(l *slog.Logger) Option {
	return func(c *Client) {
		if l != nil {
			c.logger = l
		}
	}
}

// WithBaseURL sets a custom base URL. Intended for testing.
func WithBaseURL(u string) Option {
	return func(c *Client) {
		if u != "" {
			c.baseURL = u
		}
	}
}

// New creates a Client for the named provider.
func New(name, baseURL, apiKey string, mapper TaskMapper, errs *ErrorMapper, opts ...Option) *Client {
	c := &Client{
		http:    &http.Client{Timeout: 30 * time.Second},
		logger:  slog.New(slog.NewTextHandler(io.Discard, nil)),
		baseURL: baseURL,
		apiKey:  apiKey,
		name:    name,
		errors:  errs,
		mapTask: mapper,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// CreateTask submits a captcha task and returns the provider task ID.
func (c *Client) CreateTask(ctx context.Context, task unicap.Task) (string, error) {
	body, err := c.buildTask(task)
	if err != nil {
		return "", fmt.Errorf("mapping task: %w", err)
	}

	req := createTaskRequest{
		ClientKey: c.apiKey,
		Task:      body,
	}

	var resp createTaskResponse
	if err := c.doJSON(ctx, "/createTask", req, &resp); err != nil {
		return "", err
	}

	if resp.ErrorID != 0 {
		return "", c.errors.Error(resp.ErrorCode, resp.ErrorDescription)
	}

	return resp.TaskID.String(), nil
}

// GetTaskResult retrieves the result for the given provider task ID.
func (c *Client) GetTaskResult(ctx context.Context, taskID string) (*unicap.TaskResult, error) {
	req := getTaskResultRequest{
		ClientKey: c.apiKey,
		TaskID:    taskID,
	}

	var resp getTaskResultResponse
	if err := c.doJSON(ctx, "/getTaskResult", req, &resp); err != nil {
		return nil, err
	}

	if resp.ErrorID != 0 {
		return &unicap.TaskResult{
			Status: unicap.TaskStatusFailed,
			Error:  c.errors.Error(resp.ErrorCode, resp.ErrorDescription),
		}, nil
	}

	return &unicap.TaskResult{
		Status:   mapStatus(resp.Status),
		Solution: mapSolution(resp.Solution),
	}, nil
}

// Name returns the provider identifier.
func (c *Client) Name() string {
	return c.name
}

// buildTask maps a task to its provider payload, passing raw tasks through
// unchanged.
func (c *Client) buildTask(task unicap.Task) (any, error) {
	if raw, ok := task.(*tasks.RawTask); ok {
		return raw.Payload(), nil
	}

	return c.mapTask(task)
}

func (c *Client) doJSON(ctx context.Context, endpoint string, reqBody, respBody any) error {
	data, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+endpoint, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// The request body carries the client API key, so it is intentionally not
	// logged here; the endpoint alone is sufficient for debugging.
	c.logger.DebugContext(ctx, "sending request",
		slog.String("endpoint", endpoint),
	)

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}

	c.logger.DebugContext(ctx, "received response",
		slog.Int("status_code", resp.StatusCode),
		slog.String("body", string(body)),
	)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.Unmarshal(body, respBody); err != nil {
		return fmt.Errorf("unmarshaling response: %w", err)
	}

	return nil
}

type createTaskRequest struct {
	ClientKey string `json:"clientKey"`
	Task      any    `json:"task"`
}

type createTaskResponse struct {
	ErrorID          int    `json:"errorId"`
	ErrorCode        string `json:"errorCode,omitempty"`
	ErrorDescription string `json:"errorDescription,omitempty"`
	TaskID           TaskID `json:"taskId,omitempty"`
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
