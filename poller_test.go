package unicap

import (
	"context"
	"errors"
	"testing"
	"time"
)

// fakeProvider returns a scripted sequence of results, repeating the last entry
// once exhausted.
type fakeProvider struct {
	steps []step
	calls int
}

type step struct {
	result *TaskResult
	err    error
}

func (f *fakeProvider) CreateTask(context.Context, Task) (string, error) {
	return "task-1", nil
}

func (f *fakeProvider) GetTaskResult(context.Context, string) (*TaskResult, error) {
	i := min(f.calls, len(f.steps)-1)
	f.calls++

	return f.steps[i].result, f.steps[i].err
}

func (f *fakeProvider) Name() string {
	return "fake"
}

func testConfig() PollerConfig {
	return PollerConfig{
		InitialInterval: time.Millisecond,
		MaxInterval:     2 * time.Millisecond,
		Timeout:         500 * time.Millisecond,
		Multiplier:      1,
	}
}

func ready() *TaskResult {
	return &TaskResult{Status: TaskStatusReady, Solution: Solution{Token: "solved"}}
}

func processing() *TaskResult {
	return &TaskResult{Status: TaskStatusProcessing}
}

func TestPollerPoll(t *testing.T) {
	transient := errors.New("transient")

	tests := []struct {
		name      string
		steps     []step
		wantToken string
		wantErr   bool
		wantErrIs error
	}{
		{
			name:      "ready immediately",
			steps:     []step{{result: ready()}},
			wantToken: "solved",
		},
		{
			name:      "processing then ready",
			steps:     []step{{result: processing()}, {result: ready()}},
			wantToken: "solved",
		},
		{
			name: "failed returns provider error",
			steps: []step{{result: &TaskResult{
				Status: TaskStatusFailed,
				Error:  NewError("ERROR_X", "boom", "fake", false, ErrInvalidTask),
			}}},
			wantErr:   true,
			wantErrIs: ErrInvalidTask,
		},
		{
			name: "tolerates transient errors then succeeds",
			steps: []step{
				{err: transient},
				{err: transient},
				{err: transient},
				{result: ready()},
			},
			wantToken: "solved",
		},
		{
			name: "gives up after too many transient errors",
			steps: []step{
				{err: transient},
				{err: transient},
				{err: transient},
				{err: transient},
			},
			wantErr:   true,
			wantErrIs: transient,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &fakeProvider{steps: tt.steps}
			poller := NewPoller(provider, testConfig())

			result, err := poller.Poll(context.Background(), "task-1")

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				if tt.wantErrIs != nil && !errors.Is(err, tt.wantErrIs) {
					t.Errorf("errors.Is(%v, %v) = false, want true", err, tt.wantErrIs)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.Solution.Token != tt.wantToken {
				t.Errorf("token = %q, want %q", result.Solution.Token, tt.wantToken)
			}
		})
	}
}

func TestPollerPollTimeout(t *testing.T) {
	provider := &fakeProvider{steps: []step{{result: processing()}}}

	config := testConfig()
	config.Timeout = 15 * time.Millisecond

	poller := NewPoller(provider, config)

	_, err := poller.Poll(context.Background(), "task-1")
	if !errors.Is(err, ErrTimeout) {
		t.Fatalf("errors.Is(%v, ErrTimeout) = false, want true", err)
	}
}

func TestPollerPollContextCancelled(t *testing.T) {
	provider := &fakeProvider{steps: []step{{result: processing()}}}
	poller := NewPoller(provider, testConfig())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := poller.Poll(ctx, "task-1")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("errors.Is(%v, context.Canceled) = false, want true", err)
	}
}
