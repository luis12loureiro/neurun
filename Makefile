.DEFAULT_GOAL := build

proto:
	@echo "Generating protobuf and gRPC code..."
	./scripts/gen_proto.sh
.PHONY:proto

fmt: proto
	@echo "Formatting Go code..."
	go fmt ./...
.PHONY:fmt

lint: fmt
	@echo "Linting Go code..."
	golint ./...
.PHONY:lint

build: lint
	@echo "Building application..."
	go build -o bin/neurun ./cmd/main.go
.PHONY:build