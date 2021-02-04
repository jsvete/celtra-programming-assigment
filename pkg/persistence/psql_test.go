// Package persistence contains database logic.
package persistence

import (
	"fmt"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

func TestMain(m *testing.M) {
	// start postgres-test docker container and connect to it
	pool, err := dockertest.NewPool("")
	if err != nil {
		panic(fmt.Sprintf("docker start: %v\n", err))
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13.1",
		Env: []string{
			"POSTGRES_USER=tracker",
			"POSTGRES_PASSWORD=tracker",
			"POSTGRES_DB=tracker",
		},
		Name: "postgres-test",
		ExposedPorts: []string{
			"5432",
		},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: "5432"},
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
	resource.Expire(60)

	os.Setenv("DB_NAME", "tracker")
	os.Setenv("DB_USER", "tracker")
	os.Setenv("DB_PWD", "tracker")
	os.Setenv("DB_ADDR", "localhost")
	os.Setenv("DB_PORT", "5432")

	if err = pool.Retry(func() error {
		if err := NewPostgres(); err != nil {
			fmt.Printf("error connecting to db: %v\n", err)
			return err
		}

		return nil
	}); err != nil {
		panic(fmt.Sprintf("couldn't connect to postgres container: %v", err))
	}

	// run tests
	m.Run()

	if err = pool.Purge(resource); err != nil {
		panic(fmt.Sprintf("container stop: %v", err))
	}
}

func Test_Account(t *testing.T) {
	// test insert
	account, err := DB.CreateAccount("test account", false)
	if err != nil {
		t.Fatalf("failed to create account: %v", err)
	}

	// we already populated the database with 1000 records
	if account.ID != 1001 {
		t.Fatalf("account id, expected %d, was %d", 1001, account.ID)
	}

	if account.Name != "test account" {
		t.Fatalf("account name, expected %s, was %s", "test account", account.Name)
	}

	if account.IsActive {
		t.Fatalf("account isActive, expected %t, was %t", false, account.IsActive)
	}

	// test select
	account, err = DB.GetAccount(account.ID)
	if err != nil {
		t.Fatalf("failed to get account: %v", err)
	}

	if account.ID != 1001 {
		t.Fatalf("account id, expected %d, was %d", 1001, account.ID)
	}

	if account.Name != "test account" {
		t.Fatalf("account name, expected %s, was %s", "test account", account.Name)
	}

	if account.IsActive {
		t.Fatalf("account isActive, expected %t, was %t", false, account.IsActive)
	}

	// test isActive check
	isActive, err := DB.IsActiveAccount(account.ID)
	if err != nil {
		t.Fatalf("failed to get account: %v", err)
	}

	if isActive {
		t.Fatalf("account isActive, expected %t, was %t", false, isActive)
	}

	// test get active account
	isActive, err = DB.IsActiveAccount(1)
	if err != nil {
		t.Fatalf("failed to get account: %v", err)
	}

	if !isActive {
		t.Fatalf("account isActive, expected %t, was %t", false, isActive)
	}

	// test get unknown account
	isActive, err = DB.IsActiveAccount(9999)
	if err == nil {
		t.Fatalf("IsActiveAccount(9999) should have returned an error")
	}

	account, err = DB.GetAccount(9999)
	if err == nil {
		t.Fatalf("GetAccount(9999) should have returned an error")
	}
}
