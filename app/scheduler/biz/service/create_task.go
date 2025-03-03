package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/trashwbin/dymall/app/scheduler/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/scheduler/biz/dal/redis"
	"github.com/trashwbin/dymall/app/scheduler/biz/model"
	scheduler "github.com/trashwbin/dymall/rpc_gen/kitex_gen/scheduler"
)

type CreateTaskService struct {
	ctx       context.Context
	mysqlRepo *mysql.TaskRepo
	redisRepo *redis.TaskRepo
}

// NewCreateTaskService new CreateTaskService
func NewCreateTaskService(ctx context.Context) *CreateTaskService {
	return &CreateTaskService{
		ctx:       ctx,
		mysqlRepo: mysql.NewTaskRepo(),
		redisRepo: redis.NewTaskRepo(),
	}
}

// Run create task
func (s *CreateTaskService) Run(req *scheduler.CreateTaskReq) (resp *scheduler.CreateTaskResp, err error) {
	// 1. 参数校验
	if err := s.validateRequest(req); err != nil {
		return nil, err
	}

	// 2. 创建任务领域模型
	task := &model.Task{
		ID:        uuid.New().String(),
		Type:      model.TaskType(req.Type),
		Status:    model.TaskStatusPending,
		TargetID:  req.TargetId,
		ExecuteAt: time.Unix(req.ExecuteAt, 0),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Metadata:  req.Metadata,
	}

	// 3. 验证任务模型
	if !task.IsValid() {
		return nil, fmt.Errorf("invalid task: %v", task)
	}

	// 4. 保存到MySQL
	if err := s.mysqlRepo.CreateTask(task); err != nil {
		return nil, fmt.Errorf("failed to create task in mysql: %w", err)
	}

	// 5. 保存到Redis缓存
	if err := s.redisRepo.SetTask(s.ctx, task); err != nil {
		// Redis缓存失败不影响主流程，只记录错误
		fmt.Printf("failed to cache task in redis: %v\n", err)
	}

	// 6. 构建响应
	resp = &scheduler.CreateTaskResp{
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
	}

	return resp, nil
}

// validateRequest 验证请求参数
func (s *CreateTaskService) validateRequest(req *scheduler.CreateTaskReq) error {
	if req == nil {
		return fmt.Errorf("request is nil")
	}

	if scheduler.TaskType(req.Type) != scheduler.TaskType(model.TaskTypeOrderExpiration) {
		return fmt.Errorf("unsupported task type: %d", req.Type)
	}

	if req.TargetId == "" {
		return fmt.Errorf("target_id is required")
	}

	if req.ExecuteAt <= 0 {
		return fmt.Errorf("invalid execute_at: %d", req.ExecuteAt)
	}

	// 如果是测试任务，跳过执行时间验证
	if req.Metadata != nil && req.Metadata["test"] == "true" {
		return nil
	}

	executeAt := time.Unix(req.ExecuteAt, 0)
	if executeAt.Before(time.Now()) {
		return fmt.Errorf("execute_at must be in the future")
	}

	return nil
}
