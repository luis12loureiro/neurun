package grpc

import (
	"context"
	"fmt"

	pb "github.com/luis12loureiro/neurun/api/gen"
)

type server struct {
	pb.UnimplementedWorkflowServiceServer
}

func NewServer() pb.WorkflowServiceServer {
	return &server{}
}

func (s *server) CreateWorkflow(_ context.Context, in *pb.CreateWorkflowRequest) (*pb.WorkflowResponse, error) {
	fmt.Printf("Received: %v", in.GetName())
	return &pb.WorkflowResponse{Id: "123", Name: in.GetName()}, nil
}
