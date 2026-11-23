package noptcha

import (
	"strconv"
	"time"

	"github.com/aarock1234/unicap/v3"
)

// newResult builds the unicap result of a task that is processing or solved.
// The solution is decoded only once the task is solved.
func newResult(task taskObject) *unicap.Result {
	result := &unicap.Result{
		TaskID:    task.ID,
		Status:    taskStatus(task.Status),
		Cost:      dollars(task.Cost.Amount),
		CreatedAt: timestamp(task.CreatedAt),
		SolvedAt:  timestamp(task.CompletedAt),
	}

	if result.Ready() {
		result.Solution = decodeSolution(task.Solution)
	}

	return result
}

// dollars parses a decimal US-dollar amount. A missing or malformed amount
// is zero, which unicap reads as "not reported".
func dollars(amount string) float64 {
	value, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return 0
	}

	return value
}

// timestamp parses an RFC 3339 time. A missing or malformed time is the zero
// time, which unicap reads as "unknown".
func timestamp(text string) time.Time {
	parsed, err := time.Parse(time.RFC3339, text)
	if err != nil {
		return time.Time{}
	}

	return parsed
}
