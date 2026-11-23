package easycaptcha

import (
	"context"
	"net/http"
	"strconv"

	"github.com/aarock1234/unicap/v3/internal/httpapi"
)

// balanceCurrency is the only currency whose balance unicap can report.
const balanceCurrency = "USD"

// microcentsPerDollar is how many of easycaptcha.to's microcents make one US
// dollar.
const microcentsPerDollar = 100_000_000

// Balance returns the account balance in US dollars.
func (c *Client) Balance(ctx context.Context) (float64, error) {
	request := httpapi.Request{
		Method: http.MethodGet,
		Path:   balancePath,
	}

	response, err := c.api.Do(ctx, request)
	if err != nil {
		return 0, err
	}

	account, err := classifyReply[balance](c, request, response, balanceStatuses)
	if err != nil {
		return 0, err
	}

	if account.Currency != balanceCurrency {
		return 0, c.api.Error("UNEXPECTED_CURRENCY", "the balance is not in us dollars", nil)
	}

	microcents, err := strconv.ParseInt(account.Microcents, 10, 64)
	if err != nil {
		return 0, c.api.Error("INVALID_BALANCE", "the balance is not a whole number of microcents", nil)
	}

	return float64(microcents) / microcentsPerDollar, nil
}
