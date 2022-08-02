FROM golang:alpine

# Install required running dependencies
RUN apk update \
    && apk add bash \
    && apk add postgresql-client

# Do not use gcc
ENV CGO_ENABLED 0

# Set testing dir/entrypoint
WORKDIR /code
ENTRYPOINT [ "/bin/bash", "/code/tests/integration/entrypoint.sh" ]
