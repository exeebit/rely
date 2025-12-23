package journal

// EventType represents the type of event in the journal.
type EventType string

const (
	EventWorkflowStarted   EventType = "WorkflowStarted"
	EventStepStarted       EventType = "StepStarted"
	EventStepCompleted     EventType = "StepCompleted"
	EventWorkflowCompleted EventType = "WorkflowCompleted"
	EventWorkflowFailed    EventType = "WorkflowFailed"
)

// Event is a single immutable record in the journal.
type Event struct {
	ID        int       `json:"id"`
	Type      EventType `json:"type"`
	Workflow  string    `json:"workflow"` // Workflow Definition Name
	StepName  string    `json:"step_name,omitempty"`
	Payload   []byte    `json:"payload,omitempty"`   // Serialized input/output
	Error     string    `json:"error,omitempty"`     // Error message if failed
	Timestamp int64     `json:"timestamp"`
}

// Journal defines the storage interface for the execution engine.
type Journal interface {
	// Append adds a new event to the journal.
	Append(event Event) error

	// Read returns all events for a given workflow instance history.
	// In a real system, this might take an InstanceID.
	// For this MVP, we assume a single-threaded access or single instance per journal for simplicity,
	// or we can add InstanceID to the Event struct later.
	// For now, let's treat the Journal as a stream for *one* execution or filter by workflow.
	Read() ([]Event, error)
	
	// Close cleans up resources working with the journal.
	Close() error
}
