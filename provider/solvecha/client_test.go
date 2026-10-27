package solvecha_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/provider/solvecha"
	"github.com/aarock1234/unicap/v3/tasks"
)

func newServer(t *testing.T, status int, reply string, got *map[string]any) *solvecha.Client {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got != nil {
			_ = json.NewDecoder(r.Body).Decode(got)
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(reply))
	}))
	t.Cleanup(server.Close)

	client, err := solvecha.New("test-key", solvecha.WithBaseURL(server.URL))
	if err != nil {
		t.Fatal(err)
	}

	return client
}

func TestImageToText(t *testing.T) {
	var got map[string]any
	client := newServer(t, 200, `{"errorId":0,"status":"ready","taskId":"1","solution":{"text":"abc12"}}`, &got)

	result, err := client.CreateTask(context.Background(), &tasks.ImageToTextTask{Body: "AAAA"})
	if err != nil {
		t.Fatal(err)
	}

	if result.Status != unicap.TaskStatusReady || result.Solution.Text != "abc12" {
		t.Fatalf("unexpected result: %+v", result)
	}

	task, _ := got["task"].(map[string]any)
	if got["clientKey"] != "test-key" || task["type"] != "ImageToTextTask" || task["body"] != "AAAA" {
		t.Fatalf("unexpected request: %v", got)
	}
}

func TestReCaptchaClassification(t *testing.T) {
	var got map[string]any
	client := newServer(t, 200, `{"errorId":0,"status":"ready","taskId":"2","solution":{"objects":[0,2,5]}}`, &got)

	result, err := client.CreateTask(context.Background(), &tasks.ReCaptchaClassificationTask{Image: "BBBB", Question: "/m/04_sv"})
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Solution.Indexes) != 3 || result.Solution.Indexes[2] != 5 {
		t.Fatalf("unexpected indexes: %v", result.Solution.Indexes)
	}

	task, _ := got["task"].(map[string]any)
	if task["type"] != "ReCaptchaV2Classification" || task["question"] != "/m/04_sv" {
		t.Fatalf("unexpected request: %v", got)
	}
}

func TestPendingTaskThenResult(t *testing.T) {
	client := newServer(t, 200, `{"errorId":0,"status":"processing","taskId":"3"}`, nil)

	result, err := client.CreateTask(context.Background(), &tasks.ImageToTextTask{Body: "AAAA"})
	if err != nil {
		t.Fatal(err)
	}

	if result.Status != unicap.TaskStatusPending || result.TaskID != "3" {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestErrorBodyOn4xx(t *testing.T) {
	client := newServer(t, 401, `{"errorId":1,"errorCode":"ERROR_KEY_DENIED_ACCESS","errorDescription":"Invalid key"}`, nil)

	_, err := client.CreateTask(context.Background(), &tasks.ImageToTextTask{Body: "AAAA"})
	if !errors.Is(err, unicap.ErrInvalidAPIKey) {
		t.Fatalf("want ErrInvalidAPIKey, got %v", err)
	}

	var apiErr *unicap.Error
	if !errors.As(err, &apiErr) || apiErr.Message != "Invalid key" {
		t.Fatalf("message not kept: %v", err)
	}
}

func TestBalanceZeroIsNotAnError(t *testing.T) {
	client := newServer(t, 200, `{"errorId":0,"balance":0,"packages":[]}`, nil)

	balance, err := client.Balance(context.Background())
	if err != nil || balance != 0 {
		t.Fatalf("balance = %v, err = %v", balance, err)
	}
}

func TestTokenTasksUnsupported(t *testing.T) {
	client := newServer(t, 200, `{}`, nil)

	for _, taskType := range []unicap.TaskType{
		unicap.TaskTypeReCaptchaV2,
		unicap.TaskTypeReCaptchaV3,
		unicap.TaskTypeTurnstile,
		unicap.TaskTypeGeeTest,
		unicap.TaskTypeAWSWAF,
	} {
		if client.Supports(taskType) {
			t.Errorf("%s should not be supported", taskType)
		}
	}

	if !client.Supports(unicap.TaskTypeImageToText) || !client.Supports(unicap.TaskTypeReCaptchaClassification) {
		t.Error("image tasks should be supported")
	}

	_, err := client.CreateTask(context.Background(), &tasks.TurnstileTask{WebsiteURL: "https://example.com", WebsiteKey: "k"})
	if !errors.Is(err, unicap.ErrUnsupportedTask) {
		t.Fatalf("want ErrUnsupportedTask, got %v", err)
	}
}
