FROM golang:alpine

WORKDIR /app

# Dependencies
COPY clients/go ./clients/go
COPY internal/go ./internal/go

# Code
COPY services/exchanges ./services/exchanges/

# Build from code
WORKDIR /app/services/exchanges
RUN go build ./cmd/exchanges-server

ENV CGO_ENABLED 0

FROM alpine

# Install required running dependencies
RUN apk update \
    && apk add bash \
    && apk add postgresql-client

# Install running files
COPY services/exchanges/entrypoint.sh .
COPY scripts/runtime /scripts/
COPY --from=0 /app/services/exchanges/exchanges-server /usr/local/bin

ENTRYPOINT [ "/bin/bash", "/entrypoint.sh" ]
CMD [ "exchanges-server" ]
