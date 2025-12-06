.DEFAULT_GOAL := build

proto:
	@echo "Generating protobuf and gRPC code..."
	../../scripts/gen_proto.sh workflow
.PHONY:proto

fmt: proto
	@echo "Formatting Go code..."
	go fmt ./...
.PHONY:fmt

lint: fmt
	@echo "Linting Go code..."
	golint ./...
.PHONY:lint

vet: lint
	@echo "Vetting Go code..."
	go vet ./...
.PHONY:vet

build: vet
	@echo "Building application..."
	go build -o bin/neurun ./main.go
.PHONY:build