# Celtra Programming Assigment

## Requrements

- [Git](https://git-scm.com/downloads) to pull code from this repository
- [Golang](https://golang.org/dl/) to run integration tests for REST and database
- [Docker](https://hub.docker.com/search?q=&type=edition&offering=community) to build Docker images
- [Docker Compose](https://docs.docker.com/compose/install/) to start all the services at once and the correct order 

## Prepare the environment

1. Clone the repository
```
git clone git@github.com:jsvete/celtra-programming-assigment.git
```

2. Build small `cli` and `tracker` images
```
docker build --build-arg SERVICE=cli -t cli:latest .
```
```
docker build --build-arg SERVICE=tracker -t tracker:latest .
```

3. Start everything up
