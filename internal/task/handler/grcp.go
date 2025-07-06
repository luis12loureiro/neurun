package handler

import (
	"context"
	"fmt"

	pb "github.com/luis12loureiro/neurun/api/gen"
	t "github.com/luis12loureiro/neurun/internal/task"
	task "github.com/luis12loureiro/neurun/internal/task/domain"
)

type handler struct {
	pb.UnimplementedWorkflowServiceServer
	s t.Service
}

func NewServer(s t.Service) pb.WorkflowServiceServer {
	return &handler{s: s}
}

func (h *handler) CreateWorkflow(_ context.Context, in *pb.CreateWorkflowRequest) (*pb.WorkflowResponse, error) {
	fmt.Printf("Received: %v", in.GetName())
	h.s.Create(task.Task{Id: "123", Name: "task1", Type: "log", Delay: 2})
	return &pb.WorkflowResponse{Id: "123", Name: in.GetName()}, nil
}
