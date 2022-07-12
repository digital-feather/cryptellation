#!/bin/bash
set -e

function generate_golang {
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

function generate_python {
  local readonly SERVICE_NAME=$1
  local readonly OUTPUT_DIR=clients/python/cryptellation/_genproto

  mkdir -p $OUTPUT_DIR

  python3 \
    -m grpc_tools.protoc \
    --proto_path=api/protobuf \
    --python_out=$OUTPUT_DIR \
    --grpc_python_out=$OUTPUT_DIR \
    "api/protobuf/$SERVICE_NAME.proto"

    sed -i 's/^import .*_pb2 as/from . \0/' clients/python/cryptellation/_genproto/*pb2_grpc.py
}

for SERVICE_NAME in "$@"; do
  generate_golang $SERVICE_NAME
  generate_python $SERVICE_NAME
done
