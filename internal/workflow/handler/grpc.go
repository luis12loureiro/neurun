package handler

import (
	"context"

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
	tasks := []*domain.Task{}
	for _, t := range in.Tasks {
		task, err := TaskFromProto(t)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	w, err := domain.NewWorkflow(in.GetName(), in.GetDescription(), tasks)
	if err != nil {
		return nil, err
	}
	wf, err := h.s.Create(w)
	if err != nil {
		return nil, err
	}
	return &pb.WorkflowResponse{
		Id:          wf.ID,
		Name:        wf.Name,
		Description: wf.Description,
		Status:      string(wf.Status),
		Tasks:       convertNextToProto(wf.Tasks),
	}, nil
}

func (h *handler) GetWorkflow(_ context.Context, in *pb.GetWorkflowRequest) (*pb.WorkflowResponse, error) {
	wf, err := h.s.Get(in.GetId())
	if err != nil {
		return nil, err
	}
	tasks := make([]*pb.Task, len(wf.Tasks))
	for i, t := range wf.Tasks {
		tasks[i] = TaskToProto(t)
	}
	return &pb.WorkflowResponse{
		Id:          wf.ID,
		Name:        wf.Name,
		Description: wf.Description,
		Status:      string(wf.Status),
		Tasks:       tasks,
	}, nil
}
