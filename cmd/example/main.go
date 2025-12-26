package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aarock1234/unicap/providers/capsolver"
	"github.com/aarock1234/unicap/unicap"
	"github.com/aarock1234/unicap/unicap/tasks"

	_ "github.com/aarock1234/unicap/internal/logger"
)

func main() {
	apiKey := os.Getenv("CAPSOLVER_API_KEY")
	if apiKey == "" {
		slog.Error("CAPSOLVER_API_KEY environment variable is required")
		os.Exit(1)
	}

	provider, err := capsolver.NewCapSolverProvider(apiKey)
	if err != nil {
		slog.Error("failed to create provider", slog.Any("error", err))
		os.Exit(1)
	}

	client, err := unicap.NewClient(
		provider,
		unicap.WithLogger(slog.Default()),
	)
	if err != nil {
		slog.Error("failed to create client", slog.Any("error", err))
		os.Exit(1)
	}

	exampleSynchronous(client)
	exampleAsynchronous(client)
}

func exampleSynchronous(client *unicap.Client) {
	slog.Info("synchronous api example")

	task := &tasks.ReCaptchaV2Task{
		WebsiteURL: "https://www.google.com/recaptcha/api2/demo",
		WebsiteKey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	solution, err := client.Solve(ctx, task)
	if err != nil {
		slog.ErrorContext(ctx, "failed to solve captcha", slog.Any("error", err))
		return
	}

	slog.InfoContext(ctx, "captcha solved", slog.String("token", solution.Token))

	if err := verifyReCaptchaV2(ctx, solution.Token); err != nil {
		slog.ErrorContext(ctx, "verification failed", slog.Any("error", err))
		return
	}
}

func exampleAsynchronous(client *unicap.Client) {
	slog.Info("asynchronous api example")

	task := &tasks.ReCaptchaV3Task{
		WebsiteURL: "https://recaptcha-demo.appspot.com/recaptcha-v3-request-scores.php",
		WebsiteKey: "6LdKlZEpAAAAAAOQjzC2v_d36tWxCl6dWsozdSy9",
		PageAction: "examples/v3scores",
		// MinScore:   0.7, // CapSolver does not support MinScore
	}

	ctx := context.Background()

	taskID, err := client.CreateTask(ctx, task)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create task", slog.Any("error", err))
		return
	}

	slog.InfoContext(ctx, "task created", slog.String("task_id", taskID))

	time.Sleep(10 * time.Second)

	result, err := client.GetTaskResult(ctx, taskID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get task result", slog.Any("error", err))
		return
	}

	switch result.Status {
	case unicap.TaskStatusReady:
		slog.InfoContext(ctx, "captcha solved", slog.String("token", result.Solution.Token))
		if err := verifyReCaptchaV3(ctx, result.Solution.Token); err != nil {
			slog.ErrorContext(ctx, "verification failed", slog.Any("error", err))
		}
	case unicap.TaskStatusProcessing:
		slog.InfoContext(ctx, "still processing, check again later")
	case unicap.TaskStatusFailed:
		slog.ErrorContext(ctx, "task failed", slog.Any("error", result.Error))
	}
}

func verifyReCaptchaV2(ctx context.Context, token string) error {
	slog.InfoContext(ctx, "verifying token with recaptcha v2 demo")

	data := url.Values{
		"g-recaptcha-response": {token},
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://www.google.com/recaptcha/api2/demo", strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}

	slog.InfoContext(ctx, "verification response", slog.String("body", string(body)))
	if strings.Contains(string(body), "Verification Success") {
		slog.InfoContext(ctx, "verification successful")
	} else {
		slog.WarnContext(ctx, "verification may have failed - check response above")
	}

	return nil
}

type verifyReCaptchaV3Response struct {
	Success    bool     `json:"success"`
	Score      float64  `json:"score"`
	ErrorCodes []string `json:"error-codes"`
}

func verifyReCaptchaV3(ctx context.Context, token string) error {
	slog.InfoContext(ctx, "verifying token with recaptcha v3 demo")

	data := url.Values{
		"token": {token},
	}

	verifyURL, err := url.Parse("https://recaptcha-demo.appspot.com/recaptcha-v3-verify.php")
	if err != nil {
		return fmt.Errorf("parsing URL: %w", err)
	}

	verifyURL.RawQuery = data.Encode()

	req, err := http.NewRequestWithContext(ctx, "POST", verifyURL.String(), nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	var verifyResponse verifyReCaptchaV3Response
	if err := json.NewDecoder(resp.Body).Decode(&verifyResponse); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}

	if !verifyResponse.Success {
		return fmt.Errorf("verification failed: %s", verifyResponse.ErrorCodes)
	}

	slog.InfoContext(ctx, "verification successful", slog.Float64("score", verifyResponse.Score))

	return nil
}
