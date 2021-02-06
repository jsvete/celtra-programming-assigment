package main

import (
	"celtra-programming-assigment/cmd/tracker/rest"
	"celtra-programming-assigment/pkg/persistence"
	"celtra-programming-assigment/pkg/pubsub"
	"net/http"

	"github.com/rs/zerolog"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func main() {
	// init persistence layer
	if err := persistence.NewPostgres(); err != nil {
		panic(err)
	}

	// init pubsub
	if err := pubsub.NewRedis(); err != nil {
		panic(err)
	}

	// init REST API
	if err := http.ListenAndServe(":8080", rest.CreateRouter()); err != nil {
		panic(err)
	}
}
