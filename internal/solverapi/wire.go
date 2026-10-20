package solverapi

import (
	"encoding/json"
)

// Provider task status strings.
const (
	statusProcessing = "processing"
	statusReady      = "ready"
	statusFailed     = "failed"
)

// Auth carries the client key. Embed it in every request type; Call fills it
// in before sending.
type Auth struct {
	ClientKey string `json:"clientKey"`
}

// authenticate sets the client key.
func (a *Auth) authenticate(apiKey string) {
	a.ClientKey = apiKey
}

// Request is a provider request body. Types satisfy it by embedding Auth.
type Request interface {
	authenticate(apiKey string)
}

// Envelope is the error header every provider response carries. Embed it in
// every response type.
type Envelope struct {
	ErrorID          int    `json:"errorId"`
	ErrorCode        string `json:"errorCode,omitempty"`
	ErrorDescription string `json:"errorDescription,omitempty"`
}

// envelope returns the embedded envelope.
func (e *Envelope) envelope() *Envelope {
	return e
}

// Response is a provider response body. Types satisfy it by embedding
// Envelope.
type Response interface {
	envelope() *Envelope
}

type createTaskRequest struct {
	Auth

	Task         any    `json:"task"`
	AppID        string `json:"appId,omitempty"`
	SoftID       int    `json:"softId,omitempty"`
	CallbackURL  string `json:"callbackUrl,omitempty"`
	LanguagePool string `json:"languagePool,omitempty"`
}

type createTaskResponse struct {
	Envelope

	TaskID   TaskID          `json:"taskId"`
	Status   string          `json:"status"`
	Solution json.RawMessage `json:"solution"`
}

type getTaskResultRequest struct {
	Auth

	TaskID any `json:"taskId"`
}

type getTaskResultResponse struct {
	Envelope

	Status     string          `json:"status"`
	Solution   json.RawMessage `json:"solution"`
	Cost       Number          `json:"cost"`
	CreateTime Number          `json:"createTime"`
	EndTime    Number          `json:"endTime"`
}

type balanceResponse struct {
	Envelope

	Balance Number `json:"balance"`
}
