package solverapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aarock1234/unicap/v3/internal/httpapi"
)

// Call posts request to the named API method and decodes the reply into
// response. A non-zero errorId becomes a *unicap.Error classified by the
// provider's error rules, whatever the HTTP status. Server errors and rate
// limits without an error envelope are returned as plain transport errors so
// pollers treat them as transient, as are failures to send the request or
// read the reply. A request that cannot be encoded or built, and a response
// over the size limit, are final provider errors.
func (c *Client) Call(ctx context.Context, method string, request Request, response Response) error {
	request.authenticate(c.config.APIKey)

	reply, err := c.transport.Do(ctx, httpapi.Request{
		Method: http.MethodPost,
		Path:   "/" + method,
		Body:   request,
	})
	if err != nil {
		return err
	}

	var envelope Envelope
	if err := json.Unmarshal(reply.Body, &envelope); err == nil && envelope.ErrorID != 0 {
		return c.config.Errors.classify(c.config.Name, &envelope)
	}

	if reply.StatusCode != http.StatusOK {
		return c.transport.StatusError(method, reply.StatusCode)
	}

	if err := json.Unmarshal(reply.Body, response); err != nil {
		return fmt.Errorf("%s %s: decoding response: %w", c.config.Name, method, err)
	}

	return nil
}
