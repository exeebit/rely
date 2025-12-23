package main

import (
	"context"
	"fmt"
	"log"

	"rely"
	"github.com/<your-username>/rely/journal"
)

func main() {
	// 1. Create In-Memory Journal
	// In a real app, you'd use a File or DB journal.
	j := journal.NewMemoryJournal()

	engine := rely.New(j)

	// 2. Define Workflow
	workflow := engine.Define("UserOnboarding", func(ctx rely.Context, args ...interface{}) error {
		emailStr := args[0].(string)
		fmt.Printf("Workflow Started for %s\n", emailStr)

		// Step 1: Charge
		var chargeID string
		if err := ctx.Step("ChargeCard", func() (interface{}, error) {
			fmt.Println(" -> EXEC: Charging Card...")
			return "ch_12345", nil
		}).Result(&chargeID); err != nil {
			return err
		}
		fmt.Printf("   [Done] Charge ID: %s\n", chargeID)

		// Step 2: Provision
		var serverID string
		if err := ctx.Step("ProvisionServer", func() (interface{}, error) {
			fmt.Println(" -> EXEC: Provisioning Server...")
			return "srv_9999", nil
		}).Result(&serverID); err != nil {
			return err
		}
		fmt.Printf("   [Done] Server ID: %s\n", serverID)

		// Step 3: Send Email
		if err := ctx.Step("SendEmail", func() (interface{}, error) {
			fmt.Printf(" -> EXEC: Sending Email to %s about server %s...\n", emailStr, serverID)
			return nil, nil // No return value
		}).Err(); err != nil {
			return err
		}
		fmt.Println("   [Done] Email Sent")

		return nil
	})

	// 3. Execute
	ctx := context.Background()
	fmt.Println("--- RUN 1 (Fresh) ---")
	if err := workflow.Execute(ctx, "user@example.com"); err != nil {
		log.Fatal(err)
	}

	// 4. Simulate Crash & Restart
	// We re-use the SAME journal instance to simulate persistence.
	fmt.Println("\n--- RUN 2 (Replay) ---")
	// The workflow should NOT print "-> EXEC" lines again, but should print "[Done]" lines.

	// Note: In a real "Crash", the memory journal would be lost.
	// But `rely` relies on the journal being persistent (e.g. file).
	// Since we used MemoryJournal here, we are simulating "restart with restored state".

	if err := workflow.Execute(ctx, "user@example.com"); err != nil {
		log.Fatal(err)
	}
}
