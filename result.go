package unicap

// TaskResult represents the result of a captcha solving task
type TaskResult struct {
	// Status indicates the current state of the task
	Status TaskStatus

	// Solution contains the captcha solution when ready
	Solution Solution

	// Error contains error details if the task failed
	Error *Error
}

// TaskStatus represents the state of a task
type TaskStatus string

const (
	// TaskStatusPending reports that a task has been accepted but not started.
	TaskStatusPending TaskStatus = "pending"
	// TaskStatusProcessing reports that a task is still being solved.
	TaskStatusProcessing TaskStatus = "processing"
	// TaskStatusReady reports that a solution is available.
	TaskStatusReady TaskStatus = "ready"
	// TaskStatusFailed reports that the provider could not solve the task.
	TaskStatusFailed TaskStatus = "failed"
)

// Solution contains the captcha solution data
type Solution struct {
	// Token is the primary solution (for token-based captchas)
	Token string

	// Text is the solution for image-based captchas
	Text string

	// Coordinates are for coordinate-based captchas
	Coordinates []Coordinate

	// Extra holds provider-specific additional data
	Extra map[string]any
}

// Coordinate represents a point selection
type Coordinate struct {
	X int
	Y int
}
