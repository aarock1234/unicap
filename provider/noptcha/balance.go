package noptcha

import (
	"context"
	"net/http"
	"strconv"

	"github.com/aarock1234/unicap/v3/internal/httpapi"
)

// Balance returns the amount new tasks can spend, in US dollars. The API key
// needs the balance:read scope.
func (c *Client) Balance(ctx context.Context) (float64, error) {
	request := httpapi.Request{
		Method: http.MethodGet,
		Path:   balancePath,
	}

	response, err := c.api.Do(ctx, request)
	if err != nil {
		return 0, err
	}

	account, err := classifyReply[balance](c, request, response)
	if err != nil {
		return 0, err
	}

	available, err := strconv.ParseFloat(account.Available, 64)
	if err != nil {
		return 0, c.api.Error("INVALID_BALANCE", "the available balance is not a decimal amount", nil)
	}

	return available, nil
}
