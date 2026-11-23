package captchaai

import (
	"fmt"
	"maps"
	"slices"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/tasks"
)

// providerFields are the form fields the provider sets itself, which a raw
// task therefore cannot carry.
var providerFields = []string{
	fieldKey,
	fieldMethod,
	fieldJSON,
}

// submission is a task mapped for in.php: its form fields and the shape its
// solution will arrive in.
type submission struct {
	shape  solutionShape
	fields form
}

// ref returns the reference of the task CaptchaAI created as vendorTaskID.
func (s submission) ref(vendorTaskID string) taskRef {
	return taskRef{
		shape: s.shape,
		id:    vendorTaskID,
	}
}

// newSubmission maps task to its submission. A raw task names its own
// method; every other task is mapped by its registered mapper.
func newSubmission(task unicap.Task) (submission, error) {
	if raw, ok := task.(*tasks.RawTask); ok {
		return newRawSubmission(raw)
	}

	mapper, ok := mappers[task.Type()]
	if !ok {
		return submission{}, fmt.Errorf("%s: %w", task.Type(), unicap.ErrUnsupportedTask)
	}

	mapped, err := mapper(task)
	if err != nil {
		return submission{}, err
	}

	mappedSubmission, ok := mapped.(submission)
	if !ok {
		return submission{}, fmt.Errorf("%s mapper returned unexpected type %T: %w", task.Type(), mapped, unicap.ErrUnsupportedTask)
	}

	return mappedSubmission, nil
}

// newRawSubmission passes a raw task through: its TaskType is the method,
// such as "bls", and each of its Params is a form field, in name order. A
// param named like a field the provider sets itself is rejected, because the
// form would carry that field twice.
func newRawSubmission(raw *tasks.RawTask) (submission, error) {
	fields := newForm(raw.TaskType)

	for _, name := range slices.Sorted(maps.Keys(raw.Params)) {
		if slices.Contains(providerFields, name) {
			return submission{}, fmt.Errorf("raw task param %q is set by the provider: %w", name, unicap.ErrInvalidTask)
		}

		fields = fields.with(name, formValue(raw.Params[name]))
	}

	return submission{
		shape:  shapeRaw,
		fields: fields,
	}, nil
}

// formValue returns the form value of a raw task's param: a string as it is,
// and any other value in its default format.
func formValue(param any) string {
	if text, ok := param.(string); ok {
		return text
	}

	return fmt.Sprint(param)
}
