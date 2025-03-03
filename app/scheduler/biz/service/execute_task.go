package service

import (
	"context"
	"fmt"
	"time"

	"github.com/trashwbin/dymall/app/scheduler/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/scheduler/biz/dal/redis"
	"github.com/trashwbin/dymall/app/scheduler/infra/rpc"
	scheduler "github.com/trashwbin/dymall/rpc_gen/kitex_gen/scheduler"
)

type ExecuteTaskService struct {
	ctx       context.Context
	mysqlRepo *mysql.TaskRepo
	redisRepo *redis.TaskRepo
	executor  *TaskExecutor
}

func NewExecuteTaskService(ctx context.Context) *ExecuteTaskService {
	return &ExecuteTaskService{
		ctx:       ctx,
		mysqlRepo: mysql.NewTaskRepo(),
		redisRepo: redis.NewTaskRepo(),
		executor: &TaskExecutor{
			mysqlRepo:     mysql.NewTaskRepo(),
			redisRepo:     redis.NewTaskRepo(),
			orderClient:   rpc.OrderClient,
			paymentClient: rpc.PaymentClient,
		},
	}
}

func (s *ExecuteTaskService) Run(req *scheduler.ExecuteTaskReq) (resp *scheduler.ExecuteTaskResp, err error) {
	// 1. 获取任务
	task, err := s.mysqlRepo.GetTask(req.TaskId)
	if err != nil {
		return nil, fmt.Errorf("get task failed: %w", err)
	}

	// 2. 检查任务状态
	if !task.IsExecutable() {
		return nil, fmt.Errorf("task is not executable")
	}

	// 3. 尝试获取任务锁
	locked, err := s.redisRepo.AcquireLock(s.ctx, task.ID)
	if err != nil {
		return nil, fmt.Errorf("acquire task lock failed: %w", err)
	}
	if !locked {
		return nil, fmt.Errorf("task is being executed by another instance")
	}
	defer s.redisRepo.ReleaseLock(s.ctx, task.ID)

	// 4. 更新任务状态为执行中
	task.MarkRunning()
	if err := s.mysqlRepo.UpdateTask(task); err != nil {
		return nil, fmt.Errorf("update task status failed: %w", err)
	}

	// 5. 执行任务（同步执行）
	err = s.executor.handleOrderExpiration(s.ctx, task)
	if err != nil {
		task.MarkFailed(err.Error())
	} else {
		task.MarkCompleted()
	}

	// 6. 更新任务状态
	if err := s.mysqlRepo.UpdateTask(task); err != nil {
		return nil, fmt.Errorf("update task status failed: %w", err)
	}

	// 7. 等待一段时间确保数据库更新完成
	time.Sleep(time.Millisecond * 100)

	// 8. 获取最新的任务状态（确保读取到最新状态）
	task, err = s.mysqlRepo.GetTask(req.TaskId)
	if err != nil {
		return nil, fmt.Errorf("get task status failed: %w", err)
	}

	return &scheduler.ExecuteTaskResp{
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
