package service

import (
	"context"
	scheduler "github.com/trashwbin/dymall/rpc_gen/kitex_gen/scheduler"
)

type CreateTaskService struct {
	ctx context.Context
} // NewCreateTaskService new CreateTaskService
func NewCreateTaskService(ctx context.Context) *CreateTaskService {
	return &CreateTaskService{ctx: ctx}
}

// Run create note info
func (s *CreateTaskService) Run(req *scheduler.CreateTaskReq) (resp *scheduler.CreateTaskResp, err error) {
	// Finish your business logic.

	return
}
