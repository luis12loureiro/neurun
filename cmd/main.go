package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/luis12loureiro/neurun/api/gen"
	grpcimpl "github.com/luis12loureiro/neurun/internal/task"
	"github.com/luis12loureiro/neurun/internal/task/repository"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	repo := repository.NewJSONRepository("", "")
	svc := grpcimpl.NewService(repo)
	pb.RegisterWorkflowServiceServer(s, svc)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
