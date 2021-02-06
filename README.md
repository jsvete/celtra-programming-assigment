# Celtra Programming Assigment

This is my implementation of the Celtra's backend programming assigment. It was written in Go and is packaged together with Docker to form two main containers. It uses Docker Compose to set up and run the main `tracking` service and its dependencies (database, pubsub, reverse proxy/load balancer). A normal Docker `run` command is used to start the `cli` client.

## Main components of the system

The application has a couple of different components:
- `tracker` service,
- `cli` client,
- PostgreSQL database,
- Redis for messaging (publish/subscribe),
- [nginx-proxy](https://github.com/nginx-proxy/nginx-proxy) for acting as areverse proxy and load balancing (round robin).

`tracker` service is the core part of the system and it was designed with scaling in mind. If you uncomment the additional services in the `docker-compose.yml` file and run it again, it will deploy 3 different instances of the `tracker` service. `nginx-proxy` will take care of the registration of the new containers automatically. It also exposes a single port (`8080`) where all the REST API requests forwarded by the `nginx-proxy` are made.

The `cli` client was made to be fault tolerant. If you kill the rest of the system (e.g. `docker-compose stop`), it will print out an error but it won't crash. After the restart ti will continue to receive the events.

Redis is used as the messaging pipeline. I chose it because it is simple to deploy and use. Its [Go client](https://github.com/go-redis/redis) is also very easy to use and is fault tolerant. It will check for downed connections and will resubscribe when it detects that the Redis is back online.

PostgreSQL was chosen for database layer. This was mostly due to being used to it since I've been working with it for couple of years. At the startup the `tracker` service will create the `account` table and insert 1000 records with randomly generated data into it.

While doing some research on how to make `tracker` service more scalable, I found `nginx-proxy`. Its main feature is that it can automatically update its configuration when it detects that a new Docker container was deployed. With correct configuration it also works as a simple load balancer.


## Requrements

- [Git](https://git-scm.com/downloads) to pull code from this repository
- [Golang 1.15.x](https://golang.org/dl/) to run unit/integration tests for REST API and database
- [Docker](https://hub.docker.com/search?q=&type=edition&offering=community) to build Docker images
- [Docker Compose](https://docs.docker.com/compose/install/) to start all the services

## Setup

### Prepare the environment

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

1. Start `tracker` service, Redis, PostgreSQL and nginx-proxy containers
```
docker-compose up -d
```

2. Start `cli` service in a separate Docker container
```
docker run --rm -ti --network celtra-programming-assigment cli
```

## Cleanup
Remove containers and custom network
```
docker-compose down
```


## Unit/Integration tests

Sadly I couldn't dockerize unit and integration tests due to some connection issues between containers.
Because of this, you can only run them locally via `go test`

Integration tests check the SQL queries and Publish/Subscribe mechanism against the real service without mocking the API. This mean that the test will bootstrap a PostgreSQL/Redis container and run the test.

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

## CLI client
### Account selection
After you start it, you will se the follwoing prompt:
```
Connecting to Redis@redis:6379
options: [accounts events]
>
```
After typing in `accounts`, you will se already selected accounts and input guidelines:
```
>accounts
 already selected accounts: []
 input space separated account IDs
accounts >
```
When you input the space-separated account IDs, you will be returned to the main prompt where you can modify your ID selection or continue to event subscription:
```
accounts >1 2 3
 current selected accounts: [1 2 3]
options: [accounts events]
>
```
If you go back to `accounts` you can select additional accounts:
```
>accounts
 already selected accounts: [1 2 3]
 input space separated account IDs
accounts >
```
If you don't with to select any more account IDs, just press `<Enter>` and you will be returned to the main prompt:
```
accounts >
 "" should be an (positive) number
 current selected accounts: [1 2 3]
options: [accounts events]
>
```
### Event listening
In the main prompt, type in `events`
```
>events
listening for events from: [1 2 3]
```
You will get a message that the client has started listening to the selected events.

When the client will receive the events, it will output them to the terminal:
```
<2021-02-06 17:35:30:000>: [1]: "test data" [290ad619a440]
<2021-02-06 17:35:35:000>: [2]: "test data" [290ad619a440]
<2021-02-06 17:35:39:000>: [3]: "test data" [290ad619a440]
```
Structure of the message is: `<UTC_TIMESTAMP>: [ACCOUNT_ID]: "RECEIVED_DATA" [SERVICE_HOSTNAME]`

`SERVICE_HOSTNAME` can be used to identify which instance of the `tracker` service sent the event.

If you kill the rest of the system, you should se an error message in the terminal. But don't be discuraged, because once you restart the system, you should againg start receiving events without restarting the client:
```
redis: 2021/02/06 17:39:02 pubsub.go:168: redis: discarding bad PubSub connection: EOF
<2021-02-06 17:40:56:000>: [3]: "test data" [7f76a48100a6]
<2021-02-06 17:41:04:000>: [1]: "test data" [7f76a48100a6]
<2021-02-06 17:41:08:000>: [2]: "test data" [7f76a48100a6]
```

To exit the application, use `Ctrl+C`. This will also remove the container so no additional cleanup is required.
## REST API
### Fetch account information:
```
GET: localhost:8080/<accountID>
```
Response:
```
{
    "ID": 1000,
    "Name": "4cc2f6d0bed451c2e0c0b6b6aa21bfca",
    "IsActive": true
}
```
### Send an event for a specific account ID.
```
PUT: localhost:8080/<accountID>?data="<data>"
```
### Create a new account
```
POST: localhost:8080/
Content-Type: application/json
```
Body:
```
{
    "Name":"ACCOUNT_NAME", 
    "IsActive": true/false
}
```
### Get rate counter information
```
GET: localhost:8080
```
Response:
```
{
    "Rate":997
}
```

