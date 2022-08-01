#!/bin/bash
set -e

function generate {
  local readonly SERVICE_NAME=$1
  local readonly SERVICE_PATH=services/$SERVICE_NAME

  local readonly PROTO_PATH=$SERVICE_PATH/api
  local readonly INPUT_FILE=$PROTO_PATH/$SERVICE_NAME.proto
  local readonly OUTPUT_DIR=$SERVICE_PATH/pkg/client/proto

  mkdir -p $OUTPUT_DIR

  protoc \
    --proto_path=$PROTO_PATH \
    --go_out=$OUTPUT_DIR \
    --go_opt=paths=source_relative \
    --go-grpc_opt=require_unimplemented_servers=false \
    --go-grpc_out=$OUTPUT_DIR \
    --go-grpc_opt=paths=source_relative \
    "$INPUT_FILE"
}

for SERVICE_NAME in "$@"; do
  generate $SERVICE_NAME
done
