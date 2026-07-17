// Package main demonstrates implementing a custom provider for unicap using
// only the standard library and the public unicap API.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/aarock1234/unicap"
	"github.com/aarock1234/unicap/tasks"
)

// customProvider implements the unicap.Provider interface against a fictional
// captcha service.
type customProvider struct {
	apiKey  string
	http    *http.Client
	baseURL string
}

var _ unicap.Provider = (*customProvider)(nil)

// newCustomProvider creates a custom provider.
func newCustomProvider(apiKey string) (*customProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("api key: %w", unicap.ErrInvalidAPIKey)
	}

	return &customProvider{
		apiKey:  apiKey,
		http:    &http.Client{Timeout: 30 * time.Second},
		baseURL: "https://api.customservice.com",
	}, nil
}

// CreateTask submits a captcha task and returns its ID.
func (p *customProvider) CreateTask(ctx context.Context, task unicap.Task) (string, error) {
	body, err := mapTask(task)
	if err != nil {
		return "", fmt.Errorf("mapping task: %w", err)
	}

	req := createTaskRequest{
		APIKey: p.apiKey,
		Task:   body,
	}

	var resp createTaskResponse
	if err := p.doJSON(ctx, "/create", req, &resp); err != nil {
		return "", err
	}

	if resp.Error != "" {
		return "", p.mapError(resp.ErrorCode, resp.Error)
	}

	return resp.TaskID, nil
}

// GetTaskResult retrieves the result for a task ID.
func (p *customProvider) GetTaskResult(ctx context.Context, taskID string) (*unicap.TaskResult, error) {
	req := getResultRequest{
		APIKey: p.apiKey,
		TaskID: taskID,
	}

	var resp getResultResponse
	if err := p.doJSON(ctx, "/result", req, &resp); err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return &unicap.TaskResult{
			Status: unicap.TaskStatusFailed,
			Error:  p.mapError(resp.ErrorCode, resp.Error),
		}, nil
	}

	status := mapStatus(resp.Status)

	var solution unicap.Solution
	if status == unicap.TaskStatusReady {
		solution.Token = resp.Solution.Token
	}

	return &unicap.TaskResult{
		Status:   status,
		Solution: solution,
	}, nil
}

// Name returns the provider identifier.
func (p *customProvider) Name() string {
	return "customservice"
}

func (p *customProvider) mapError(code, message string) *unicap.Error {
	switch code {
	case "ERROR_INVALID_KEY":
		return unicap.NewError(code, message, p.Name(), false, unicap.ErrInvalidAPIKey)
	case "ERROR_NO_FUNDS":
		return unicap.NewError(code, message, p.Name(), false, unicap.ErrInsufficientFunds)
	default:
		return unicap.NewError(code, message, p.Name(), true, nil)
	}
}

func (p *customProvider) doJSON(ctx context.Context, endpoint string, reqBody, respBody any) error {
	data, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+endpoint, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.http.Do(req)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(respBody); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}

	return nil
}

type createTaskRequest struct {
	APIKey string         `json:"apiKey"`
	Task   customTaskData `json:"task"`
}

type createTaskResponse struct {
	TaskID    string `json:"taskId,omitempty"`
	Error     string `json:"error,omitempty"`
	ErrorCode string `json:"errorCode,omitempty"`
}

type getResultRequest struct {
	APIKey string `json:"apiKey"`
	TaskID string `json:"taskId"`
}

type getResultResponse struct {
	Status    string       `json:"status,omitempty"`
	Solution  solutionData `json:"solution,omitempty"`
	Error     string       `json:"error,omitempty"`
	ErrorCode string       `json:"errorCode,omitempty"`
}

type solutionData struct {
	Token string `json:"token"`
}

type customTaskData struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteUrl"`
	SiteKey    string `json:"siteKey"`
}

func mapTask(task unicap.Task) (customTaskData, error) {
	switch t := task.(type) {
	case *tasks.ReCaptchaV2Task:
		return customTaskData{
			Type:       "recaptcha_v2",
			WebsiteURL: t.WebsiteURL,
			SiteKey:    t.WebsiteKey,
		}, nil
	default:
		return customTaskData{}, fmt.Errorf("%s: %w", task.Type(), unicap.ErrUnsupportedTask)
	}
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

func main() {
	if err := run(); err != nil {
		slog.Error("custom provider example failed", slog.Any("error", err))
		os.Exit(1)
	}
}

func run() error {
	provider, err := newCustomProvider("your-api-key")
	if err != nil {
		return fmt.Errorf("creating provider: %w", err)
	}

	client, err := unicap.New(provider)
	if err != nil {
		return fmt.Errorf("creating client: %w", err)
	}

	task := &tasks.ReCaptchaV2Task{
		WebsiteURL: "https://example.com",
		WebsiteKey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-",
	}

	ctx := context.Background()
	solution, err := client.Solve(ctx, task)
	if err != nil {
		return err
	}

	fmt.Printf("solved: %s\n", solution.Token)

	return nil
}
