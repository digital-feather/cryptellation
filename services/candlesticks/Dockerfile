FROM golang:alpine

WORKDIR /app

# Code
COPY / .

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
COPY services/candlesticks/scripts /scripts/
COPY --from=0 /app/services/candlesticks/candlesticks-server /usr/local/bin

ENTRYPOINT [ "/bin/bash", "/entrypoint.sh" ]
CMD [ "candlesticks-server" ]
