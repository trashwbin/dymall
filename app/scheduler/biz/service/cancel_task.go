package service

import (
	"context"
	scheduler "github.com/trashwbin/dymall/rpc_gen/kitex_gen/scheduler"
)

type CancelTaskService struct {
	ctx context.Context
} // NewCancelTaskService new CancelTaskService
func NewCancelTaskService(ctx context.Context) *CancelTaskService {
	return &CancelTaskService{ctx: ctx}
}

// Run create note info
func (s *CancelTaskService) Run(req *scheduler.CancelTaskReq) (resp *scheduler.CancelTaskResp, err error) {
	// Finish your business logic.

	return
}
