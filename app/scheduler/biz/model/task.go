package model

import (
	"time"
)

// TaskType 任务类型
type TaskType int32

const (
	TaskTypeUnspecified     TaskType = 0
	TaskTypeOrderExpiration TaskType = 1 // 订单过期（同时处理支付过期）
)

// TaskStatus 任务状态
type TaskStatus int32

const (
	TaskStatusUnspecified TaskStatus = 0
	TaskStatusPending     TaskStatus = 1 // 待执行
	TaskStatusRunning     TaskStatus = 2 // 执行中
	TaskStatusCompleted   TaskStatus = 3 // 已完成
	TaskStatusFailed      TaskStatus = 4 // 执行失败
	TaskStatusCancelled   TaskStatus = 5 // 已取消
)

// Task 定时任务领域模型
type Task struct {
	ID           string
	Type         TaskType
	Status       TaskStatus
	TargetID     string    // 目标ID（订单ID）
	ExecuteAt    time.Time // 执行时间
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ErrorMessage string            // 执行失败原因
	Metadata     map[string]string // 额外元数据
}

// IsValid 验证任务是否有效
func (t *Task) IsValid() bool {
	return t.Type != TaskTypeUnspecified &&
		t.Status != TaskStatusUnspecified &&
		t.TargetID != "" &&
		!t.ExecuteAt.IsZero()
}

// IsExecutable 检查任务是否可执行
func (t *Task) IsExecutable() bool {
	return t.Status == TaskStatusPending &&
		time.Now().After(t.ExecuteAt)
}

// MarkRunning 标记任务为执行中
func (t *Task) MarkRunning() {
	t.Status = TaskStatusRunning
	t.UpdatedAt = time.Now()
}

// MarkCompleted 标记任务为已完成
func (t *Task) MarkCompleted() {
	t.Status = TaskStatusCompleted
	t.UpdatedAt = time.Now()
}

// MarkFailed 标记任务为执行失败
func (t *Task) MarkFailed(errMsg string) {
	t.Status = TaskStatusFailed
	t.ErrorMessage = errMsg
	t.UpdatedAt = time.Now()
}

// MarkCancelled 标记任务为已取消
func (t *Task) MarkCancelled() {
	t.Status = TaskStatusCancelled
	t.UpdatedAt = time.Now()
}
