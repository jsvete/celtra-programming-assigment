// Package pubsub contains code for publishing or subscribing data.
package pubsub

import (
	"time"
)

// Bus is an active messaging bus connection
var Bus PubSub

// Event struct wraps the data that is received when subscribing to an account event stream.
type Event struct {
	ID        int
	Timestamp time.Time
	Data      string
}

// PubSub interface represents the connection to the messaging bus
// and defines methods that can be implemented by various messaging providers.
//
// It can also be used to create a mocked implementation for testing purposes.
type PubSub interface {
	// Publish publishes the account's data to the "events" subscriptiono.
	Publish(accountID int, data string) error
	// Subscribe is used to subscribe to an "events" channel.
	//
	// Returns a channel where you can receive those events.
	Subscribe() chan *Event
}
