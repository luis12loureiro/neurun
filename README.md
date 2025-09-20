# Neurun

A **workflow automation engine** implemented as a **gRPC server in Go**. It enables users to define workflows composed of multiple tasks with event triggers, retries, delays, and conditional execution.

## Features

- Define workflows with multiple tasks  
- Event-driven triggers to start workflows  
- Support for retries and delays on tasks  
- Conditional task execution  
- gRPC API for easy integration and extensibility  

## Getting Started

### Prerequisites

- Go 1.20+  
- Protocol Buffers Compiler (`protoc`)  
- `protoc-gen-go` and `protoc-gen-go-grpc` plugins installed  

### Install Protocol Buffers Plugins

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Generate gRPC Code

```bash
./scripts/generate_proto.sh
```

### Build & Run

```bash
go build -o bin/server ./cmd/server
./bin/server
```

## Project Structure

- `api/` – Protobuf definitions and generated code
  - `proto/` – `.proto` files
  - `gen/` – Generated Go code from `protoc`
- `cmd/server/` – Main entry point for the gRPC server (`main.go`)
- `internal/` – Core business logic and gRPC service implementation
- `scripts/` – Helper scripts (e.g., code generation)
- `go.mod` – Go module definition
- `go.sum` – Dependency checksums
- `README.md` – Project documentation

## License

This project is licensed under the [Attribution-NonCommercial 4.0 International (CC BY-NC 4.0)](https://creativecommons.org/licenses/by-nc/4.0/) License.

© 2025 Luís Loureiro


