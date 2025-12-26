package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/aarock1234/unicap/unicap"
	"github.com/aarock1234/unicap/unicap/tasks"
)

// customProvider implements the unicap.Provider interface
type customProvider struct {
	apiKey string
	client *unicap.BaseHTTPClient
	errors *unicap.ErrorMapper
}

// NewCustomProvider creates a new custom provider
func NewCustomProvider(apiKey string) (unicap.Provider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("api key: %w", unicap.ErrInvalidAPIKey)
	}

	return &customProvider{
		apiKey: apiKey,
		client: &unicap.BaseHTTPClient{
			HTTPClient: &http.Client{Timeout: 30 * time.Second},
			Logger:     slog.New(slog.NewTextHandler(io.Discard, nil)),
			BaseURL:    "https://api.customservice.com",
		},
		errors: unicap.StandardErrorMapper(
			"customservice",
			[]string{"ERROR_INVALID_KEY"},
			[]string{"ERROR_NO_FUNDS"},
			[]string{"ERROR_TASK_NOT_FOUND"},
			[]string{"ERROR_BAD_REQUEST"},
		),
	}, nil
}

func (p *customProvider) CreateTask(ctx context.Context, task unicap.Task) (string, error) {
	customTask, err := mapTask(task)
	if err != nil {
		return "", fmt.Errorf("mapping task: %w", err)
	}

	req := createTaskRequest{
		APIKey: p.apiKey,
		Task:   customTask,
	}

	var resp createTaskResponse
	if err := p.client.DoJSON(ctx, "/create", req, &resp); err != nil {
		return "", err
	}

	if resp.Error != "" {
		return "", p.errors.MapError(resp.ErrorCode, resp.Error)
	}

	p.client.Logger.InfoContext(ctx, "task created",
		slog.String("task_id", resp.TaskID),
		slog.String("task_type", string(task.Type())),
	)

	return resp.TaskID, nil
}

func (p *customProvider) GetTaskResult(ctx context.Context, taskID string) (*unicap.TaskResult, error) {
	req := getResultRequest{
		APIKey: p.apiKey,
		TaskID: taskID,
	}

	var resp getResultResponse
	if err := p.client.DoJSON(ctx, "/result", req, &resp); err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return &unicap.TaskResult{
			Status: unicap.TaskStatusFailed,
			Error: &unicap.Error{
				Code:     resp.ErrorCode,
				Message:  resp.Error,
				Provider: "customservice",
			},
		}, nil
	}

	status := mapStatus(resp.Status)
	solution := unicap.Solution{}

	if status == unicap.TaskStatusReady {
		solution.Token = resp.Solution.Token
	}

	return &unicap.TaskResult{
		Status:   status,
		Solution: solution,
	}, nil
}

func (p *customProvider) Name() string {
	return "customservice"
}

// Request/Response types for the custom API
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

// Helper functions
func mapTask(task unicap.Task) (customTaskData, error) {
	switch t := task.(type) {
	case *tasks.ReCaptchaV2Task:
		return customTaskData{
			Type:       "recaptcha_v2",
			WebsiteURL: t.WebsiteURL,
			SiteKey:    t.WebsiteKey,
		}, nil
	default:
		return customTaskData{}, fmt.Errorf("unsupported task type: %s", task.Type())
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
	provider, err := NewCustomProvider("your-api-key")
	if err != nil {
		panic(err)
	}

	client, err := unicap.NewClient(provider)
	if err != nil {
		panic(err)
	}

	task := &tasks.ReCaptchaV2Task{
		WebsiteURL: "https://example.com",
		WebsiteKey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-",
	}

	ctx := context.Background()
	solution, err := client.Solve(ctx, task)
	if err != nil {
		panic(err)
	}

	fmt.Printf("solved: %s\n", solution.Token)
}
