package teecp

import "sync"

// Clients maintains a slice of receivers for teecp.
type Clients struct {
	mu        sync.Mutex
	receivers []Receiver
}

// Broadcast sends a message to every knwon receiver. If the receiver is no longer active,
// it is removed from the slice.
func (c *Clients) Broadcast(msg string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i := 0; i < len(c.receivers); i++ {
		r := c.receivers[i]

		if !r(msg) {
			// Replace the current receive with the last one in the slice and reset the index.
			// This allows for in-place replacement.
			c.receivers[i] = c.receivers[len(c.receivers)-1]
			i--
		}
	}
}

// Attach adds a receiver as a client.
func (c *Clients) Attach(receiver Receiver) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.receivers = append(c.receivers, receiver)
}

type Receiver func(msg string) bool
