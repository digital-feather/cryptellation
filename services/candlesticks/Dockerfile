FROM golang:alpine

WORKDIR /app

# Dependencies
COPY clients/go ./clients/go
COPY internal/go ./internal/go

# Code
COPY services/candlesticks ./services/candlesticks/

# Build from code
WORKDIR /app/services/candlesticks
RUN go build ./cmd/candlesticks-server

ENV CGO_ENABLED 0

FROM alpine

# Install required running dependencies
RUN apk update \
    && apk add bash \
    && apk add postgresql-client

# Install running files
COPY services/candlesticks/entrypoint.sh .
COPY scripts/runtime /scripts/
COPY --from=0 /app/services/candlesticks/candlesticks-server /usr/local/bin

ENTRYPOINT [ "/bin/bash", "/entrypoint.sh" ]
CMD [ "candlesticks-server" ]
