#!/bin/bash

set -eo pipefail

/docker/wait-required-services.sh

go test $(go list ./... | grep -e /adapters -e /service$ ) -coverprofile cover.out
