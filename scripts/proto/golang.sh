#!/bin/bash
set -e

function generate {
  local readonly SERVICE_NAME=$1
  local readonly OUTPUT_DIR=internal/controllers/grpc/genproto/$SERVICE_NAME

  mkdir -p $OUTPUT_DIR

  protoc \
    --proto_path=api/protobuf \
    --go_out=$OUTPUT_DIR \
    --go_opt=paths=source_relative \
    --go-grpc_opt=require_unimplemented_servers=false \
    --go-grpc_out=$OUTPUT_DIR \
    --go-grpc_opt=paths=source_relative \
    "api/protobuf/$SERVICE_NAME.proto"
}

for SERVICE_NAME in "$@"; do
  generate $SERVICE_NAME
done
