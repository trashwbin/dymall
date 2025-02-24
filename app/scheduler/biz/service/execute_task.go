package service

import (
	"context"
	scheduler "github.com/trashwbin/dymall/rpc_gen/kitex_gen/scheduler"
)

type ExecuteTaskService struct {
	ctx context.Context
} // NewExecuteTaskService new ExecuteTaskService
func NewExecuteTaskService(ctx context.Context) *ExecuteTaskService {
	return &ExecuteTaskService{ctx: ctx}
}

// Run create note info
func (s *ExecuteTaskService) Run(req *scheduler.ExecuteTaskReq) (resp *scheduler.ExecuteTaskResp, err error) {
	// Finish your business logic.

	return
}
