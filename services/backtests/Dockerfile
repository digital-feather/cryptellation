FROM golang:alpine

WORKDIR /app

# Dependencies
COPY clients/go ./clients/go
COPY internal/go ./internal/go
COPY services/candlesticks ./services/candlesticks/

# Code
COPY services/backtests ./services/backtests/

# Build from code
WORKDIR /app/services/backtests
RUN go build ./cmd/backtests-server

ENV CGO_ENABLED 0

FROM alpine

# Install required running dependencies
RUN apk update \
    && apk add bash \
    && apk add postgresql-client

# Install running files
COPY services/backtests/entrypoint.sh .
COPY scripts/runtime /scripts/
COPY --from=0 /app/services/backtests/backtests-server /usr/local/bin

ENTRYPOINT [ "/bin/bash", "/entrypoint.sh" ]
CMD [ "backtests-server" ]
