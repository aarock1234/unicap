// Package solverapi implements the createTask / getTaskResult JSON protocol
// shared by the built-in providers. A provider supplies a Config with its
// endpoint, task mappers, error rules, and polling advice; this package owns
// the HTTP exchange, envelope and error decoding, solution normalisation,
// and the raw-task passthrough.
package solverapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"time"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/internal/httpapi"
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

	// Extras are the optional top-level createTask fields the provider
	// accepts.
	Extras []Extra

	// RefineSolution, when set, fills solution fields the shared decoder
	// does not know from the provider's raw solution payload.
	RefineSolution func(raw json.RawMessage, solution *unicap.Solution)
}

// Extra is an optional top-level createTask field.
type Extra string

const (
	// ExtraAppID is the developer appId.
	ExtraAppID Extra = "app id"

	// ExtraSoftID is the developer softId.
	ExtraSoftID Extra = "soft id"

	// ExtraCallbackURL is the URL the provider notifies when a task finishes.
	ExtraCallbackURL Extra = "callback url"

	// ExtraLanguagePool is the worker language pool.
	ExtraLanguagePool Extra = "language pool"
)

// Client is a provider client for the createTask / getTaskResult protocol.
// It is safe for concurrent use.
type Client struct {
	config    Config
	http      *http.Client
	logger    *slog.Logger
	transport *httpapi.Client
	extras    taskExtras
}

// taskExtras are optional top-level createTask fields.
type taskExtras struct {
	appID        string
	softID       int
	callbackURL  string
	languagePool string
}

// set returns the extras that carry a value, in the order their constants
// are declared. An empty string or a zero number means the extra is not set.
func (e taskExtras) set() []Extra {
	var extras []Extra

	if e.appID != "" {
		extras = append(extras, ExtraAppID)
	}

	if e.softID != 0 {
		extras = append(extras, ExtraSoftID)
	}

	if e.callbackURL != "" {
		extras = append(extras, ExtraCallbackURL)
	}

	if e.languagePool != "" {
		extras = append(extras, ExtraLanguagePool)
	}

	return extras
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

	for _, extra := range c.extras.set() {
		if !slices.Contains(config.Extras, extra) {
			return nil, fmt.Errorf("%s does not accept a %s: %w", config.Name, extra, unicap.ErrUnsupportedOption)
		}
	}

	c.transport = httpapi.New(
		httpapi.Config{
			Name:    c.config.Name,
			BaseURL: c.config.BaseURL,
		},
		httpapi.WithHTTPClient(c.http),
		httpapi.WithLogger(c.logger),
	)

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
			Solution: c.solution(response.Solution),
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
		result.Solution = c.solution(response.Solution)
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

// EncodeTaskID converts a task ID to the JSON value the provider expects. A
// malformed ID is a final provider error wrapping unicap.ErrTaskNotFound.
func (c *Client) EncodeTaskID(taskID string) (any, error) {
	encoded, err := c.config.TaskIDFormat.encode(taskID)
	if err != nil {
		return nil, c.providerError("INVALID_TASK_ID", err.Error(), unicap.ErrTaskNotFound)
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

// solution decodes a provider solution payload and applies the provider's
// refinement, if any.
func (c *Client) solution(raw json.RawMessage) unicap.Solution {
	solution := decodeSolution(raw)

	if c.config.RefineSolution != nil {
		c.config.RefineSolution(raw, &solution)
	}

	return solution
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
