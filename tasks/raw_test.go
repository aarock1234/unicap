package tasks

import (
	"errors"
	"testing"

	"github.com/aarock1234/unicap"
)

func TestRawTaskValidate(t *testing.T) {
	tests := []struct {
		name    string
		task    RawTask
		wantErr bool
	}{
		{
			name:    "valid",
			task:    RawTask{TaskType: "AntiGateTask"},
			wantErr: false,
		},
		{
			name:    "missing type",
			task:    RawTask{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.task.Validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && !errors.Is(err, unicap.ErrInvalidTask) {
				t.Errorf("error = %v, want wrapped ErrInvalidTask", err)
			}
		})
	}
}

func TestRawTaskPayload(t *testing.T) {
	task := RawTask{
		TaskType: "AntiGateTask",
		Params: map[string]any{
			"websiteURL": "https://example.com",
			"type":       "should-be-overridden",
		},
	}

	payload := task.Payload()

	if got := payload["type"]; got != "AntiGateTask" {
		t.Errorf("type = %v, want AntiGateTask", got)
	}

	if got := payload["websiteURL"]; got != "https://example.com" {
		t.Errorf("websiteURL = %v, want https://example.com", got)
	}
}
