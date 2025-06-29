protoc \                     
  --proto_path=api/proto \
  --go_out=api/gen \
  --go-grpc_out=api/gen \
  --go_opt=paths=source_relative \
  --go-grpc_opt=paths=source_relative \
  neurun.proto