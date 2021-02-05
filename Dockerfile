FROM golang:1.15.7-alpine3.13 AS builder

# SERVICE defines a directory name inside /code/cmd that we want to build (e.g. tracker or cli)
ARG SERVICE

WORKDIR $GOPATH/src/celtra-programming-assigment/

# use go modules
COPY go.mod .

ENV GO111MODULE=on
RUN go mod download
RUN go mod verify

# copy the code
COPY . .

RUN mkdir -p build

# build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' -a \
    -o /go/bin/${SERVICE}  celtra-programming-assigment/cmd/${SERVICE}

# build the image (can't use scratch because we need /bin/sh for entrypoint)
# with scratch the image could be half the size but we would lose the ability to parametrize the build
FROM alpine:3.13

# SERVICE defines a directory name inside /code/cmd that we want to build (e.g. tracker or cli)
ARG SERVICE

# copy the built binary from the previous step
COPY --from=builder /go/bin/${SERVICE} /go/bin/${SERVICE}

# we need to set an environment variable since we can't use docker build arguments in the entrypoint
ENV BINARY=${SERVICE}

# start the binary from shell so we can pass the $BINARY environment variable to it
ENTRYPOINT [ "sh", "-c", "/go/bin/$BINARY" ]