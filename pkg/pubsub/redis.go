// Package pubsub contains code for publishing or subscribing data.
package pubsub

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

var redisAddr string // REDIS_ADDR

// Redis struct is an implementation of PubSub interface
// and wraps Redis client implementation plus it holds subscription callback.
type Redis struct {
	client *redis.Client
}

// NewRedis creates a new PubSub client that uses Redis for publishing and subscribing to events.
func NewRedis() error {
	redisBus := redis.NewClient(&redis.Options{
		Addr: redisAddr,
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
		timestamp: time.Now().UTC(),
		data:      data,
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return r.client.Publish(context.Background(), strconv.Itoa(accountID), string(eventData)).Err()
}

// Subscribe is used to subscribe to one or multiple accounts.
//
// Returns a channel where you can receive those events.
func (r *Redis) Subscribe(accountID ...int) (chan *Event, error) {
	eventChan := make(chan *Event)

	if len(accountID) < 1 {
		return nil, errors.New("need to specify one or more account IDs to subscribe to")
	}

	subscriptions := []string{}
	for _, id := range accountID {
		subscriptions = append(subscriptions, strconv.Itoa(id))
	}

	sub := r.client.Subscribe(context.Background(), subscriptions...)

	go func() {
		for msg := range sub.Channel() {
			payload := []byte(msg.Payload)

			event := &Event{}
			if err := json.Unmarshal(payload, event); err != nil {
				log.Warn().Msgf("error while deserializing event, sending error on event chan: %v", err)

				eventChan <- &Event{
					timestamp: time.Now().UTC(),
					data:      err.Error(),
				}

				continue
			}

			eventChan <- event
		}
	}()

	return eventChan, nil
}

func init() {
	redisAddr = os.Getenv("REDIS_ADDR")
}
