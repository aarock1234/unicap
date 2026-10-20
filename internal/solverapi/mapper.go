package solverapi

import (
	"errors"
	"fmt"

	"github.com/aarock1234/unicap/v3"
)

// Mapper converts a unicap task into a provider task object.
type Mapper func(task unicap.Task) (any, error)

// Mappers registers a provider's supported task types. Its keys are the
// single source of truth for Client.Supports.
type Mappers map[unicap.TaskType]Mapper

// Map adapts a typed mapping function to a Mapper. The adapter asserts the
// concrete task type, so each provider mapper can accept and return precise
// types.
func Map[T unicap.Task, W any](mapTask func(T) (W, error)) Mapper {
	return func(task unicap.Task) (any, error) {
		typed, ok := task.(T)
		if !ok {
			return nil, fmt.Errorf("%s task has unexpected type %T: %w", task.Type(), task, unicap.ErrUnsupportedTask)
		}

		return mapTask(typed)
	}
}

// Variant names a task type's proxyless and proxy wire types. Leave Proxy
// empty when the provider has no proxy variant, or Proxyless empty when it
// requires a proxy.
type Variant struct {
	Proxyless string
	Proxy     string
}

// For returns the wire type for the task's proxy setting. It reports
// ErrUnsupportedOption when a proxy is set but the provider has no proxy
// variant, and ErrInvalidTask when the provider requires a proxy that is not
// set.
func (v Variant) For(proxy *unicap.Proxy) (string, error) {
	switch {
	case proxy != nil && v.Proxy == "":
		return "", fmt.Errorf("proxy: %w", unicap.ErrUnsupportedOption)
	case proxy == nil && v.Proxyless == "":
		return "", fmt.Errorf("proxy is required: %w", unicap.ErrInvalidTask)
	case proxy != nil:
		return v.Proxy, nil
	default:
		return v.Proxyless, nil
	}
}

// Field describes whether a task field carries a value.
type Field struct {
	name string
	set  bool
}

// Setting describes the task field name and whether it is set.
func Setting(name string, set bool) Field {
	return Field{
		name: name,
		set:  set,
	}
}

// Unsupported reports every set field the provider cannot honour, wrapped in
// unicap.ErrUnsupportedOption.
func Unsupported(fields ...Field) error {
	var errs []error
	for _, field := range fields {
		if field.set {
			errs = append(errs, fmt.Errorf("%s: %w", field.name, unicap.ErrUnsupportedOption))
		}
	}

	return errors.Join(errs...)
}

// Required reports every unset field the provider needs, wrapped in
// unicap.ErrInvalidTask.
func Required(fields ...Field) error {
	var errs []error
	for _, field := range fields {
		if !field.set {
			errs = append(errs, fmt.Errorf("%s is required by this provider: %w", field.name, unicap.ErrInvalidTask))
		}
	}

	return errors.Join(errs...)
}

// Check joins the results of Unsupported, Required, and similar checks.
func Check(checks ...error) error {
	return errors.Join(checks...)
}
