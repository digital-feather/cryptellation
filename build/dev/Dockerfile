FROM golang:alpine

RUN go install github.com/cespare/reflex@latest

RUN apk update \
    && apk add inotify-tools bash \
    && apk add postgresql-client

COPY build/dev/ /docker

VOLUME /code
WORKDIR /code

COPY go.mod go.sum ./
RUN go mod download

ENV CGO_ENABLED 0

ENTRYPOINT [ "/bin/bash", "/docker/entrypoint.sh" ]