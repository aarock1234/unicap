package nonecap

import (
	"net/url"
)

// solvesPath is the path of the solves collection in the native v1 API.
const solvesPath = "/v1/solves"

// Solve statuses that decide how a solve is read. Every other status, such
// as "solving", means the solve is in progress.
const (
	statusPending   = "pending"
	statusSolved    = "solved"
	statusFailed    = "failed"
	statusCancelled = "cancelled"
	statusExpired   = "expired"
)

// Feedback outcomes that tell NoneCap whether a token worked.
const (
	outcomeAccepted = "accepted"
	outcomeRejected = "rejected"
)

// solvePath returns the path of one solve. The ID is escaped so that it
// stays a single path segment.
func solvePath(solveID string) string {
	return solvesPath + "/" + url.PathEscape(solveID)
}

// feedbackPath returns the path that takes feedback on one solve.
func feedbackPath(solveID string) string {
	return solvePath(solveID) + "/feedback"
}

// solveObject is the solve NoneCap returns when one is created or fetched.
// It is flat: the token sits beside the status rather than in a solution
// object. Its timestamps are RFC 3339, and the fields that are null until
// the solve ends decode as empty.
type solveObject struct {
	ID         string     `json:"id"`
	Status     string     `json:"status"`
	Token      string     `json:"token"`
	UserAgent  string     `json:"user_agent"`
	Error      solveError `json:"error"`
	CreatedAt  string     `json:"created_at"`
	FinishedAt string     `json:"finished_at"`
}

// solveError explains why a solve ended without a token. It is empty unless
// the solve failed, was cancelled, or expired. Reason is an optional typed
// sub-reason of Code, and Retryable says whether the same request can
// succeed when sent again.
type solveError struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Reason    string `json:"reason"`
	Retryable bool   `json:"retryable"`
}

// errorEnvelope is the body of a request error.
type errorEnvelope struct {
	Error requestFailure `json:"error"`
}

// requestFailure explains why NoneCap rejected a request.
type requestFailure struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// feedbackRequest tells NoneCap how a solve's token fared.
type feedbackRequest struct {
	Outcome string `json:"outcome"`
}
