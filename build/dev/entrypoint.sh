#!/bin/bash

set -eo pipefail

# Wait for cockroachdb if required
if [ ! -z ${COCKROACHDB_HOST+x} ]; then
  export PGHOST=$COCKROACHDB_HOST
  export PGPASSWORD=$COCKROACHDB_PASSWORD
  export PGPORT=$COCKROACHDB_PORT
  export PGUSER=$COCKROACHDB_USER
  export PGDATABASE=defaultdb

  until psql -c '\q'; do
    >&2 echo "CockroachDB is unavailable - sleeping"
    sleep 1
  done

  for f in configs/sql/*.sql; do
    filename=$(basename -- "$f")
    export PGDATABASE="${filename%.*}"

    echo "# Loading '$PGDATABASE' database from $f..."
    PGOPTIONS='--client-min-messages=error' psql -f $f
  done
fi

reflex -s -r '(\.go$|go\.mod)' -- $@