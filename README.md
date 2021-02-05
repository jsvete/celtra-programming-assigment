# Celtra Programming Assigment

## Requrements

- [Git](https://git-scm.com/downloads) to pull code from this repository
- [Golang](https://golang.org/dl/) to run unit/integration tests for REST API and database
- [Docker](https://hub.docker.com/search?q=&type=edition&offering=community) to build Docker images
- [Docker Compose](https://docs.docker.com/compose/install/) to start all the services at once and the correct order 

## Prepare the environment

1. Clone the repository
```
git clone git@github.com:jsvete/celtra-programming-assigment.git
```
2. Build `cli` and `tracker` Docker images
```
docker build --build-arg SERVICE=cli -t cli:latest .
```
```
docker build --build-arg SERVICE=tracker -t tracker:latest .
```


### Start the services

1. Start `tracker` service, Redis and PostgreSQL containers
```
docker-compose up -d
```

2. Start `cli` service in a separate Docker container
```
docker run --rm -ti --network celtra-programming-assigment cli
```

### Cleanup
- Remove containers and custom network
```
docker-compose down
```

### Unit/Integration tests

Sadly I couldn't dockerize unit and integration tests due to some connection issues between containers.
Due to this, you can only run them locally via `go test`

Runs unit and integration tests with coverage:
```
go test -v -coverprofile cover.out ./...
```

Note:
Set environment variable `CGO_ENABLED` to `0` (e.g. `export CGO_ENABLED=0`) if you're having issues with Go reporting
```
# runtime/cgo
cgo: exec /missing-cc: fork/exec /missing-cc: no such file or directory
```
or a similar error.