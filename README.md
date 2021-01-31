# Celtra Programming Assigment

## Requrements

- [Git](https://git-scm.com/downloads) to pull code from this repository
- [Docker](https://hub.docker.com/search?q=&type=edition&offering=community) to build and deploy code inside Docker container

## Prepare the environment

1. Clone the repository
```
git clone git@github.com:jsvete/celtra-programming-assigment.git
```

2. Build `godev` image
```
docker build -t godev .
```

2. Chech if everything works; you should see the directory structure printed out
```
docker run --rm -v <FULL_LOCAL_PATH_TO_CLONED_REPOSITORY>:/code godev  ls -la
```
Example (on linux):
```
docker run --rm -v $PWD:/code godev  ls -la
```