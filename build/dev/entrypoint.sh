#!/bin/bash

set -eo pipefail

/docker/setup-db.sh

reflex -s -r '(\.go$|go\.mod)' -- $@
