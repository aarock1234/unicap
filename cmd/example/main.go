package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/aarock1234/unicap/internal/providers/capsolver"
	"github.com/aarock1234/unicap/pkg/upicap"
	"github.com/aarock1234/unicap/pkg/upicap/tasks"

	_ "github.com/aarock1234/unicap/internal/logger"
)

func main() {
	apiKey := os.Getenv("CAPSOLVER_API_KEY")
	if apiKey == "" {
		log.Fatal("CAPSOLVER_API_KEY environment variable is required")
	}

	provider, err := capsolver.NewCapSolverProvider(apiKey)
	if err != nil {
		log.Fatal(err)
	}

	client, err := upicap.NewClient(
		provider,
		upicap.WithLogger(slog.Default()),
	)
	if err != nil {
		log.Fatal(err)
	}

	exampleSynchronous(client)
	exampleAsynchronous(client)

	fmt.Println("\nVerify the token with:")
	fmt.Println("curl 'https://www.google.com/recaptcha/api2/demo' \\")
	fmt.Println("  -H 'content-type: application/x-www-form-urlencoded' \\")
	fmt.Println("  --data-raw \"g-recaptcha-response=${TOKEN}\"")
	fmt.Println("\nExpected: Verification Success... Hooray!")
}

func exampleSynchronous(client *upicap.Client) {
	fmt.Println("=== Synchronous API Example ===")

	task := &tasks.ReCaptchaV2Task{
		WebsiteURL: "https://example.com",
		WebsiteKey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	solution, err := client.Solve(ctx, task)
	if err != nil {
		log.Printf("failed to solve captcha: %v", err)
		return
	}

	fmt.Printf("captcha solved: %s\n", solution.Token)
}

func exampleAsynchronous(client *upicap.Client) {
	fmt.Println("\n=== Asynchronous API Example ===")

	task := &tasks.ReCaptchaV3Task{
		WebsiteURL: "https://example.com",
		WebsiteKey: "6LdyC2cUAAAAACGuDKpXeDorzUDWXmdqeg-xy696",
		PageAction: "verify",
		MinScore:   0.7,
	}

	ctx := context.Background()

	taskID, err := client.CreateTask(ctx, task)
	if err != nil {
		log.Printf("failed to create task: %v", err)
		return
	}

	fmt.Printf("task created: %s\n", taskID)

	time.Sleep(10 * time.Second)

	result, err := client.GetTaskResult(ctx, taskID)
	if err != nil {
		log.Printf("failed to get task result: %v", err)
		return
	}

	switch result.Status {
	case upicap.TaskStatusReady:
		fmt.Printf("captcha solved: %s\n", result.Solution.Token)
	case upicap.TaskStatusProcessing:
		fmt.Println("still processing, check again later")
	case upicap.TaskStatusFailed:
		fmt.Printf("task failed: %v\n", result.Error)
	}
}
