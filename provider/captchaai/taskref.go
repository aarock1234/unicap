package captchaai

import (
	"strings"

	"github.com/aarock1234/unicap/v3"
)

// taskRefSeparator parts the solution shape from CaptchaAI's task ID in a
// unicap task ID.
const taskRefSeparator = ":"

// taskRef names a CaptchaAI task. A reply to res.php does not say which kind
// of answer it holds, so CaptchaAI's task ID alone does not tell how to
// decode it. As a string, which is the unicap task ID, it reads
// "<shape>:<task id>".
type taskRef struct {
	shape solutionShape
	id    string
}

// String returns the unicap task ID of the task.
func (t taskRef) String() string {
	return string(t.shape) + taskRefSeparator + t.id
}

// fields returns the form fields that fetch the task's result. Only
// CaptchaAI's own task ID is sent.
func (t taskRef) fields() form {
	return form{}.
		with(fieldAction, actionGet).
		with(fieldID, t.id)
}

// parseTaskRef reads a unicap task ID. The ID is cut at the first separator;
// the shape must be one a decoder is registered for, and the rest must have
// the form of a task ID CaptchaAI issues. A malformed ID is a final provider
// error wrapping unicap.ErrTaskNotFound.
func parseTaskRef(taskID string) (taskRef, error) {
	shapeName, id, found := strings.Cut(taskID, taskRefSeparator)
	shape := solutionShape(shapeName)

	_, known := decoders[shape]
	if !found || !known || !vendorTaskIDPattern.MatchString(id) {
		return taskRef{}, &unicap.Error{
			Provider: Name,
			Code:     "INVALID_TASK_ID",
			Message:  `task id must read "<shape>:<task id>" with a known shape and a task id of decimal digits`,
			Err:      unicap.ErrTaskNotFound,
		}
	}

	return taskRef{
		shape: shape,
		id:    id,
	}, nil
}
