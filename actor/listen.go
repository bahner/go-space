package actor

import (
	"fmt"
)

func (a *Actor) Listen(outputChannel chan<- string) error {
	// Subscribe to Inbox topic
	inboxSub, err := a.Inbox.Subscribe()
	if err != nil {
		return fmt.Errorf("failed to subscribe to Inbox topic: %v", err)
	}
	defer inboxSub.Cancel()

	// Subscribe to Space topic
	spaceSub, err := a.Space.Subscribe()
	if err != nil {
		return fmt.Errorf("failed to subscribe to Space topic: %v", err)
	}
	defer spaceSub.Cancel()

	// Start a goroutine for Inbox subscription
	go a.handleInboxSubscription(inboxSub)

	// Start a goroutine for Space subscription
	// Assuming you have a similar function for Space
	go a.handleSpaceSubscription(spaceSub)

	// Start a goroutine for REPL
	go a.HandleREPL(outputChannel)

	// Wait for context cancellation (or other exit conditions)
	<-a.Ctx.Done()
	return a.Ctx.Err()
}
