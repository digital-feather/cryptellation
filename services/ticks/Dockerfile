FROM golang:alpine

WORKDIR /app

# Code
COPY / .

# Build from code
WORKDIR /app/services/ticks
RUN go build ./cmd/ticks-server

ENV CGO_ENABLED 0

FROM alpine

# Install required running dependencies
RUN apk update \
    && apk add bash

# Install running files
COPY services/ticks/entrypoint.sh .
COPY --from=0 /app/services/ticks/ticks-server /usr/local/bin

ENTRYPOINT [ "/bin/bash", "/entrypoint.sh" ]
CMD [ "ticks-server" ]
