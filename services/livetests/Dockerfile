FROM golang:alpine

WORKDIR /app

# Dependencies
COPY clients/go ./clients/go
COPY internal/go ./internal/go
COPY services/backtests ./services/backtests/
COPY services/candlesticks ./services/candlesticks/

# Code
COPY services/livetests ./services/livetests/

# Build from code
WORKDIR /app/services/livetests
RUN go build ./cmd/livetests-server

ENV CGO_ENABLED 0

FROM alpine

# Install required running dependencies
RUN apk update \
    && apk add bash \
    && apk add postgresql-client

# Install running files
COPY services/livetests/entrypoint.sh .
COPY scripts/runtime /scripts/
COPY --from=0 /app/services/livetests/livetests-server /usr/local/bin

ENTRYPOINT [ "/bin/bash", "/entrypoint.sh" ]
CMD [ "livetests-server" ]
