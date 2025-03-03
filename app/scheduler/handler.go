package main

import (
	"context"

	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/trashwbin/dymall/app/scheduler/biz/service"
	"github.com/trashwbin/dymall/app/scheduler/middleware"
	scheduler "github.com/trashwbin/dymall/rpc_gen/kitex_gen/scheduler"
)

// SchedulerServiceImpl implements the last service interface defined in the IDL.
type SchedulerServiceImpl struct{}

// GetMiddlewares 获取中间件配置
func (s *SchedulerServiceImpl) GetMiddlewares() []endpoint.Middleware {
	return []endpoint.Middleware{
		middleware.UserAuthMiddleware(),
	}
}

// GetTaskMiddlewares 获取任务相关接口的中间件
func (s *SchedulerServiceImpl) GetTaskMiddlewares() []endpoint.Middleware {
	return []endpoint.Middleware{
		middleware.ServiceAuthMiddleware(middleware.OrderService, middleware.PaymentService),
	}
}

// CreateTask implements the SchedulerServiceImpl interface.
func (s *SchedulerServiceImpl) CreateTask(ctx context.Context, req *scheduler.CreateTaskReq) (resp *scheduler.CreateTaskResp, err error) {
	// 只允许订单服务和支付服务创建任务
	for _, m := range s.GetTaskMiddlewares() {
		if m != nil {
			endpoint := m(func(ctx context.Context, req, resp interface{}) error {
				return nil
			})
			if err := endpoint(ctx, req, resp); err != nil {
				return nil, err
			}
		}
	}

	resp, err = service.NewCreateTaskService(ctx).Run(req)
	return resp, err
}

// CancelTask implements the SchedulerServiceImpl interface.
func (s *SchedulerServiceImpl) CancelTask(ctx context.Context, req *scheduler.CancelTaskReq) (resp *scheduler.CancelTaskResp, err error) {
	// 只允许订单服务和支付服务取消任务
	for _, m := range s.GetTaskMiddlewares() {
		if m != nil {
			endpoint := m(func(ctx context.Context, req, resp interface{}) error {
				return nil
			})
			if err := endpoint(ctx, req, resp); err != nil {
				return nil, err
			}
		}
	}

	resp, err = service.NewCancelTaskService(ctx).Run(req)
	return resp, err
}

// GetTask implements the SchedulerServiceImpl interface.
func (s *SchedulerServiceImpl) GetTask(ctx context.Context, req *scheduler.GetTaskReq) (resp *scheduler.GetTaskResp, err error) {
	// 需要用户认证
	for _, m := range s.GetMiddlewares() {
		if m != nil {
			endpoint := m(func(ctx context.Context, req, resp interface{}) error {
				return nil
			})
			if err := endpoint(ctx, req, resp); err != nil {
				return nil, err
			}
		}
	}

	resp, err = service.NewGetTaskService(ctx).Run(req)
	return resp, err
}

// ExecuteTask implements the SchedulerServiceImpl interface.
func (s *SchedulerServiceImpl) ExecuteTask(ctx context.Context, req *scheduler.ExecuteTaskReq) (resp *scheduler.ExecuteTaskResp, err error) {
	// 只允许订单服务和支付服务执行任务
	for _, m := range s.GetTaskMiddlewares() {
		if m != nil {
			endpoint := m(func(ctx context.Context, req, resp interface{}) error {
				return nil
			})
			if err := endpoint(ctx, req, resp); err != nil {
				return nil, err
			}
		}
	}

	resp, err = service.NewExecuteTaskService(ctx).Run(req)
	return resp, err
}
