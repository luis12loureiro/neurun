#!/bin/bash

# Script to generate TypeScript gRPC-Web client
PROTO_DIR="../proto"
OUT_DIR="../apps/frontend/src/app/proto-gen"

mkdir -p $OUT_DIR

# Generate JavaScript and TypeScript definitions
protoc \
  --plugin=protoc-gen-grpc-web=/usr/local/bin/protoc-gen-grpc-web \
  --js_out=import_style=commonjs:$OUT_DIR \
  --grpc-web_out=import_style=typescript,mode=grpcwebtext:$OUT_DIR \
  -I $PROTO_DIR \
  $PROTO_DIR/workflow.proto \
  $PROTO_DIR/task.proto