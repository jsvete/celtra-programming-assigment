package main

import (
	"celtra-programming-assigment/pkg/pubsub"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/peterh/liner"
)

var (
	redisAddr = flag.String("addr", "redis", "Redis address")
	redisPort = flag.String("port", "6379", "Redis port")

	commands = []string{"accounts", "events"}
)

func main() {
	flag.Parse()

	addr := *redisAddr + ":" + *redisPort
	os.Setenv("REDIS_ADDR", addr)
	fmt.Printf("Connecting to Redis@%s\n", addr)

	if err := pubsub.NewRedis(); err != nil {
		panic(err)
	}

	cli := liner.NewLiner()
	defer cli.Close()

	cli.SetCtrlCAborts(true)

	cli.SetCompleter(func(line string) (c []string) {
		for _, command := range commands {
			if strings.HasPrefix(command, strings.ToLower(line)) {
				c = append(c, command)
			}
		}

		return
	})

	for {
		fmt.Printf("options: %s\n", commands)
		command, err := cli.Prompt(">")
		if err != nil {
			if err == liner.ErrPromptAborted {
				fmt.Println("Goodbye!")
				break
			} else {
				fmt.Printf("Error reading line: %v\n", err)
			}
		}

		switch command {
		case "accounts":
			fmt.Printf(" already selected accounts: %d\n", selectedAccounts())
			fmt.Printf(" input space separated account IDs\n")
			accounts, err := cli.Prompt("accounts >")
			if err != nil {
				if err == liner.ErrPromptAborted {
					break
				} else {
					fmt.Printf("Error reading line: %v\n", err)
				}
			}

			fmt.Printf(" current selected accounts: %d\n", selectAccounts(strings.Split(accounts, " ")...))
		case "events":
			if len(selectedIds) < 1 {
				fmt.Printf(" no account IDs selected\n")
				break
			}

			fmt.Printf("listening for events from: %d\n", selectedAccounts())

			listenForEvents()
		default:
			fmt.Printf("unrecognized command: %s\n", command)
		}
	}
}
