FROM golang:alpine

WORKDIR /app

# Code
COPY / .

# Build from code
WORKDIR /app/services/backtests
RUN go build ./cmd/backtests-server

ENV CGO_ENABLED 0

FROM alpine

# Install required running dependencies
RUN apk update \
    && apk add bash

# Install running files
COPY services/backtests/entrypoint.sh .
COPY --from=0 /app/services/backtests/backtests-server /usr/local/bin

ENTRYPOINT [ "/bin/bash", "/entrypoint.sh" ]
CMD [ "backtests-server" ]
