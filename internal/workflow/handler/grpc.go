package handler

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/luis12loureiro/neurun/api/gen"
	"github.com/luis12loureiro/neurun/internal/workflow"
	"github.com/luis12loureiro/neurun/internal/workflow/domain"
)

type handler struct {
	pb.UnimplementedWorkflowServiceServer
	s workflow.Service
}

func NewServer(s workflow.Service) pb.WorkflowServiceServer {
	return &handler{s: s}
}

func (h *handler) CreateWorkflow(_ context.Context, in *pb.CreateWorkflowRequest) (*pb.WorkflowResponse, error) {
	err := h.s.Create(domain.Worklow{
		ID:          uuid.NewString(),
		Name:        in.Name,
		Description: *in.Description,
	})
	if err != nil {
		return nil, err
	}
	return &pb.WorkflowResponse{Id: "123", Name: in.GetName()}, nil
}

func (h *handler) GetWorkflow(_ context.Context, in *pb.GetWorkflowRequest) (*pb.WorkflowResponse, error) {
	return &pb.WorkflowResponse{Id: "123", Name: "test123"}, nil
}
