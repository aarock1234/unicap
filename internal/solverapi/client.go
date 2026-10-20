// Package solverapi implements the createTask / getTaskResult JSON protocol
// shared by the built-in providers. A provider supplies a Config with its
// endpoint, task mappers, error rules, and polling advice; this package owns
// the HTTP exchange, envelope and error decoding, solution normalisation,
// and the raw-task passthrough.
package solverapi

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/tasks"
)

var (
	_ unicap.Provider       = (*Client)(nil)
	_ unicap.BalanceChecker = (*Client)(nil)
	_ unicap.TaskSupporter  = (*Client)(nil)
	_ unicap.PollingAdvisor = (*Client)(nil)
)

// defaultHTTPTimeout bounds a single API request.
const defaultHTTPTimeout = 30 * time.Second

// Config describes a provider that speaks the createTask / getTaskResult
// protocol.
type Config struct {
	// Name is the provider identifier reported in errors and logs.
	Name string

	// BaseURL is the API root, such as "https://api.capsolver.com".
	BaseURL string

	// APIKey is the account's client key.
	APIKey string

	// TaskIDFormat is the JSON type the provider expects for taskId.
	TaskIDFormat TaskIDFormat

	// Mappers convert unicap tasks into provider task objects.
	Mappers Mappers

	// Errors maps provider error codes to unicap sentinel errors.
	Errors ErrorRules

	// Polling is the provider's recommended polling pace.
	Polling unicap.Polling
}

// Client is a provider client for the createTask / getTaskResult protocol.
// It is safe for concurrent use.
type Client struct {
	config Config
	http   *http.Client
	logger *slog.Logger
	extras taskExtras
}

// taskExtras are optional top-level createTask fields.
type taskExtras struct {
	appID        string
	softID       int
	callbackURL  string
	languagePool string
}

// New creates a client for the provider described by config.
func New(config Config, opts ...Option) (*Client, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("%s api key is empty: %w", config.Name, unicap.ErrInvalidAPIKey)
	}

	c := &Client{
		config: config,
		http: &http.Client{
			Timeout: defaultHTTPTimeout,
		},
		logger: slog.New(slog.DiscardHandler),
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.extras.softID < 0 {
		return nil, errors.New("soft id must not be negative")
	}

	return c, nil
}

// Name returns the provider identifier.
func (c *Client) Name() string {
	return c.config.Name
}

// CreateTask submits task. Recognition tasks that the provider solves
// synchronously come back already ready.
func (c *Client) CreateTask(ctx context.Context, task unicap.Task) (*unicap.Result, error) {
	payload, err := c.taskPayload(task)
	if err != nil {
		return nil, fmt.Errorf("%s: mapping task: %w", c.config.Name, err)
	}

	request := &createTaskRequest{
		Task:         payload,
		AppID:        c.extras.appID,
		SoftID:       c.extras.softID,
		CallbackURL:  c.extras.callbackURL,
		LanguagePool: c.extras.languagePool,
	}

	var response createTaskResponse
	if err := c.Call(ctx, "createTask", request, &response); err != nil {
		return nil, err
	}

	if response.Status == statusReady {
		return &unicap.Result{
			TaskID:   response.TaskID.String(),
			Status:   unicap.TaskStatusReady,
			Solution: decodeSolution(response.Solution),
		}, nil
	}

	if response.TaskID == "" {
		return nil, c.providerError("NO_TASK_ID", "createTask returned no task id", unicap.ErrTaskFailed)
	}

	return &unicap.Result{
		TaskID: response.TaskID.String(),
		Status: unicap.TaskStatusPending,
	}, nil
}

// GetTaskResult returns the current state of a submitted task.
func (c *Client) GetTaskResult(ctx context.Context, taskID string) (*unicap.Result, error) {
	encodedID, err := c.EncodeTaskID(taskID)
	if err != nil {
		return nil, err
	}

	request := &getTaskResultRequest{
		TaskID: encodedID,
	}

	var response getTaskResultResponse
	if err := c.Call(ctx, "getTaskResult", request, &response); err != nil {
		return nil, err
	}

	status, err := c.taskStatus(response.Status)
	if err != nil {
		return nil, err
	}

	result := &unicap.Result{
		TaskID:    taskID,
		Status:    status,
		Cost:      float64(response.Cost),
		CreatedAt: response.CreateTime.Time(),
		SolvedAt:  response.EndTime.Time(),
	}

	if status == unicap.TaskStatusReady {
		result.Solution = decodeSolution(response.Solution)
	}

	return result, nil
}

// Balance returns the account balance in US dollars.
func (c *Client) Balance(ctx context.Context) (float64, error) {
	var response balanceResponse
	if err := c.Call(ctx, "getBalance", &Auth{}, &response); err != nil {
		return 0, err
	}

	return float64(response.Balance), nil
}

// Supports reports whether the provider has a mapper for taskType. Raw
// tasks are always accepted.
func (c *Client) Supports(taskType unicap.TaskType) bool {
	if taskType == unicap.TaskTypeRaw {
		return true
	}

	_, ok := c.config.Mappers[taskType]

	return ok
}

// RecommendedPolling returns the provider's recommended polling pace.
func (c *Client) RecommendedPolling() unicap.Polling {
	if c.config.Polling == (unicap.Polling{}) {
		return unicap.DefaultPolling()
	}

	return c.config.Polling
}

// TaskPayload maps task to the provider task object, as sent in createTask.
func (c *Client) TaskPayload(task unicap.Task) (any, error) {
	return c.taskPayload(task)
}

// EncodeTaskID converts a task ID to the JSON value the provider expects.
func (c *Client) EncodeTaskID(taskID string) (any, error) {
	encoded, err := c.config.TaskIDFormat.encode(taskID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", c.config.Name, err)
	}

	return encoded, nil
}

// taskPayload passes raw tasks through unchanged and maps every other task
// with the provider's registered mapper.
func (c *Client) taskPayload(task unicap.Task) (any, error) {
	if raw, ok := task.(*tasks.RawTask); ok {
		return raw.Payload(), nil
	}

	mapper, ok := c.config.Mappers[task.Type()]
	if !ok {
		return nil, fmt.Errorf("%s: %w", task.Type(), unicap.ErrUnsupportedTask)
	}

	return mapper(task)
}

// taskStatus converts a provider status string. A "failed" status without
// an error code becomes ErrTaskFailed.
func (c *Client) taskStatus(status string) (unicap.TaskStatus, error) {
	switch status {
	case statusReady:
		return unicap.TaskStatusReady, nil
	case statusProcessing:
		return unicap.TaskStatusProcessing, nil
	case statusFailed:
		return "", c.providerError("TASK_FAILED", "provider marked the task as failed", unicap.ErrTaskFailed)
	default:
		return unicap.TaskStatusPending, nil
	}
}

// providerError builds a provider error that unicap did not receive as an
// error code.
func (c *Client) providerError(code, message string, sentinel error) *unicap.Error {
	return &unicap.Error{
		Provider: c.config.Name,
		Code:     code,
		Message:  message,
		Err:      sentinel,
	}
}
