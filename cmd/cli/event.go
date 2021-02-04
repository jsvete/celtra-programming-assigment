package main

import (
	"celtra-programming-assigment/pkg/pubsub"
	"fmt"
)

func listenForEvents() {
	defer fmt.Printf("stopped listening\n")

	for event := range pubsub.Bus.Subscribe() {
		if event.ID < 1 {
			fmt.Printf("<%s>: error: %s", event.Timestamp.Format("2006-01-02 15:04:05:000"), event.Data)
		} else if _, ok := selectedIds[event.ID]; ok {
			fmt.Printf("<%s>: [%d]: %s", event.Timestamp.Format("2006-01-02 15:04:05:000"), event.ID, event.Data)
		}
	}
}
