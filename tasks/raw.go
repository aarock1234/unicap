package tasks

import (
	"maps"

	"github.com/aarock1234/unicap/v3"
)

// RawTask is a provider-specific passthrough for captcha types or options
// unicap does not model. The built-in providers send Params as the task
// object with its "type" field set to TaskType, so the caller is responsible
// for matching the payload to the selected provider.
type RawTask struct {
	// TaskType is the provider's task type string, such as "AntiGateTask".
	TaskType string

	// Params holds the remaining provider task fields. Their names and types
	// are defined by the provider, so the map is intentionally dynamic.
	Params map[string]any
}

// Type returns unicap.TaskTypeRaw.
func (t *RawTask) Type() unicap.TaskType {
	return unicap.TaskTypeRaw
}

// Validate reports a missing provider task type.
func (t *RawTask) Validate() error {
	return validate(
		required("TaskType", t.TaskType),
	)
}

// Payload returns the provider task object: a copy of Params with the
// "type" field set to TaskType, which takes precedence over any "type" key
// in Params.
func (t *RawTask) Payload() map[string]any {
	payload := make(map[string]any, len(t.Params)+1)
	maps.Copy(payload, t.Params)
	payload["type"] = t.TaskType

	return payload
}
