FROM golang:1.18.5-alpine3.15 AS builder

RUN apk update && \
    apk add ca-certificates git curl gcc musl-dev build-base autoconf automake libtool

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

FROM builder as dev 

WORKDIR /app

CMD ["air"]

