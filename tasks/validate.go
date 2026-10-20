package tasks

import (
	"errors"
	"fmt"

	"github.com/aarock1234/unicap/v3"
)

// validate joins every failed check into one error so callers see all
// missing inputs at once.
func validate(checks ...error) error {
	return errors.Join(checks...)
}

// required reports an error when a required text field is empty.
func required(field, value string) error {
	if value != "" {
		return nil
	}

	return fmt.Errorf("%s is required: %w", field, unicap.ErrInvalidTask)
}

// requiredItems reports an error when a required list is empty.
func requiredItems[T any](field string, values []T) error {
	if len(values) > 0 {
		return nil
	}

	return fmt.Errorf("%s needs at least one item: %w", field, unicap.ErrInvalidTask)
}

// exactlyOne reports an error unless exactly one of two alternative fields
// is set.
func exactlyOne(firstField, firstValue, secondField, secondValue string) error {
	switch {
	case firstValue == "" && secondValue == "":
		return fmt.Errorf("%s or %s is required: %w", firstField, secondField, unicap.ErrInvalidTask)
	case firstValue != "" && secondValue != "":
		return fmt.Errorf("%s and %s are mutually exclusive: %w", firstField, secondField, unicap.ErrInvalidTask)
	default:
		return nil
	}
}

// requiredInstruction reports an error when neither a text nor an image
// instruction is set.
func requiredInstruction(comment, imgInstructions string) error {
	if comment != "" || imgInstructions != "" {
		return nil
	}

	return fmt.Errorf("an instruction is required in Comment or ImgInstructions: %w", unicap.ErrInvalidTask)
}

// inRange reports an error when value falls outside [low, high].
func inRange(field string, value, low, high float64) error {
	if value >= low && value <= high {
		return nil
	}

	return fmt.Errorf("%s must be between %g and %g: %w", field, low, high, unicap.ErrInvalidTask)
}

// optionalProxy validates a proxy when one is set.
func optionalProxy(proxy *unicap.Proxy) error {
	if proxy == nil {
		return nil
	}

	if err := proxy.Validate(); err != nil {
		return fmt.Errorf("proxy: %w", errors.Join(unicap.ErrInvalidTask, err))
	}

	return nil
}

// requiredProxy validates a proxy and reports an error when none is set.
func requiredProxy(proxy *unicap.Proxy) error {
	if proxy == nil {
		return fmt.Errorf("proxy is required: %w", unicap.ErrInvalidTask)
	}

	return optionalProxy(proxy)
}
