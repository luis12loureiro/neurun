#!/bin/bash

# Set root directory to the script location's parent
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

protoc \
  --proto_path="${ROOT_DIR}/api/proto" \
  --go_out="${ROOT_DIR}/api/gen" \
  --go-grpc_out="${ROOT_DIR}/api/gen" \
  --go_opt=paths=source_relative \
  --go-grpc_opt=paths=source_relative \
  "${ROOT_DIR}/api/proto/common.proto" \
  "${ROOT_DIR}/api/proto/workflow.proto"