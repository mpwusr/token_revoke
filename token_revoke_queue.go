package main

import (
	"container/list"
	"fmt"
	"time"
)

package main

import (
"container/list"
"fmt"
"time"
)

type Token struct {
	id        int
	createdAt time.Time
}

func main() {
	// Create a queue to hold the access tokens
	queue := list.New()

	// Add some tokens to the queue
	for i := 0; i < 5; i++ {
		queue.PushBack(&Token{id: i, createdAt: time.Now()})
	}

	// Revoke tokens based on TTL
	for queue.Len() > 0 {
		// Get the oldest token from the queue
		oldestToken := queue.Front().Value.(*Token)

		// Calculate the age of the token
		age := time.Since(oldestToken.createdAt)

		// Check if the token has expired (TTL is 10 seconds in this example)
		if age.Seconds() > 10 {
			fmt.Printf("Revoking token with id %d\n", oldestToken.id)
			queue.Remove(queue.Front())
		} else {
			// If the token hasn't expired, exit the loop (because the tokens are ordered by age in the queue)
			break
		}
	}
}

