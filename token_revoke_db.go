package token_revoke

import (
	"container/list"
	"fmt"
	"strconv"
	"sync"
	"time"
)

package main

import (
"container/list"
"fmt"
"sync"
"time"

"github.com/aws/aws-sdk-go/aws"
"github.com/aws/aws-sdk-go/aws/session"
"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Token struct {
	ID         string
	Expiration time.Time
}

func main() {
	// Initialize AWS session
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	// Specify the name of the DynamoDB table
	tableName := "my-tokens-table"

	// Specify the TTL attribute name
	ttlAttributeName := "expiration"

	// Specify the lease duration in seconds
	leaseDuration := 60

	// Initialize a list to hold the tokens
	tokenList := list.New()

	// Initialize a mutex to protect the list
	var tokenListMutex sync.Mutex

	// Start a goroutine to periodically check for expired tokens
	go func() {
		for {
			// Calculate the current time plus the lease duration
			expirationTime := time.Now().Add(time.Duration(leaseDuration) * time.Second)

			// Create a DynamoDB query input
			input := &dynamodb.QueryInput{
				TableName: aws.String(tableName),
				KeyConditions: map[string]*dynamodb.Condition{
					ttlAttributeName: {
						ComparisonOperator: aws.String("LT"),
						AttributeValueList: []*dynamodb.AttributeValue{
							{
								N: aws.String(fmt.Sprintf("%d", expirationTime.Unix())),
							},
						},
					},
				},
			}

			// Query the DynamoDB table for expired tokens
			result, err := svc.Query(input)
			if err != nil {
				fmt.Println("Error querying DynamoDB:", err)
				continue
			}

			// Add the expired tokens to the list
			tokenListMutex.Lock()
			for _, item := range result.Items {
				id := *item["id"].S
				expirationTimeUnix, err := item["expiration"].N()
				if err != nil {
					fmt.Println("Error parsing expiration time:", err)
					continue
				}
				expirationTimeUnixInt64, err := strconv.ParseInt(expirationTimeUnix, 10, 64)
				if err != nil {
					fmt.Println("Error parsing expiration time:", err)
					continue
				}
				expirationTime := time.Unix(expirationTimeUnixInt64, 0)
				tokenList.PushBack(Token{
					ID:         id,
					Expiration: expirationTime,
				})
			}
			tokenListMutex.Unlock()

			// Sleep for the lease duration before checking again
			time.Sleep(time.Duration(leaseDuration) * time.Second)
		}
	}()

	// Start a goroutine to revoke expired tokens
	go func() {
		for {
			// Get the current time


