package rely

import (
	"context"
)

// Context wraps standard context and adds durable execution capabilities.
type Context interface {
	context.Context

	// Step executes a unit of work durably.
	// name: Unique identifier for the step within the workflow.
	// fn: The function to execute. Must return (result, error).
	Step(name string, fn func() (interface{}, error), opts ...StepOption) StepResult
}

// StepResult allows retrieving the result of a step or its error.
type StepResult interface {
	// Result unmarshals the success value into target (pointer).
	Result(target interface{}) error
	// Err returns the error if the step failed.
	Err() error
}

type stepResult struct {
	err     error
	payload []byte
}

func (r *stepResult) Result(target interface{}) error {
	if r.err != nil {
		return r.err
	}
	if target == nil || r.payload == nil {
		return nil
	}
	return unmarshalResult(r.payload, target)
}

func (r *stepResult) Err() error {
	return r.err
}
