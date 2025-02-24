package main

import (
	"context"

	"github.com/trashwbin/dymall/app/scheduler/biz/service"
	scheduler "github.com/trashwbin/dymall/rpc_gen/kitex_gen/scheduler"
)

// SchedulerServiceImpl implements the last service interface defined in the IDL.
type SchedulerServiceImpl struct{}

// CreateTask implements the SchedulerServiceImpl interface.
func (s *SchedulerServiceImpl) CreateTask(ctx context.Context, req *scheduler.CreateTaskReq) (resp *scheduler.CreateTaskResp, err error) {
	resp, err = service.NewCreateTaskService(ctx).Run(req)

	return resp, err
}

// CancelTask implements the SchedulerServiceImpl interface.
func (s *SchedulerServiceImpl) CancelTask(ctx context.Context, req *scheduler.CancelTaskReq) (resp *scheduler.CancelTaskResp, err error) {
	resp, err = service.NewCancelTaskService(ctx).Run(req)

	return resp, err
}

// GetTask implements the SchedulerServiceImpl interface.
func (s *SchedulerServiceImpl) GetTask(ctx context.Context, req *scheduler.GetTaskReq) (resp *scheduler.GetTaskResp, err error) {
	resp, err = service.NewGetTaskService(ctx).Run(req)

	return resp, err
}

// ExecuteTask implements the SchedulerServiceImpl interface.
func (s *SchedulerServiceImpl) ExecuteTask(ctx context.Context, req *scheduler.ExecuteTaskReq) (resp *scheduler.ExecuteTaskResp, err error) {
	resp, err = service.NewExecuteTaskService(ctx).Run(req)

	return resp, err
}
