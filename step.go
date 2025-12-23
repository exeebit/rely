package rely

import (
	"encoding/json"
)

// StepOption allows configuring a step (e.g., retries).
type StepOption func(*StepConfig)

type StepConfig struct {
	MaxRetries int
}

func Retry(times int) StepOption {
	return func(c *StepConfig) {
		c.MaxRetries = times
	}
}

// Step represents a unit of execution.
type Step struct {
	Name string
	Fn   func() error
}

// ResultContainer is a helper to unmarshal step results.
type ResultContainer struct {
	Value interface{}
	Err   error
}

// unmarshalResult helper to decode payload into a target pointer.
func unmarshalResult(payload []byte, target interface{}) error {
	return json.Unmarshal(payload, target)
}
