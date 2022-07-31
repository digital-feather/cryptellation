#!/bin/bash

set -eo pipefail

/scripts/wait-cockroachdb.sh

$@
