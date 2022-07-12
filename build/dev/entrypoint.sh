#!/bin/bash

set -eo pipefail

/docker/wait-required-services.sh

reflex -s -r '(\.go$|go\.mod)' -- $@
