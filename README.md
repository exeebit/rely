# rely

`rely` is a lightweight, durable execution engine for Go. It allows you to define workflows that can survive crashes and restarts by replaying history from a journal.

## Installation

```bash
go get github.com/exeebit/rely
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/exeebit/rely"
	"github.com/exeebit/rely/journal"
)

func main() {
	// 1. Initialize Journal
	j := journal.NewMemoryJournal()
	engine := rely.New(j)

	// 2. Define Workflow
	workflow := engine.Define("MyWorkflow", func(ctx rely.Context, args ...interface{}) error {
		var result string
		// Define a durable step
		if err := ctx.Step("Step1", func() (interface{}, error) {
			return "done", nil
		}).Result(&result); err != nil {
			return err
		}
		fmt.Println("Result:", result)
		return nil
	})

	// 3. Execute
	if err := workflow.Execute(context.Background(), "arg1"); err != nil {
		log.Fatal(err)
	}
}
```

## Features

- **Durable Steps**: Steps are executed once and their results are persisted.
- **Journaling**: Pluggable journal backend (currently supports in-memory).
- **Replay**: Automatically restores state on restart.
