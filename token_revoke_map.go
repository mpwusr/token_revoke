package main

import (
	"fmt"
	"time"
)

func main() {
	// Start a goroutine to revoke tokens
	go revokeTokens()

	// Wait for some time to simulate token activity
	time.Sleep(10 * time.Second)
}

func revokeTokens() {
	// Map to store tokens and their TTLs
	tokens := make(map[string]time.Time)

	for {
		// Check TTLs of tokens and revoke expired ones
		for token, ttl := range tokens {
			if time.Now().After(ttl) {
				delete(tokens, token)
				fmt.Printf("Token %s has been revoked\n", token)
			}
		}

		// Wait for some time before checking again
		time.Sleep(1 * time.Second)
	}
}
