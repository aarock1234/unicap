package unicap

import (
	"errors"
	"math"
	"time"
)

// Polling controls how Client waits for a task to finish.
type Polling struct {
	// InitialDelay is how long to wait after creating a task before the
	// first result check. Most providers need several seconds for token
	// captchas.
	InitialDelay time.Duration

	// Interval is the wait between the first checks.
	Interval time.Duration

	// MaxInterval caps the wait between checks as it grows by Multiplier.
	MaxInterval time.Duration

	// Multiplier grows the interval after every unfinished check. A value of
	// 1 keeps a constant interval.
	Multiplier float64

	// Timeout bounds the total wait for one task. Zero relies on the caller's
	// context alone.
	Timeout time.Duration

	// MaxTransientErrors is the number of consecutive transport or temporary
	// provider failures tolerated before waiting gives up.
	MaxTransientErrors int
}

// DefaultPolling returns polling settings that suit most providers.
func DefaultPolling() Polling {
	return Polling{
		InitialDelay:       3 * time.Second,
		Interval:           3 * time.Second,
		MaxInterval:        10 * time.Second,
		Multiplier:         1.5,
		Timeout:            5 * time.Minute,
		MaxTransientErrors: 3,
	}
}

// ErrInvalidPolling reports that polling settings are unusable.
var ErrInvalidPolling = errors.New("invalid polling settings")

// Validate reports whether the settings can drive a polling loop without
// spinning or stalling.
func (p Polling) Validate() error {
	switch {
	case p.InitialDelay < 0:
		return errors.Join(ErrInvalidPolling, errors.New("initial delay must not be negative"))
	case p.Interval <= 0:
		return errors.Join(ErrInvalidPolling, errors.New("interval must be positive"))
	case p.MaxInterval < p.Interval:
		return errors.Join(ErrInvalidPolling, errors.New("max interval must not be less than interval"))
	case math.IsNaN(p.Multiplier) || math.IsInf(p.Multiplier, 0) || p.Multiplier < 1:
		return errors.Join(ErrInvalidPolling, errors.New("multiplier must be finite and at least 1"))
	case p.Timeout < 0:
		return errors.Join(ErrInvalidPolling, errors.New("timeout must not be negative"))
	case p.MaxTransientErrors < 0:
		return errors.Join(ErrInvalidPolling, errors.New("max transient errors must not be negative"))
	default:
		return nil
	}
}

// next returns the interval that follows current, capped at MaxInterval.
func (p Polling) next(current time.Duration) time.Duration {
	grown := float64(current) * p.Multiplier
	if grown >= float64(p.MaxInterval) {
		return p.MaxInterval
	}

	return time.Duration(grown)
}
