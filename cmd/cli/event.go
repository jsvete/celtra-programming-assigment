package main

import (
	"celtra-programming-assigment/pkg/pubsub"
	"fmt"
)

func listenForEvents() {
	defer fmt.Printf("stopped listening\n")
	events := pubsub.Bus.Subscribe()

	for event := range events {
		if event.ID < 1 {
			fmt.Printf("<%s>: [error] %s\n", event.Timestamp.Format("2006-01-02 15:04:05:000"), event.Data)
		} else if _, ok := selectedIds[event.ID]; ok {
			fmt.Printf("<%s>: [%d]: %s\n", event.Timestamp.Format("2006-01-02 15:04:05:000"), event.ID, event.Data)
		}
	}
}
