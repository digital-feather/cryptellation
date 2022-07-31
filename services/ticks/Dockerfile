FROM golang:alpine

WORKDIR /app

# Dependencies
COPY clients/go ./clients/go
COPY internal/go ./internal/go
COPY services/candlesticks ./services/candlesticks/

# Code
COPY services/ticks ./services/ticks/

# Build from code
WORKDIR /app/services/ticks
RUN go build ./cmd/ticks-server

ENV CGO_ENABLED 0

FROM alpine

# Install required running dependencies
RUN apk update \
    && apk add bash \
    && apk add postgresql-client

# Install running files
COPY services/ticks/entrypoint.sh .
COPY scripts/runtime /scripts/
COPY --from=0 /app/services/ticks/ticks-server /usr/local/bin

ENTRYPOINT [ "/bin/bash", "/entrypoint.sh" ]
CMD [ "ticks-server" ]
