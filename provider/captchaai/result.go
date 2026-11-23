package captchaai

import (
	"github.com/aarock1234/unicap/v3"
)

// newResult builds the unicap result of a task that is not solved. The cost
// and the timestamps stay zero because CaptchaAI sells threads and reports
// no times.
func newResult(ref taskRef, status unicap.TaskStatus) *unicap.Result {
	return &unicap.Result{
		TaskID: ref.String(),
		Status: status,
	}
}

// newReadyResult builds the unicap result of a solved task; answer is its
// reply and body is that reply as received.
func newReadyResult(ref taskRef, answer reply, body []byte) *unicap.Result {
	result := newResult(ref, unicap.TaskStatusReady)
	result.Solution = decodeSolution(ref.shape, answer, body)

	return result
}
