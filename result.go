package unicap

// TaskResult represents the result of a captcha solving task.
type TaskResult struct {
	// Status indicates the current state of the task.
	Status TaskStatus

	// Solution contains the captcha solution when ready.
	Solution Solution

	// Error contains error details if the task failed.
	Error *Error
}

// TaskStatus represents the state of a task.
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

// Solution contains the captcha solution data. Which field is populated depends
// on the captcha type: token-based captchas set Token, cookie-based challenges
// (DataDome, Cloudflare) set Cookie, image captchas set Text, and coordinate
// captchas set Coordinates.
type Solution struct {
	// Token is the primary solution for token-based captchas.
	Token string

	// Cookie is the solution for cookie-based challenges such as DataDome and
	// Cloudflare.
	Cookie string

	// Text is the solution for image-based captchas.
	Text string

	// Coordinates are the solution for coordinate-based captchas.
	Coordinates []Coordinate

	// Extra holds the raw, provider-specific solution payload.
	Extra map[string]any
}

// Coordinate represents a point selection within an image.
type Coordinate struct {
	X int
	Y int
}
