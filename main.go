package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/luis12loureiro/neurun/api/gen"
	ws "github.com/luis12loureiro/neurun/internal/workflow"
	wh "github.com/luis12loureiro/neurun/internal/workflow/handler"
	wr "github.com/luis12loureiro/neurun/internal/workflow/repository"
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
	repo := wr.NewJSONRepository("./internal/workflow/repository/storage", "data.ndjson")
	svc := ws.NewService(repo)
	handler := wh.NewServer(svc)
	pb.RegisterWorkflowServiceServer(s, handler)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
