FROM golang:alpine

RUN go get github.com/cespare/reflex

RUN apk update \
    && apk add inotify-tools bash \
    && apk add postgresql-client

COPY build/dev/ /docker

VOLUME /code
WORKDIR /code

ENV CGO_ENABLED 0

ENTRYPOINT [ "/bin/bash", "/docker/entrypoint.sh" ]