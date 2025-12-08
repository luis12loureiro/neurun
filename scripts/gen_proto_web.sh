#!/bin/bash

PROTO_DIR="../proto"
OUT_DIR="../apps/frontend/src/app/proto-gen"

mkdir -p $OUT_DIR
rm -rf $OUT_DIR/*

# Use protoc bundled with grpc-web (not system protoc)
PROTOC=../apps/frontend/node_modules/.bin/protoc
GRPC_WEB_PLUGIN=../apps/frontend/node_modules/.bin/protoc-gen-grpc-web

$PROTOC -I=$PROTO_DIR \
  --js_out=import_style=commonjs,binary:$OUT_DIR \
  --grpc-web_out=import_style=typescript,mode=grpcweb:$OUT_DIR \
  $PROTO_DIR/*.proto
