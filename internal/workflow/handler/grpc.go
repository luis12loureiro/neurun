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
	var tasks []domain.Task
	for _, t := range in.Tasks {
		task, err := TaskFromProto(t)
		if err != nil {
			return nil, err
		}
		task.ID = uuid.NewString()
		tasks = append(tasks, task)
	}
	id := uuid.NewString()
	err := h.s.Create(domain.Worklow{
		ID:          id,
		Name:        in.GetName(),
		Description: in.GetDescription(),
		Status:      domain.WorkflowStatusIDLE,
		Tasks:       tasks,
	})
	if err != nil {
		return nil, err
	}
	return &pb.WorkflowResponse{Id: id, Name: in.GetName()}, nil
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
