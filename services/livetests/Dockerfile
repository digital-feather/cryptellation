FROM golang:alpine

WORKDIR /app

# Code
COPY / .

# Build from code
WORKDIR /app/services/livetests
RUN go build ./cmd/livetests-server

ENV CGO_ENABLED 0

FROM alpine

# Install required running dependencies
RUN apk update \
    && apk add bash

# Install running files
COPY services/livetests/entrypoint.sh .
COPY --from=0 /app/services/livetests/livetests-server /usr/local/bin

ENTRYPOINT [ "/bin/bash", "/entrypoint.sh" ]
CMD [ "livetests-server" ]
