package journal

import "sync"

// MemoryJournal is a simple thread-safe in-memory journal.
type MemoryJournal struct {
	mu     sync.Mutex
	events []Event
}

// NewMemoryJournal creates a new empty in-memory journal.
func NewMemoryJournal() *MemoryJournal {
	return &MemoryJournal{
		events: make([]Event, 0),
	}
}

// Append adds an event to the history.
func (m *MemoryJournal) Append(event Event) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Auto-increment ID based on position
	event.ID = len(m.events) + 1
	m.events = append(m.events, event)
	return nil
}

// Read returns a copy of the current history.
func (m *MemoryJournal) Read() ([]Event, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Return a copy to avoid race conditions if the caller modifies it (though they shouldn't)
	dst := make([]Event, len(m.events))
	copy(dst, m.events)
	return dst, nil
}

func (m *MemoryJournal) Close() error {
	return nil
}
