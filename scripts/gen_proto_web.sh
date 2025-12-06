#!/bin/bash

PROTO_DIR="../proto"
OUT_DIR="../apps/frontend/src/app/proto-gen"

mkdir -p $OUT_DIR
rm -rf $OUT_DIR/*

protoc -I=$PROTO_DIR \
  --js_out=import_style=commonjs:$OUT_DIR \
  --grpc-web_out=import_style=typescript,mode=grpcwebtext:$OUT_DIR \
  $PROTO_DIR/*.proto