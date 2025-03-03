package service

import (
	"context"
	"fmt"

	"github.com/trashwbin/dymall/app/scheduler/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/scheduler/biz/dal/redis"
	scheduler "github.com/trashwbin/dymall/rpc_gen/kitex_gen/scheduler"
)

type GetTaskService struct {
	ctx       context.Context
	mysqlRepo *mysql.TaskRepo
	redisRepo *redis.TaskRepo
}

func NewGetTaskService(ctx context.Context) *GetTaskService {
	return &GetTaskService{
		ctx:       ctx,
		mysqlRepo: mysql.NewTaskRepo(),
		redisRepo: redis.NewTaskRepo(),
	}
}

func (s *GetTaskService) Run(req *scheduler.GetTaskReq) (resp *scheduler.GetTaskResp, err error) {
	// 1. 尝试从缓存获取
	task, err := s.redisRepo.GetTask(s.ctx, req.TaskId)
	if err == nil {
		// 缓存命中
		return &scheduler.GetTaskResp{
			Task: &scheduler.Task{
				TaskId:       task.ID,
				Type:         scheduler.TaskType(task.Type),
				Status:       scheduler.TaskStatus(task.Status),
				TargetId:     task.TargetID,
				ExecuteAt:    task.ExecuteAt.Unix(),
				CreatedAt:    task.CreatedAt.Unix(),
				UpdatedAt:    task.UpdatedAt.Unix(),
				ErrorMessage: task.ErrorMessage,
				Metadata:     task.Metadata,
			},
		}, nil
	}

	// 2. 从数据库获取
	task, err = s.mysqlRepo.GetTask(req.TaskId)
	if err != nil {
		return nil, fmt.Errorf("get task failed: %w", err)
	}

	// 3. 更新缓存
	if err := s.redisRepo.SetTask(s.ctx, task); err != nil {
		// 缓存更新失败只记录日志
		fmt.Printf("set task cache failed: %v\n", err)
	}

	return &scheduler.GetTaskResp{
		Task: &scheduler.Task{
			TaskId:       task.ID,
			Type:         scheduler.TaskType(task.Type),
			Status:       scheduler.TaskStatus(task.Status),
			TargetId:     task.TargetID,
			ExecuteAt:    task.ExecuteAt.Unix(),
			CreatedAt:    task.CreatedAt.Unix(),
			UpdatedAt:    task.UpdatedAt.Unix(),
			ErrorMessage: task.ErrorMessage,
			Metadata:     task.Metadata,
		},
	}, nil
}
