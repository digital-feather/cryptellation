#!/bin/bash

set -eo pipefail

# Wait cockroachdb
if [ ! -z ${COCKROACHDB_HOST+x} ]; then
  export PGHOST=$COCKROACHDB_HOST
  export PGPASSWORD=$COCKROACHDB_PASSWORD
  export PGPORT=$COCKROACHDB_PORT
  export PGUSER=$COCKROACHDB_USER
  export PGDATABASE=defaultdb

  until psql -c '\q' &> /dev/null; do
    echo "CockroachDB is unavailable - sleeping"
    sleep 3
  done
fi

# Launch tests
go test $(go list ./... | grep -e /adapters -e /service$ ) -coverprofile cover.out
