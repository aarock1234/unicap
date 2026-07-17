package tasks

import (
	"fmt"

	"github.com/aarock1234/unicap"
)

// TextCaptchaTask represents a text captcha solving task, where a worker
// answers a natural-language question.
type TextCaptchaTask struct {
	Question string
}

// Type returns the SDK task type identifier.
func (t *TextCaptchaTask) Type() unicap.TaskType {
	return unicap.TaskTypeText
}

// Validate ensures required fields are present.
func (t *TextCaptchaTask) Validate() error {
	if t.Question == "" {
		return fmt.Errorf("question: %w", unicap.ErrInvalidTask)
	}

	return nil
}
