package tasks

import (
	"fmt"
	"maps"

	"github.com/aarock1234/unicap"
)

// RawTask is a provider-specific passthrough task. It lets callers submit a
// captcha type that the SDK does not model directly by supplying the provider's
// task type string and raw parameters. The caller is responsible for matching
// the payload to the target provider; unicap sends Params as-is with the type
// field set to TaskType.
type RawTask struct {
	// TaskType is the provider-specific task type string, e.g. "AntiGateTask".
	TaskType string

	// Params holds the remaining provider-specific task fields.
	Params map[string]any
}

// Type returns the SDK task type identifier.
func (t *RawTask) Type() unicap.TaskType {
	return unicap.TaskTypeRaw
}

// Validate ensures the provider task type is set.
func (t *RawTask) Validate() error {
	if t.TaskType == "" {
		return fmt.Errorf("task_type: %w", unicap.ErrInvalidTask)
	}

	return nil
}

// Payload returns the provider task object: Params with the "type" field set to
// TaskType. TaskType takes precedence over any "type" key in Params.
func (t *RawTask) Payload() map[string]any {
	payload := make(map[string]any, len(t.Params)+1)
	maps.Copy(payload, t.Params)
	payload["type"] = t.TaskType

	return payload
}
