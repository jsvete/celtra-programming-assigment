package main

import (
	"celtra-programming-assigment/pkg/persistence"
	"celtra-programming-assigment/pkg/pubsub"
	"celtra-programming-assigment/cmd/tracker/rest"
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
	if err := http.ListenAndServe(":8081", rest.CreateRouter()); err != nil {
		panic(err)
	}
}
