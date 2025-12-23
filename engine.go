package rely

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/exeebit/rely/journal"
)

// WorkflowFunc is the function signature for a workflow.
type WorkflowFunc func(ctx Context, input ...interface{}) error

// Engine manages the durable execution.
type Engine struct {
	journal   journal.Journal
	workflows map[string]WorkflowFunc
}

// New creates a new Engine with the given journal backend.
func New(j journal.Journal) *Engine {
	return &Engine{
		journal:   j,
		workflows: make(map[string]WorkflowFunc),
	}
}

// Define registers a workflow definition.
func (e *Engine) Define(name string, fn WorkflowFunc) *Workflow {
	e.workflows[name] = fn
	return &Workflow{
		name:   name,
		engine: e,
	}
}

type Workflow struct {
	name   string
	engine *Engine
}

// Execute starts or follows a workflow execution.
func (w *Workflow) Execute(ctx context.Context, args ...interface{}) error {
	// 1. Replay History to build current state
	events, err := w.engine.journal.Read()
	if err != nil {
		return fmt.Errorf("failed to read journal: %w", err)
	}

	// 2. Create the Replay Context
	rctx := &replayContext{
		Context:      ctx,
		journal:      w.engine.journal,
		history:      make(map[string]*journal.Event),
		workflowName: w.name,
	}

	// Index history by step name
	for i := range events {
		ev := events[i]
		if ev.Type == journal.EventStepCompleted {
			rctx.history[ev.StepName] = &ev
		}
	}

	// 3. Log Workflow Started (omitted for MVP)

	// 4. Run the Workflow Function
	fn, ok := w.engine.workflows[w.name]
	if !ok {
		return fmt.Errorf("workflow %s not found", w.name)
	}

	return fn(rctx, args...)
}

// replayContext implements Context
type replayContext struct {
	context.Context
	journal      journal.Journal
	history      map[string]*journal.Event
	workflowName string
}

func (c *replayContext) Step(name string, fn func() (interface{}, error), opts ...StepOption) StepResult {
	// 1. Check if step already completed
	if event, ok := c.history[name]; ok {
		// Replay: Return cached result
		return &stepResult{
			payload: event.Payload,
			err:     nil,
		}
	}

	// 2. Not found? Execute it.
	val, err := fn()

	// 3. Handle Failure
	if err != nil {
		return &stepResult{err: err}
	}

	// 4. Handle Success: Serialize
	payload, err := json.Marshal(val)
	if err != nil {
		return &stepResult{err: fmt.Errorf("failed to marshal step result: %w", err)}
	}

	// 5. Persist to Journal
	je := journal.Event{
		Type:      journal.EventStepCompleted,
		Workflow:  c.workflowName,
		StepName:  name,
		Payload:   payload,
		Timestamp: 0, // Should be time.Now().Unix()
	}

	if err := c.journal.Append(je); err != nil {
		return &stepResult{err: fmt.Errorf("failed to write to journal: %w", err)}
	}

	return &stepResult{
		payload: payload,
		err:     nil,
	}
}
