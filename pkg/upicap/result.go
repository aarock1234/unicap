package upicap

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
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusProcessing TaskStatus = "processing"
	TaskStatusReady      TaskStatus = "ready"
	TaskStatusFailed     TaskStatus = "failed"
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
