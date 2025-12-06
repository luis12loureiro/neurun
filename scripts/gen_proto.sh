#!/bin/bash

# Check if app name is provided
if [ -z "$1" ]; then
  echo "Usage: $0 <app-name>"
  echo "Example: $0 workflow"
  exit 1
fi

APP_NAME=$1

# Set root directory to the script location's parent
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Check if app directory exists
if [ ! -d "${ROOT_DIR}/apps/${APP_NAME}" ]; then
  echo "Error: App directory '${ROOT_DIR}/apps/${APP_NAME}' does not exist"
  exit 1
fi

# Create gen folder if it doesn't exist
mkdir -p "${ROOT_DIR}/apps/${APP_NAME}/gen"

protoc \
  --proto_path="${ROOT_DIR}/proto" \
  --go_out="${ROOT_DIR}/apps/${APP_NAME}/gen" \
  --go-grpc_out="${ROOT_DIR}/apps/${APP_NAME}/gen" \
  --go_opt=paths=source_relative \
  --go-grpc_opt=paths=source_relative \
  "${ROOT_DIR}/proto/task.proto" \
  "${ROOT_DIR}/proto/workflow.proto"