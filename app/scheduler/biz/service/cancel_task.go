package service

import (
	"context"
	"fmt"

	"github.com/trashwbin/dymall/app/scheduler/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/scheduler/biz/dal/redis"
	"github.com/trashwbin/dymall/app/scheduler/biz/model"
	scheduler "github.com/trashwbin/dymall/rpc_gen/kitex_gen/scheduler"
)

type CancelTaskService struct {
	ctx       context.Context
	mysqlRepo *mysql.TaskRepo
	redisRepo *redis.TaskRepo
}

func NewCancelTaskService(ctx context.Context) *CancelTaskService {
	return &CancelTaskService{
		ctx:       ctx,
		mysqlRepo: mysql.NewTaskRepo(),
		redisRepo: redis.NewTaskRepo(),
	}
}

func (s *CancelTaskService) Run(req *scheduler.CancelTaskReq) (resp *scheduler.CancelTaskResp, err error) {
	// 1. 获取任务
	task, err := s.mysqlRepo.GetTask(req.TaskId)
	if err != nil {
		return nil, fmt.Errorf("get task failed: %w", err)
	}

	// 2. 检查任务状态
	if task.Status == model.TaskStatusCompleted || task.Status == model.TaskStatusCancelled {
		return nil, fmt.Errorf("task already completed or cancelled")
	}

	// 3. 标记任务为已取消
	task.MarkCancelled()
	if err := s.mysqlRepo.UpdateTask(task); err != nil {
		return nil, fmt.Errorf("update task status failed: %w", err)
	}

	// 4. 删除缓存
	if err := s.redisRepo.DeleteTask(s.ctx, task.ID); err != nil {
		// 缓存删除失败只记录日志
		fmt.Printf("delete task cache failed: %v\n", err)
	}

	return &scheduler.CancelTaskResp{}, nil
}
