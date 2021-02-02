package main

import (
	"celtra-programming-assigment/cmd/tracker/persistence"
	"celtra-programming-assigment/cmd/tracker/rest"
	"celtra-programming-assigment/pkg/pubsub"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func main() {
	if err := persistence.NewPostgres(); err != nil {
		panic(err)
	}

	if err := pubsub.NewRedis(); err != nil {
		panic(err)
	}

	router := httprouter.New()
	router.HandlerFunc("GET", "/:accountId", rest.HandleGet)
	router.HandlerFunc("POST", "/", rest.HandlePost)
	router.HandlerFunc("PUT", "/:accountId", rest.HandlePut)

	if err := http.ListenAndServe(":8081", router); err != nil {
		panic(err)
	}
}
