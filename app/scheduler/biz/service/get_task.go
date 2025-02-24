package service

import (
	"context"
	scheduler "github.com/trashwbin/dymall/rpc_gen/kitex_gen/scheduler"
)

type GetTaskService struct {
	ctx context.Context
} // NewGetTaskService new GetTaskService
func NewGetTaskService(ctx context.Context) *GetTaskService {
	return &GetTaskService{ctx: ctx}
}

// Run create note info
func (s *GetTaskService) Run(req *scheduler.GetTaskReq) (resp *scheduler.GetTaskResp, err error) {
	// Finish your business logic.

	return
}
