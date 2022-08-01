#!/bin/bash
set -e

function generate {
  local readonly SERVICE_NAME=$1
  local readonly SERVICE_PATH=services/$SERVICE_NAME

  local readonly PROTO_PATH=$SERVICE_PATH/api
  local readonly INPUT_FILE=$PROTO_PATH/$SERVICE_NAME.proto
  local readonly OUTPUT_DIR=clients/python/cryptellation/_genproto

  mkdir -p $OUTPUT_DIR

  python3 \
    -m grpc_tools.protoc \
    --proto_path=$PROTO_PATH \
    --python_out=$OUTPUT_DIR \
    --grpc_python_out=$OUTPUT_DIR \
    "$INPUT_FILE"

    sed -i 's/^import .*_pb2 as/from . \0/' clients/python/cryptellation/_genproto/*pb2_grpc.py
}

for SERVICE_NAME in "$@"; do
  generate $SERVICE_NAME
done
