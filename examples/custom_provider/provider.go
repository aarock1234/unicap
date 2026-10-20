package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/tasks"
)

const customProviderName = "customservice"

var (
	_ unicap.Provider      = (*customProvider)(nil)
	_ unicap.TaskSupporter = (*customProvider)(nil)
)

// customProvider implements unicap.Provider against a fictional captcha
// service, using only the standard library and the public unicap API.
type customProvider struct {
	apiKey  string
	baseURL string
	http    *http.Client
}

// newCustomProvider creates a custom provider.
func newCustomProvider(apiKey string) (*customProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("api key is empty: %w", unicap.ErrInvalidAPIKey)
	}

	return &customProvider{
		apiKey:  apiKey,
		baseURL: "https://api.customservice.example",
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// Name returns the provider identifier.
func (p *customProvider) Name() string {
	return customProviderName
}

// Supports reports whether the service solves tasks of the given type.
func (p *customProvider) Supports(taskType unicap.TaskType) bool {
	return taskType == unicap.TaskTypeReCaptchaV2
}

// CreateTask submits a task and returns its ID.
func (p *customProvider) CreateTask(ctx context.Context, task unicap.Task) (*unicap.Result, error) {
	payload, err := mapTask(task)
	if err != nil {
		return nil, err
	}

	request := createTaskRequest{
		APIKey: p.apiKey,
		Task:   payload,
	}

	var response createTaskResponse
	if err := p.post(ctx, "/create", request, &response); err != nil {
		return nil, err
	}

	if response.ErrorCode != "" {
		return nil, providerError(response.ErrorCode, response.Error)
	}

	return &unicap.Result{
		TaskID: response.TaskID,
		Status: unicap.TaskStatusPending,
	}, nil
}

// GetTaskResult returns the current state of a task.
func (p *customProvider) GetTaskResult(ctx context.Context, taskID string) (*unicap.Result, error) {
	request := getResultRequest{
		APIKey: p.apiKey,
		TaskID: taskID,
	}

	var response getResultResponse
	if err := p.post(ctx, "/result", request, &response); err != nil {
		return nil, err
	}

	if response.ErrorCode != "" {
		return nil, providerError(response.ErrorCode, response.Error)
	}

	if response.Status != "ready" {
		return &unicap.Result{
			TaskID: taskID,
			Status: unicap.TaskStatusProcessing,
		}, nil
	}

	return &unicap.Result{
		TaskID: taskID,
		Status: unicap.TaskStatusReady,
		Solution: unicap.Solution{
			Token: response.Solution.Token,
		},
	}, nil
}

// providerError classifies the service's error codes so callers can match
// them with errors.Is.
func providerError(code, message string) *unicap.Error {
	sentinels := map[string]error{
		"ERROR_INVALID_KEY": unicap.ErrInvalidAPIKey,
		"ERROR_NO_FUNDS":    unicap.ErrInsufficientFunds,
		"ERROR_UNSOLVABLE":  unicap.ErrUnsolvable,
	}

	return &unicap.Error{
		Provider: customProviderName,
		Code:     code,
		Message:  message,
		Err:      sentinels[code],
	}
}

// post sends a JSON request and decodes the JSON response.
func (p *customProvider) post(ctx context.Context, path string, request, response any) error {
	body, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("encoding request: %w", err)
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err := p.http.Do(httpRequest)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected http status %d", httpResponse.StatusCode)
	}

	if err := json.NewDecoder(httpResponse.Body).Decode(response); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}

	return nil
}

type createTaskRequest struct {
	APIKey string   `json:"apiKey"`
	Task   taskData `json:"task"`
}

type createTaskResponse struct {
	TaskID    string `json:"taskId"`
	Error     string `json:"error"`
	ErrorCode string `json:"errorCode"`
}

type getResultRequest struct {
	APIKey string `json:"apiKey"`
	TaskID string `json:"taskId"`
}

type getResultResponse struct {
	Status    string       `json:"status"`
	Solution  solutionData `json:"solution"`
	Error     string       `json:"error"`
	ErrorCode string       `json:"errorCode"`
}

type solutionData struct {
	Token string `json:"token"`
}

type taskData struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteUrl"`
	SiteKey    string `json:"siteKey"`
}

// mapTask converts the tasks this service supports.
func mapTask(task unicap.Task) (taskData, error) {
	reCaptcha, ok := task.(*tasks.ReCaptchaV2Task)
	if !ok {
		return taskData{}, fmt.Errorf("%s: %w", task.Type(), unicap.ErrUnsupportedTask)
	}

	if reCaptcha.Proxy != nil {
		return taskData{}, fmt.Errorf("proxy: %w", unicap.ErrUnsupportedOption)
	}

	return taskData{
		Type:       "recaptcha_v2",
		WebsiteURL: reCaptcha.WebsiteURL,
		SiteKey:    reCaptcha.WebsiteKey,
	}, nil
}
