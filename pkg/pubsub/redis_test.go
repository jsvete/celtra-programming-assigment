// Package pubsub contains code for publishing or subscribing data.
package pubsub

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

func TestMain(m *testing.M) {
	// start redis-test docker container and connect to it
	pool, err := dockertest.NewPool("")
	if err != nil {
		panic(fmt.Sprintf("docker start: %v\n", err))
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "redis",
		Tag:        "6.0.10-alpine3.12",
		Name:       "redis-test",
		ExposedPorts: []string{
			"6379",
		},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"6379": {
				{HostIP: "0.0.0.0", HostPort: "6379"},
			},
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		panic(fmt.Sprintf("container start: %v", err))
	}
	//	resource.Expire(60)

	os.Setenv("REDIS_ADDR", "localhost:6379")

	if err = pool.Retry(func() error {
		if err := NewRedis(); err != nil {
			fmt.Printf("error connecting to redis: %v\n", err)
			return err
		}

		return nil
	}); err != nil {
		panic(fmt.Sprintf("couldn't connect to redis container: %v", err))
	}

	// run tests
	m.Run()

	if err = pool.Purge(resource); err != nil {
		panic(fmt.Sprintf("container stop: %v", err))
	}
}

func Test_PubSub(t *testing.T) {
	eventChan := Bus.Subscribe()

	// need to wait a bit for Redis to register the subscription before we can publish or the event is lost
	time.Sleep(3 * time.Second)

	err := Bus.Publish(1, "test data")
	if err != nil {
		fmt.Printf("failed to publish\n: %v", err)
	}

	select {
	case event := <-eventChan:
		if event == nil {
			t.Fatalf("expected a published event")
		}
		if event.ID != 1 {
			t.Fatalf("expected %d, got %d", 1, event.ID)
		}
		if event.Data != "test data" {
			t.Fatalf("expected %s, got %s", "test data", event.Data)
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("timed out")
	}

}
