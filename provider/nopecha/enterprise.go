package nopecha

import (
	"fmt"
	"maps"
	"slices"

	"github.com/aarock1234/unicap/v3/internal/solverapi"
)

// enterpriseValue returns the string an enterprise payload holds under
// supportedKey, the only key NopeCHA has a field for on the endpoint, or an
// empty string when the payload has none. Every other key, and a value that
// is not a string, is reported as unsupported, in key order.
func enterpriseValue(payload map[string]any, supportedKey string) (string, error) {
	var (
		supported   string
		unsupported []solverapi.Field
	)

	for _, key := range slices.Sorted(maps.Keys(payload)) {
		value, isString := payload[key].(string)
		if key == supportedKey && isString {
			supported = value

			continue
		}

		unsupported = append(unsupported, solverapi.Setting(fmt.Sprintf("EnterprisePayload[%s]", key), true))
	}

	if err := solverapi.Unsupported(unsupported...); err != nil {
		return "", err
	}

	return supported, nil
}
