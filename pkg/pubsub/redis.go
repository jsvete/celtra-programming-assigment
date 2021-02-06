// Package pubsub contains code for publishing or subscribing data.
package pubsub

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

var redisAddr string // REDIS_ADDR

// Redis struct is an implementation of PubSub interface 
// and is using a Redis client for publishing and subscribing.
type Redis struct {
	client *redis.Client
}

// NewRedis creates a new PubSub client that uses Redis for publishing and subscribing to events.
func NewRedis() error {
	redisAddr = os.Getenv("REDIS_ADDR")

	redisBus := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   0,
	})

	status := redisBus.Ping(context.Background())
	if status.Err() != nil {
		return status.Err()
	}

	Bus = &Redis{
		client: redisBus,
	}

	return nil
}

// Publish publishes the account's data to the Bus.
func (r *Redis) Publish(accountID int, data string) error {
	event := Event{
		ID:        accountID,
		Timestamp: time.Now().UTC(),
		Data:      data,
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return r.client.Publish(context.Background(), "events", string(eventData)).Err()
}

// Subscribe is used to subscribe to one or multiple accounts.
//
// Returns a channel where you can receive those events.
func (r *Redis) Subscribe() chan *Event {
	eventChan := make(chan *Event)

	sub := r.client.Subscribe(context.Background(), "events")

	go func() {
		msgChan := sub.Channel()

		for msg := range msgChan {
			payload := []byte(msg.Payload)

			event := &Event{}
			if err := json.Unmarshal(payload, event); err != nil {
				log.Warn().Msgf("error while deserializing event, sending error on event chan: %v", err)

				eventChan <- &Event{
					ID:        -1,
					Timestamp: time.Now().UTC(),
					Data:      err.Error(),
				}

				continue
			}

			eventChan <- event
		}
	}()

	return eventChan
}
